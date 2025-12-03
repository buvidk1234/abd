package service

import (
	"backend/internal/dto"
	"backend/internal/model"
	"context"
	"strconv"

	"gorm.io/gorm"
)

type FriendService struct {
	db *gorm.DB
}

func NewFriendService(db *gorm.DB) *FriendService {
	return &FriendService{db: db}
}

// 申请添加好友
type ApplyToAddFriendReq struct {
	FromUserID int64  `json:"fromUserID,string" binding:"required"`
	ToUserID   int64  `json:"toUserID,string" binding:"required"`
	ReqMsg     string `json:"message"`
}

func (s *FriendService) ApplyToAddFriend(ctx context.Context, req ApplyToAddFriendReq) error {
	s.db.WithContext(ctx).Create(&model.FriendRequest{
		FromUserID: req.FromUserID,
		ToUserID:   req.ToUserID,
		ReqMsg:     req.ReqMsg,
	})
	return nil
}

type RespondFriendApplyReq struct {
	ID            int64  `json:"id,string" binding:"required"`
	HandlerUserID int64  `json:"handlerUserID,string" binding:"required"`
	HandleResult  int32  `json:"handleResult,string" binding:"required"`
	HandleMsg     string `json:"handleMsg"`
}

// 响应好友申请
func (s *FriendService) RespondFriendApply(ctx context.Context, req RespondFriendApplyReq) error {
	// 查找好友请求
	var fr model.FriendRequest
	if err := s.db.WithContext(ctx).First(&fr, req.ID).Error; err != nil {
		return err
	}
	// 更新处理结果
	updateStruct := model.FriendRequest{
		HandleResult:  req.HandleResult,
		HandlerUserID: req.HandlerUserID,
		HandleMsg:     req.HandleMsg,
		HandledAt:     s.db.NowFunc(),
	}
	if err := s.db.WithContext(ctx).Model(&fr).Updates(updateStruct).Error; err != nil {
		return err
	}
	// 如果同意，建立好友关系
	if req.HandleResult == 1 {
		s.db.WithContext(ctx).Create(&model.Friend{
			OwnerUserID:  fr.FromUserID,
			FriendUserID: fr.ToUserID,
		})
		s.db.WithContext(ctx).Create(&model.Friend{
			OwnerUserID:  fr.ToUserID,
			FriendUserID: fr.FromUserID,
		})
	}
	return nil
}

type GetPaginationFriendsReq struct {
	Page     int `json:"page,default=1"`
	PageSize int `json:"pageSize,default=10"`
}

type GetFriendListResp struct {
	Friends  []dto.FriendInfo `json:"friends"`
	Total    int              `json:"total"`
	Page     int              `json:"page"`
	PageSize int              `json:"pageSize"`
}

// 获取好友列表
func (s *FriendService) GetPaginationFriends(ctx context.Context, params GetPaginationFriendsReq, userId string) (GetFriendListResp, error) {
	// 默认值处理
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 10
	}

	var total int64
	var friends []model.Friend
	db := s.db.WithContext(ctx).Model(&model.Friend{}).Where("owner_user_id = ?", userId)
	db.Count(&total)
	db = db.Offset((params.Page - 1) * params.PageSize).Limit(params.PageSize)
	if err := db.Find(&friends).Error; err != nil {
		return GetFriendListResp{}, err
	}

	// 批量查好友用户信息
	friendIDs := make([]int64, 0, len(friends))
	for _, f := range friends {
		friendIDs = append(friendIDs, f.FriendUserID)
	}
	var users []model.User
	if len(friendIDs) > 0 {
		if err := s.db.WithContext(ctx).Where("user_id IN ?", friendIDs).Find(&users).Error; err != nil {
			return GetFriendListResp{}, err
		}
	}
	userMap := make(map[int64]model.User)
	for _, user := range users {
		userMap[user.UserID] = user
	}
	var friendInfos []dto.FriendInfo
	for _, f := range friends {
		friendInfos = append(friendInfos, dto.ConvertToFriendInfo(f, userMap[f.FriendUserID]))
	}
	return GetFriendListResp{
		Friends:  friendInfos,
		Total:    int(total),
		Page:     params.Page,
		PageSize: params.PageSize,
	}, nil
}

// 获取指定好友信息
func (s *FriendService) GetSpecifiedFriendInfo(ctx context.Context, ownerUserID, friendUserID string) (dto.FriendInfo, error) {
	var f model.Friend
	if err := s.db.WithContext(ctx).Where("owner_user_id = ? AND friend_user_id = ?", ownerUserID, friendUserID).First(&f).Error; err != nil {
		return dto.FriendInfo{}, err
	}
	var u model.User
	if err := s.db.WithContext(ctx).Where("user_id = ?", friendUserID).First(&u).Error; err != nil {
		return dto.FriendInfo{}, err
	}
	return dto.ConvertToFriendInfo(f, u), nil
}

