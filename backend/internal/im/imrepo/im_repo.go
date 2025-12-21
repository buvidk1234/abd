package imrepo

import (
	"backend/internal/model"
	"backend/internal/pkg/cache/cachekey"
	rds "backend/internal/pkg/cache/redis"
	"backend/internal/pkg/constant"
	"backend/internal/service"
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const msgCacheTimeout = time.Hour * 24

type ImRepo struct {
	db         *gorm.DB
	rdb        *redis.Client
	seqConv    *rds.SeqConversationCacheRedis
	msgService *service.MessageService
}

func NewImRepo(db *gorm.DB, rdb *redis.Client) *ImRepo {
	return &ImRepo{
		db:         db,
		rdb:        rdb,
		seqConv:    rds.NewSeqConversationCacheRedis(db, rdb),
		msgService: service.NewMessageService(db),
	}
}

func (r *ImRepo) BatchStoreMsgToRedis(ctx context.Context, conversationID string, msgs []*model.Message) (isNewConversation bool, err error) {
	len := int64(len(msgs))
	lastSeq, err := r.seqConv.Malloc(ctx, conversationID, len)
	if err != nil {
		return false, err
	}
	isNewConversation = lastSeq == 0
	for _, msg := range msgs {
		lastSeq++
		msg.Seq = lastSeq
	}

	pipe := r.rdb.Pipeline()

	for _, msg := range msgs {
		data, err := json.Marshal(msg)
		if err != nil {
			return isNewConversation, err
		}
		// 将 Set 命令加入管道，注意这里不会立即执行
		pipe.Set(ctx, cachekey.GetMsgCacheKey(conversationID, msg.Seq), string(data), msgCacheTimeout)
	}

	// 一次性执行管道中的所有命令
	_, err = pipe.Exec(ctx)
	return isNewConversation, err

}
func (r *ImRepo) BatchStoreMsgToDB(ctx context.Context, msgs []*model.Message) error {
	if len(msgs) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}}, // clientID
		DoUpdates: clause.AssignmentColumns([]string{"content", "send_time"}),
	}).Create(msgs).Error
}

func (r *ImRepo) BatchGetMsg(ctx context.Context, key string, start, end int64) ([]string, error) {
	return nil, nil
}

func (r *ImRepo) CreateConversations(ctx context.Context, req service.InitConversationReq) error {
	err := r.msgService.InitConversation(ctx, req)
	return err
}

// InvalidateConversationIDsCache 删除用户的会话ID列表缓存，确保下次拉取从DB重建。
func (r *ImRepo) InvalidateConversationIDsCache(ctx context.Context, req service.InitConversationReq) {
	switch req.ConvType {
	case constant.SingleChatType:
		owners := []int64{req.SenderID, req.TargetID}
		for _, owner := range owners {
			_ = r.rdb.Del(ctx, cachekey.GetConversationIDsKey(strconv.FormatInt(owner, 10))).Err()
		}
	case constant.GroupChatType:
		// TODO: 获取群成员并批量删除缓存，先留空以免阻塞。
	default:
	}
}
