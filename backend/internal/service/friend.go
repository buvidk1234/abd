package service

import (
	"backend/internal/dto"
	"backend/internal/model"
	"context"

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
	FromUserID string `json:"fromUserID" binding:"required"`
	ToUserID   string `json:"toUserID" binding:"required"`
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
	ID            uint   `json:"id" binding:"required"`
	HandlerUserID string `json:"handlerUserID" binding:"required"`
	HandleResult  int32  `json:"handleResult" binding:"required"` // 0未处理，1同意，2拒绝
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
	Page     int    `json:"page"`
	PageSize int    `json:"pageSize"`
	UserID   string `json:"userID" binding:"required"`
}

type GetFriendListResp struct {
	Friends  []dto.FriendInfo `json:"friends"`
	Total    int              `json:"total"`
	Page     int              `json:"page"`
	PageSize int              `json:"pageSize"`
}

// 获取好友列表
func (s *FriendService) GetPaginationFriends(ctx context.Context, req GetPaginationFriendsReq) (GetFriendListResp, error) {
	// 默认值处理
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	var total int64
	var friends []model.Friend
	db := s.db.WithContext(ctx).Model(&model.Friend{}).Where("owner_user_id = ?", req.UserID)
	db.Count(&total)
	db = db.Offset((req.Page - 1) * req.PageSize).Limit(req.PageSize)
	if err := db.Find(&friends).Error; err != nil {
		return GetFriendListResp{}, err
	}

	// 批量查好友用户信息
	friendIDs := make([]string, 0, len(friends))
	for _, f := range friends {
		friendIDs = append(friendIDs, f.FriendUserID)
	}
	var users []model.User
	if len(friendIDs) > 0 {
		if err := s.db.WithContext(ctx).Where("user_id IN ?", friendIDs).Find(&users).Error; err != nil {
			return GetFriendListResp{}, err
		}
	}
	userMap := make(map[string]model.User)
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
		Page:     req.Page,
		PageSize: req.PageSize,
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

// 删除好友
func (s *FriendService) DeleteFriend(ctx context.Context, req DeleteFriendReq) error {
	if err := s.db.WithContext(ctx).Where("owner_user_id = ? AND friend_user_id = ?", req.OwnerUserID, req.FriendUserID).Delete(&model.Friend{}).Error; err != nil {
		return err
	}
	if err := s.db.WithContext(ctx).Where("owner_user_id = ? AND friend_user_id = ?", req.FriendUserID, req.OwnerUserID).Delete(&model.Friend{}).Error; err != nil {
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
	fromIDs := make([]string, 0, len(frs))
	for _, fr := range frs {
		fromIDs = append(fromIDs, fr.FromUserID)
	}
	var users []model.User
	if len(fromIDs) > 0 {
		if err := s.db.WithContext(ctx).Where("user_id IN ?", fromIDs).Find(&users).Error; err != nil {
			return GetPaginationFriendApplyListResp{}, err
		}
	}
	userMap := make(map[string]model.User)
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
	toIDs := make([]string, 0, len(frs))
	for _, fr := range frs {
		toIDs = append(toIDs, fr.ToUserID)
	}
	var users []model.User
	if len(toIDs) > 0 {
		if err := s.db.WithContext(ctx).Find(&users, "user_id IN ?", toIDs).Error; err != nil {
			return GetSelfFriendApplyListResp{}, err
		}
	}
	userMap := make(map[string]model.User)
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
func (s *FriendService) AddBlack(ctx context.Context) {

}

// 移除黑名单
func (s *FriendService) RemoveBlack(ctx context.Context) {

}

// 获取黑名单列表
func (s *FriendService) GetPaginationBlacks(ctx context.Context) {

}