type DeleteFriendReq struct {
	OwnerUserID  string `json:"ownerUserID" binding:"required"`
	FriendUserID string `json:"friendUserID" binding:"required"`
}

// 删除好友，TODO：未测试
func (s *FriendService) DeleteFriend(ctx context.Context, ownerUserID string, friendUserID string) error {
	if err := s.db.WithContext(ctx).Where("owner_user_id = ? AND friend_user_id = ?", ownerUserID, friendUserID).Delete(&model.Friend{}).Error; err != nil {
		return err
	}
	if err := s.db.WithContext(ctx).Where("owner_user_id = ? AND friend_user_id = ?", friendUserID, ownerUserID).Delete(&model.Friend{}).Error; err != nil {
		return err
	}
	return nil
}

type GetPaginationFriendApplyListReq struct {
	Page     int    `json:"page"`
	PageSize int    `json:"pageSize"`
	ToUserID string `json:"toUserID" binding:"required"`
}
type GetPaginationFriendApplyListResp struct {
	List     []dto.FriendRequestInfo `json:"list"`
	Total    int                     `json:"total"`
	Page     int                     `json:"page"`
	PageSize int                     `json:"pageSize"`
}

// 获取收到的好友申请列表
func (s *FriendService) GetPaginationFriendApplyList(ctx context.Context, req GetPaginationFriendApplyListReq) (GetPaginationFriendApplyListResp, error) {
	// 默认值处理
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}
	var total int64
	db := s.db.WithContext(ctx).Model(&model.FriendRequest{}).Where("to_user_id = ?", req.ToUserID)
	db.Count(&total)
	var frs []model.FriendRequest
	db = db.Order("created_at desc").Offset((req.Page - 1) * req.PageSize).Limit(req.PageSize)
	if err := db.Find(&frs).Error; err != nil {
		return GetPaginationFriendApplyListResp{}, err
	}
	fromIDs := make([]int64, 0, len(frs))
	for _, fr := range frs {
		fromIDs = append(fromIDs, fr.FromUserID)
	}
	var users []model.User
	if len(fromIDs) > 0 {
		if err := s.db.WithContext(ctx).Where("user_id IN ?", fromIDs).Find(&users).Error; err != nil {
			return GetPaginationFriendApplyListResp{}, err
		}
	}
	userMap := make(map[int64]model.User)
	for _, u := range users {
		userMap[u.UserID] = u
	}
	var list []dto.FriendRequestInfo
	for _, fr := range frs {
		fromUser := userMap[fr.FromUserID]
		list = append(list, dto.ConvertToFriendRequestInfo(fr, fromUser, model.User{}))
	}
	return GetPaginationFriendApplyListResp{
		List:     list,
		Total:    int(total),
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

type GetSelfFriendApplyListResp struct {
	List     []dto.FriendRequestInfo `json:"list"`
	Total    int                     `json:"total"`
	Page     int                     `json:"page"`
	PageSize int                     `json:"pageSize"`
}

type GetPaginationSelfFriendApplyListReq struct {
	Page       int    `json:"page"`
	PageSize   int    `json:"pageSize"`
	FromUserID string `json:"fromUserID" binding:"required"`
}

// 获取自己发出的好友申请列表
func (s *FriendService) GetPaginationSelfApplyList(ctx context.Context, req GetPaginationSelfFriendApplyListReq) (GetSelfFriendApplyListResp, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}
	var total int64
	db := s.db.WithContext(ctx).Model(&model.FriendRequest{}).Where("from_user_id = ?", req.FromUserID)
	db.Count(&total)
	var frs []model.FriendRequest
	db = db.Order("created_at desc").Offset((req.Page - 1) * req.PageSize).Limit(req.PageSize)
	if err := db.Find(&frs).Error; err != nil {
		return GetSelfFriendApplyListResp{}, err
	}
	toIDs := make([]int64, 0, len(frs))
	for _, fr := range frs {
		toIDs = append(toIDs, fr.ToUserID)
	}
	var users []model.User
	if len(toIDs) > 0 {
		if err := s.db.WithContext(ctx).Find(&users, "user_id IN ?", toIDs).Error; err != nil {
			return GetSelfFriendApplyListResp{}, err
		}
	}
	userMap := make(map[int64]model.User)
	for _, u := range users {
		userMap[u.UserID] = u
	}
	var list []dto.FriendRequestInfo
	for _, fr := range frs {
		toUser := userMap[fr.ToUserID]
		list = append(list, dto.ConvertToFriendRequestInfo(fr, model.User{}, toUser))
	}
	return GetSelfFriendApplyListResp{
		List:     list,
		Total:    int(total),
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

// TODO: Black
// 添加黑名单
type AddBlackReq struct {
	OwnerUserID    string `json:"ownerUserID"`
	BlockUserID    string `json:"blockUserID" binding:"required"`
	OperatorUserID string `json:"operatorUserID"`
	AddSource      int32  `json:"addSource,string"`
}

func (s *FriendService) AddBlack(ctx context.Context, req AddBlackReq) error {
	ownerID, err := strconv.ParseInt(req.OwnerUserID, 10, 64)
	if err != nil {
		return err
	}
	blockID, err := strconv.ParseInt(req.BlockUserID, 10, 64)
	if err != nil {
		return err
	}
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var existing int64
		if err := tx.Model(&model.Black{}).Where("owner_user_id = ? AND block_user_id = ?", ownerID, blockID).Count(&existing).Error; err != nil {
			return err
		}
		if existing == 0 {
			black := model.Black{
				OwnerUserID:    ownerID,
				BlockUserID:    blockID,
				AddSource:      req.AddSource,
				OperatorUserID: ownerID,
			}
			if err := tx.Create(&black).Error; err != nil {
				return err
			}
		}
		// 同时移除好友关系，避免黑名单和好友并存
		if err := tx.Where("owner_user_id = ? AND friend_user_id = ?", ownerID, blockID).Delete(&model.Friend{}).Error; err != nil {
			return err
		}
		if err := tx.Where("owner_user_id = ? AND friend_user_id = ?", blockID, ownerID).Delete(&model.Friend{}).Error; err != nil {
			return err
		}
		return nil
	})
}

