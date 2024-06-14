package service

import (
	"context"
	"errors"
	"fmt"

	chat_dto "github.com/user2410/rrms-backend/internal/domain/chat/dto"
	chat_model "github.com/user2410/rrms-backend/internal/domain/chat/model"

	"github.com/google/uuid"
)

var (
	ErrAnonymousApplicant = errors.New("anonymous applicant")
)

func GetResourceName(aid int64) string {
	return fmt.Sprintf("[APPLICATION_%d]", aid)
}

func (s *service) CreateApplicationMsgGroup(aid int64, userId uuid.UUID) (*chat_model.MsgGroup, error) {
	a, err := s.applicationRepo.GetApplicationById(context.Background(), aid)
	if err != nil {
		return nil, err
	}
	if a.CreatorID == uuid.Nil {
		return nil, ErrAnonymousApplicant
	}

	return s.chatRepo.CreateMsgroup(context.Background(), userId, &chat_dto.CreateMsgGroup{
		Name: GetResourceName(aid),
		Members: []chat_dto.CreateMsgGroupMember{
			{UserId: userId},
			{UserId: a.CreatorID},
		},
	})
}

func (s *service) GetApplicationMsgGroup(aid int64, userId uuid.UUID) (*chat_model.MsgGroupExtended, error) {
	return s.chatRepo.GetMsgGroupByName(context.Background(), userId, GetResourceName(aid))
}
