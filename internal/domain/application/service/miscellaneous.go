package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	application_model "github.com/user2410/rrms-backend/internal/domain/application/model"
	misc_dto "github.com/user2410/rrms-backend/internal/domain/misc/dto"
	misc_service "github.com/user2410/rrms-backend/internal/domain/misc/service"
	property_model "github.com/user2410/rrms-backend/internal/domain/property/model"
	rental_model "github.com/user2410/rrms-backend/internal/domain/rental/model"
	unit_model "github.com/user2410/rrms-backend/internal/domain/unit/model"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	html_util "github.com/user2410/rrms-backend/internal/utils/template/html"
	text_util "github.com/user2410/rrms-backend/internal/utils/template/text"
)

func (s *service) GetRentalByApplicationId(aid int64) (rental_model.RentalModel, error) {
	return s.domainRepo.ApplicationRepo.GetRentalByApplicationId(context.Background(), aid)
}

func (s *service) SendNotificationOnNewApplication(am *application_model.ApplicationModel) error {
	var (
		data = struct {
			FESite      string
			Application *application_model.ApplicationModel
			Property    *property_model.PropertyModel
			Unit        *unit_model.UnitModel
		}{
			FESite:      s.feSite,
			Application: am,
		}
		err error
	)
	{
		ps, err := s.domainRepo.PropertyRepo.GetPropertiesByIds(context.Background(), []uuid.UUID{am.PropertyID}, []string{"name"})
		if err != nil {
			return err
		}
		if len(ps) == 0 {
			return database.ErrRecordNotFound
		}
		data.Property = &ps[0]

		us, err := s.domainRepo.UnitRepo.GetUnitsByIds(context.Background(), []uuid.UUID{am.UnitID}, []string{"name"})
		if err != nil {
			return err
		}
		if len(us) == 0 {
			return database.ErrRecordNotFound
		}
		data.Unit = &us[0]
	}

	title, err := text_util.RenderText(
		struct{}{},
		fmt.Sprintf("%s/title/create_application.txt", basePath),
		nil,
	)
	if err != nil {
		return err
	}
	emailContent, err := html_util.RenderHtml(
		data,
		fmt.Sprintf("%s/email/create_application_manager.html", basePath),
		nil,
	)
	if err != nil {
		return err
	}
	pushContent, err := text_util.RenderText(
		data,
		fmt.Sprintf("%s/push/create_application_manager.txt", basePath),
		nil,
	)
	if err != nil {
		return err
	}
	targets, err := s.miscService.GetNotificationManagersTargets(am.PropertyID)
	if err != nil {
		return err
	}
	cn := misc_dto.CreateNotification{
		Title:   string(title),
		Content: string(emailContent),
		Data: map[string]interface{}{
			"notificationType": misc_service.NOTIFICATIONTYPE_CREATEAPPLICATION,
			"applicationId":    am.ID,
		},
		Targets: func() []misc_dto.CreateNotificationTarget {
			var ts []misc_dto.CreateNotificationTarget
			for _, t := range targets {
				ts = append(ts, misc_dto.CreateNotificationTarget{
					UserId: t.UserId,
					Emails: t.Emails,
				})
			}
			return ts
		}(),
	}
	err = s.miscService.SendNotification(&cn)
	if err != nil {
		return err
	}

	cn.Content = string(pushContent)
	cn.Targets = func() []misc_dto.CreateNotificationTarget {
		var ts []misc_dto.CreateNotificationTarget
		for _, t := range targets {
			ts = append(ts, misc_dto.CreateNotificationTarget{
				UserId: t.UserId,
				Tokens: t.Tokens,
			})
		}
		return ts
	}()
	err = s.miscService.SendNotification(&cn)

	return err
}

func (s *service) SendNotificationOnUpdateApplication(am *application_model.ApplicationModel, status database.APPLICATIONSTATUS) error {
	var (
		data = struct {
			Application *application_model.ApplicationModel
			Property    *property_model.PropertyModel
			Unit        *unit_model.UnitModel
			Status      database.APPLICATIONSTATUS
			FESite      string
		}{
			Application: am,
			Status:      status,
			FESite:      s.feSite,
		}
		err error
	)
	{
		ps, err := s.domainRepo.PropertyRepo.GetPropertiesByIds(context.Background(), []uuid.UUID{am.PropertyID}, []string{"name"})
		if err != nil {
			return err
		}
		if len(ps) == 0 {
			return database.ErrRecordNotFound
		}
		data.Property = &ps[0]

		us, err := s.domainRepo.UnitRepo.GetUnitsByIds(context.Background(), []uuid.UUID{am.UnitID}, []string{"name"})
		if err != nil {
			return err
		}
		if len(us) == 0 {
			return database.ErrRecordNotFound
		}
		data.Unit = &us[0]
	}

	title, err := text_util.RenderText(
		data,
		fmt.Sprintf("%s/title/update_application.txt", basePath),
		nil,
	)
	if err != nil {
		return err
	}
	emailContent, err := html_util.RenderHtml(
		data,
		fmt.Sprintf("%s/email/update_application_tenant.html", basePath),
		nil,
	)
	if err != nil {
		return err
	}
	pushContent, err := text_util.RenderText(
		data,
		fmt.Sprintf("%s/push/update_application_tenant.txt", basePath),
		nil,
	)
	if err != nil {
		return err
	}

	target, err := s.miscService.GetNotificationTenantTargets(am.CreatorID, am.Email)
	if err != nil {
		return err
	}

	cn := misc_dto.CreateNotification{
		Title:   string(title),
		Content: string(emailContent),
		Data: map[string]interface{}{
			"notificationType": misc_service.NOTIFICATIONTYPE_UPDATEAPPLICATION,
			"applicationId":    am.ID,
		},
		Targets: []misc_dto.CreateNotificationTarget{
			{
				UserId: am.CreatorID,
				Emails: target.Emails,
			},
		},
	}
	err = s.miscService.SendNotification(&cn)
	if err != nil {
		return err
	}

	cn.Content = string(pushContent)
	cn.Targets = []misc_dto.CreateNotificationTarget{
		{
			UserId: am.CreatorID,
			Tokens: target.Tokens,
		},
	}
	err = s.miscService.SendNotification(&cn)

	return err
}