// 移除黑名单
type RemoveBlackReq struct {
	OwnerUserID string `json:"ownerUserID"`
	BlockUserID string `json:"blockUserID" binding:"required"`
}

func (s *FriendService) RemoveBlack(ctx context.Context, req RemoveBlackReq) error {
	ownerID, err := strconv.ParseInt(req.OwnerUserID, 10, 64)
	if err != nil {
		return err
	}
	blockID, err := strconv.ParseInt(req.BlockUserID, 10, 64)
	if err != nil {
		return err
	}
	return s.db.WithContext(ctx).Where("owner_user_id = ? AND block_user_id = ?", ownerID, blockID).Delete(&model.Black{}).Error
}

// 获取黑名单列表
type GetPaginationBlacksReq struct {
	Page        int    `form:"page,default=1"`
	PageSize    int    `form:"pageSize,default=10"`
	OwnerUserID string `form:"ownerUserID"`
}
type GetPaginationBlacksResp struct {
	Blacks   []dto.BlackInfo `json:"blacks"`
	Total    int             `json:"total"`
	Page     int             `json:"page"`
	PageSize int             `json:"pageSize"`
}

func (s *FriendService) GetPaginationBlacks(ctx context.Context, req GetPaginationBlacksReq) (GetPaginationBlacksResp, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	var total int64
	db := s.db.WithContext(ctx).Model(&model.Black{}).Where("owner_user_id = ?", req.OwnerUserID)
	if err := db.Count(&total).Error; err != nil {
		return GetPaginationBlacksResp{}, err
	}
	var blacks []model.Black
	db = db.Order("created_at desc").Offset((req.Page - 1) * req.PageSize).Limit(req.PageSize)
	if err := db.Find(&blacks).Error; err != nil {
		return GetPaginationBlacksResp{}, err
	}
	blockIDs := make([]int64, 0, len(blacks))
	for _, b := range blacks {
		blockIDs = append(blockIDs, b.BlockUserID)
	}
	var users []model.User
	if len(blockIDs) > 0 {
		if err := s.db.WithContext(ctx).Where("user_id IN ?", blockIDs).Find(&users).Error; err != nil {
			return GetPaginationBlacksResp{}, err
		}
	}
	userMap := make(map[int64]model.User)
	for _, u := range users {
		userMap[u.UserID] = u
	}
	var resp []dto.BlackInfo
	for _, b := range blacks {
		resp = append(resp, dto.ConvertToBlackInfo(b, userMap[b.BlockUserID]))
	}
	return GetPaginationBlacksResp{
		Blacks:   resp,
		Total:    int(total),
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}
