package repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/domain/chat/dto"
	"github.com/user2410/rrms-backend/internal/domain/chat/model"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
)

type Repo interface {
	CreateMsgroup(ctx context.Context, userId uuid.UUID, data *dto.CreateMsgGroup) (*model.MsgGroup, error)
	CreateMessage(ctx context.Context, groupId int64, data *dto.IncomingCreateMessageEvent) (*model.Message, error)
	UpdateMessage(ctx context.Context, userId uuid.UUID, data *database.UpdateMessageParams) (int, error)
	CheckGroupMembership(ctx context.Context, userId uuid.UUID, groupId int64) (bool, error)
	GetMsgGroupByName(ctx context.Context, userId uuid.UUID, name string) (*model.MsgGroupExtended, error)
	GetMessagesOfGroup(ctx context.Context, groupId int64, offset, limit int32) ([]model.Message, error)
}

type repo struct {
	dao database.DAO
}

func NewRepo(dao database.DAO) Repo {
	return &repo{
		dao: dao,
	}
}

func (r *repo) CreateMsgroup(ctx context.Context, userId uuid.UUID, data *dto.CreateMsgGroup) (*model.MsgGroup, error) {
	res, err := r.dao.CreateMsgGroup(ctx, database.CreateMsgGroupParams{
		CreatedBy: userId,
		Name:      data.Name,
	})
	if err != nil {
		return nil, err
	}
	gm := model.MsgGroup{
		GroupID:   res.GroupID,
		Name:      data.Name,
		CreatedAt: res.CreatedAt,
		CreatedBy: userId,
	}
	for _, m := range data.Members {
		mdb, err := r.dao.CreateMsgGroupMember(ctx, database.CreateMsgGroupMemberParams{
			GroupID: gm.GroupID,
			UserID:  m.UserId,
		})
		if err != nil {
			_ = r.dao.DeleteMsgGroup(ctx, gm.GroupID)
			return nil, err
		}
		gm.Members = append(gm.Members, model.MsgGroupMember(mdb))
	}
	return &gm, nil
}

func (r *repo) CreateMessage(ctx context.Context, groupId int64, data *dto.IncomingCreateMessageEvent) (*model.Message, error) {
	res, err := r.dao.CreateMessage(ctx, database.CreateMessageParams{
		FromUser: data.From,
		Content:  data.Content,
		GroupID:  groupId,
	})
	if err != nil {
		return nil, err
	}

	var m model.Message = model.Message(res)
	return &m, nil
}

func (r *repo) UpdateMessage(ctx context.Context, userId uuid.UUID, data *database.UpdateMessageParams) (int, error) {
	res, err := r.dao.UpdateMessage(ctx, *data)
	if err != nil {
		return 0, err
	}
	return len(res), nil
}

func (r *repo) CheckGroupMembership(ctx context.Context, userId uuid.UUID, groupId int64) (bool, error) {
	return r.dao.CheckMsgGroupMembership(ctx, database.CheckMsgGroupMembershipParams{
		UserID:  userId,
		GroupID: groupId,
	})
}

func (r *repo) GetMsgGroupByName(ctx context.Context, userId uuid.UUID, name string) (*model.MsgGroupExtended, error) {
	res, err := r.dao.GetMsgGroupByName(ctx, database.GetMsgGroupByNameParams{
		Name:   name,
		UserID: userId,
	})
	if err != nil {
		return nil, err
	}

	gm := model.MsgGroupExtended{
		GroupID:   res.GroupID,
		Name:      res.Name,
		CreatedAt: res.CreatedAt,
		CreatedBy: res.CreatedBy,
	}

	members, err := r.dao.GetMsgGroupMembers(ctx, gm.GroupID)
	if err != nil {
		return nil, err
	}
	for _, m := range members {
		gm.Members = append(gm.Members, model.MsgGroupMemberExtended(m))
	}

	return &gm, nil
}

func (r *repo) GetMessagesOfGroup(ctx context.Context, groupId int64, offset, limit int32) ([]model.Message, error) {
	res, err := r.dao.GetMessagesOfGroup(ctx, database.GetMessagesOfGroupParams{
		GroupID: groupId,
		Offset:  offset,
		Limit:   limit,
	})
	if err != nil {
		return nil, err
	}

	var msgs []model.Message
	for _, m := range res {
		msgs = append(msgs, model.Message(m))
	}
	return msgs, nil
}
