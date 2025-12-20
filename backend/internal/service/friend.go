package service

import (
	"backend/internal/dto"
	"backend/internal/model"
	"context"

	"gorm.io/gorm"
)

type FriendService struct {
	db          *gorm.DB
	userService *UserService
}

func NewFriendService(db *gorm.DB, userService *UserService) *FriendService {
	return &FriendService{db: db, userService: userService}
}

// 申请添加好友
type ApplyToAddFriendReq struct {
	ToUserID int64  `json:"toUserID,string" binding:"required"`
	ReqMsg   string `json:"message"`
}

func (s *FriendService) ApplyToAddFriend(ctx context.Context, userId int64, req ApplyToAddFriendReq) error {
	s.db.WithContext(ctx).Create(&model.FriendRequest{
		FromUserID: userId,
		ToUserID:   req.ToUserID,
		ReqMsg:     req.ReqMsg,
	})
	return nil
}

type RespondFriendApplyReq struct {
	ID           int64  `json:"id,string" binding:"required"`
	HandleResult int32  `json:"handleResult" binding:"required"`
	HandleMsg    string `json:"handleMsg"`
}

// 响应好友申请
func (s *FriendService) RespondFriendApply(ctx context.Context, req RespondFriendApplyReq, userId int64) error {
	// 查找好友请求
	var fr model.FriendRequest
	if err := s.db.WithContext(ctx).First(&fr, req.ID).Error; err != nil {
		return err
	}
	// 更新处理结果
	updateStruct := model.FriendRequest{
		HandleResult:  req.HandleResult,
		HandlerUserID: userId,
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
	PagedParams
}

// 获取好友列表
func (s *FriendService) GetPaginationFriends(ctx context.Context, params GetPaginationFriendsReq, userId int64) (PagedResp[dto.FriendInfo], error) {
	base := s.db.WithContext(ctx).Model(&model.Friend{}).Where("owner_user_id = ?", userId)
	var total int64
	if err := base.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return PagedResp[dto.FriendInfo]{}, err
	}

	page := params.Page
	pageSize := params.PageSize
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	var friends []model.Friend
	err := base.Session(&gorm.Session{}).
		Scopes(model.SelectFriendInfo).
		Preload("FriendUser", model.SelectUserInfo).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&friends).Error
	if err != nil {
		return PagedResp[dto.FriendInfo]{}, err
	}

	friendInfos := make([]dto.FriendInfo, 0, len(friends))
	for _, f := range friends {
		friendInfos = append(friendInfos, dto.ConvertToFriendInfo(f, f.FriendUser))
	}
	return PagedResp[dto.FriendInfo]{
		Page:     page,
		Total:    int(total),
		PageSize: pageSize,
		Data:     friendInfos,
	}, nil
}

// 获取指定好友信息
func (s *FriendService) GetSpecifiedFriendInfo(ctx context.Context, ownerUserID, friendUserID int64) (dto.FriendInfo, error) {
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
	OwnerUserID  int64
	FriendUserID int64 `json:"friendUserID,string" binding:"required"`
}

// 删除好友，TODO：未测试
func (s *FriendService) DeleteFriend(ctx context.Context, ownerUserID int64, friendUserID int64) error {
	if err := s.db.WithContext(ctx).Where("owner_user_id = ? AND friend_user_id = ?", ownerUserID, friendUserID).Delete(&model.Friend{}).Error; err != nil {
		return err
	}
	if err := s.db.WithContext(ctx).Where("owner_user_id = ? AND friend_user_id = ?", friendUserID, ownerUserID).Delete(&model.Friend{}).Error; err != nil {
		return err
	}
	return nil
}

type GetPaginationFriendApplyListParams struct {
	PagedParams
}

// 获取收到的好友申请列表
func (s *FriendService) GetPaginationFriendApplyList(ctx context.Context, id int64, req GetPaginationFriendApplyListParams) (PagedResp[dto.FriendRequestInfo], error) {
	base := s.db.WithContext(ctx).Model(&model.FriendRequest{}).Where("to_user_id = ? OR from_user_id = ?", id, id)
	var total int64
	if err := base.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return PagedResp[dto.FriendRequestInfo]{}, err
	}

	page := req.Page
	pageSize := req.PageSize
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	var frs []model.FriendRequest
	err := base.Session(&gorm.Session{}).
		Scopes(model.SelectFriendRequestInfo).
		Preload("FromUser", model.SelectUserInfo).
		Preload("ToUser", model.SelectUserInfo).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&frs).Error
	if err != nil {
		return PagedResp[dto.FriendRequestInfo]{}, err
	}

	list := make([]dto.FriendRequestInfo, 0, len(frs))
	for _, fr := range frs {
		list = append(list, dto.ConvertToFriendRequestInfo(fr, fr.FromUser, fr.ToUser))
	}

	return PagedResp[dto.FriendRequestInfo]{
		Page:     page,
		Total:    int(total),
		PageSize: pageSize,
		Data:     list,
	}, nil
}

type GetSelfFriendApplyListResp struct {
	List     []dto.FriendRequestInfo `json:"list"`
	Total    int                     `json:"total"`
	Page     int                     `json:"page"`
	PageSize int                     `json:"pageSize"`
}

type GetPaginationSelfFriendApplyListReq struct {
	Page       int `json:"page"`
	PageSize   int `json:"pageSize"`
	FromUserID int64
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
	OwnerUserID    int64
	BlockUserID    int64 `json:"blockUserID,string" binding:"required"`
	OperatorUserID int64 `json:"operatorUserID,string"`
	AddSource      int32 `json:"addSource"`
}

func (s *FriendService) AddBlack(ctx context.Context, req AddBlackReq) error {

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var existing int64
		if err := tx.Model(&model.Black{}).Where("owner_user_id = ? AND block_user_id = ?", req.OwnerUserID, req.BlockUserID).Count(&existing).Error; err != nil {
			return err
		}
		if existing == 0 {
			black := model.Black{
				OwnerUserID:    req.OwnerUserID,
				BlockUserID:    req.BlockUserID,
				AddSource:      req.AddSource,
				OperatorUserID: req.OperatorUserID,
			}
			if err := tx.Create(&black).Error; err != nil {
				return err
			}
		}
		// 同时移除好友关系，避免黑名单和好友并存
		if err := tx.Where("owner_user_id = ? AND friend_user_id = ?", req.OwnerUserID, req.BlockUserID).Delete(&model.Friend{}).Error; err != nil {
			return err
		}
		if err := tx.Where("owner_user_id = ? AND friend_user_id = ?", req.BlockUserID, req.OwnerUserID).Delete(&model.Friend{}).Error; err != nil {
			return err
		}
		return nil
	})
}

// 移除黑名单
type RemoveBlackReq struct {
	OwnerUserID int64
	BlockUserID int64 `json:"blockUserID,string" binding:"required"`
}

func (s *FriendService) RemoveBlack(ctx context.Context, req RemoveBlackReq) error {
	return s.db.WithContext(ctx).Where("owner_user_id = ? AND block_user_id = ?", req.OwnerUserID, req.BlockUserID).Delete(&model.Black{}).Error
}

// 获取黑名单列表
type GetPaginationBlacksReq struct {
	Page        int `form:"page,default=1"`
	PageSize    int `form:"pageSize,default=10"`
	OwnerUserID int64
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

func (s *FriendService) GetSearchedFriendInfo(ctx context.Context, searchId int64) (dto.UserInfo, error) {
	return s.userService.GetUsersPublicInfo(ctx, searchId)
}
