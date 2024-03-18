package chat

import (
	"context"

	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/domain/chat/dto"
	"github.com/user2410/rrms-backend/internal/domain/chat/model"
	"github.com/user2410/rrms-backend/internal/domain/chat/repo"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
)

type Service interface {
	CreateMsgroup(userId uuid.UUID, data *dto.CreateMsgGroup) (*model.MsgGroup, error)
	CreateMessage(groupId int64, data *dto.IncomingCreateMessageEvent) (*model.Message, error)
	GetMessagesOfGroup(groupId int64, offset, limit int32) ([]model.Message, error)
	DeleteMessage(groupId int64, data *dto.IncomingDeleteMessageEvent) (int, error)
	CheckGroupMembership(userId uuid.UUID, groupId int64) (bool, error)
}

type service struct {
	repo repo.Repo
}

func NewService(repo repo.Repo) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) CreateMsgroup(userId uuid.UUID, data *dto.CreateMsgGroup) (*model.MsgGroup, error) {
	return s.repo.CreateMsgroup(context.Background(), userId, data)
}

func (s *service) CreateMessage(groupId int64, data *dto.IncomingCreateMessageEvent) (*model.Message, error) {
	return s.repo.CreateMessage(context.Background(), groupId, data)
}

func (s *service) DeleteMessage(groupId int64, data *dto.IncomingDeleteMessageEvent) (int, error) {
	return s.repo.UpdateMessage(context.Background(), data.DeletedBy, &database.UpdateMessageParams{
		ID:       data.MessageId,
		Status:   database.MESSAGESTATUSDELETED,
		Content:  "",
		Type:     database.MESSAGETYPETEXT,
		FromUser: data.DeletedBy,
		GroupID:  groupId,
	})
}

func (s *service) CheckGroupMembership(userId uuid.UUID, groupId int64) (bool, error) {
	return s.repo.CheckGroupMembership(context.Background(), userId, groupId)
}

func (s *service) GetMessagesOfGroup(groupId int64, offset, limit int32) ([]model.Message, error) {
	return s.repo.GetMessagesOfGroup(context.Background(), groupId, offset, limit)
}
