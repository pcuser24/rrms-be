package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	application_model "github.com/user2410/rrms-backend/internal/domain/application/model"
	auth_model "github.com/user2410/rrms-backend/internal/domain/auth/model"
	misc_dto "github.com/user2410/rrms-backend/internal/domain/misc/dto"
	misc_service "github.com/user2410/rrms-backend/internal/domain/misc/service"
	property_model "github.com/user2410/rrms-backend/internal/domain/property/model"
	rental_model "github.com/user2410/rrms-backend/internal/domain/rental/model"
	unit_model "github.com/user2410/rrms-backend/internal/domain/unit/model"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	html_util "github.com/user2410/rrms-backend/internal/utils/template/html"
	text_util "github.com/user2410/rrms-backend/internal/utils/template/text"
	"github.com/user2410/rrms-backend/pkg/ds/set"
)

func (s *service) GetRentalByApplicationId(aid int64) (rental_model.RentalModel, error) {
	return s.domainRepo.ApplicationRepo.GetRentalByApplicationId(context.Background(), aid)
}

func (s *service) _sendNotification(
	am *application_model.ApplicationModel,
	title string,
	ownerEmailContent, ownerPushContent string,
	tenantEmailContent, tenantPushContent string,
	notificationType misc_service.NOTIFICATIONTYPE,
) error {
	// get all managers
	managers, err := s.domainRepo.PropertyRepo.GetPropertyManagers(context.Background(), am.PropertyID)
	if err != nil {
		return err
	}
	managerIds := set.NewSet[uuid.UUID]()
	for _, m := range managers {
		managerIds.Add(m.ManagerID)
	}

	// get all neccessary users info for rendering content
	var users []auth_model.UserModel
	if am.CreatorID != uuid.Nil {
		userIds := managerIds.Clone()
		userIds.Add(am.CreatorID)

		users, err = s.domainRepo.AuthRepo.GetUsersByIds(context.Background(), userIds.ToSlice(), []string{"first_name", "last_name", "email", "phone"})
		if err != nil {
			return err
		}
	} else {
		users, err = s.domainRepo.AuthRepo.GetUsersByIds(context.Background(), managerIds.ToSlice(), []string{"first_name", "last_name", "email", "phone"})
		if err != nil {
			return err
		}
	}

	// prepare and send notifications to managers
	emailcn := misc_dto.CreateNotification{
		Title:   string(title),
		Content: string(ownerEmailContent),
		Data: map[string]interface{}{
			"notificationType": notificationType,
			"applicationId":    am.ID,
		},
	}
	// send email notification
	for mid := range managerIds {
		// get user
		var user auth_model.UserModel
		for _, u := range users {
			if u.ID == mid {
				user = u
				break
			}
		}
		emailcn.Targets = append(emailcn.Targets, misc_dto.CreateNotificationTarget{
			UserId: mid,
			Emails: []string{user.Email},
			Tokens: nil,
		})
	}
	err = s.miscService.SendNotification(&emailcn)
	if err != nil {
		return err
	}
	// send push notifications
	pushcn := misc_dto.CreateNotification{
		Title:   string(title),
		Content: string(ownerPushContent),
		Data: map[string]interface{}{
			"notificationType": notificationType,
			"applicationId":    am.ID,
		},
	}
	for mid := range managerIds {
		// get device tokens
		devices, err := s.miscService.GetNotificationDevice(mid, uuid.Nil, "", "")
		if err != nil {
			return err
		}
		tokens := make([]string, len(devices))
		for i, d := range devices {
			tokens[i] = d.Token
		}
		pushcn.Targets = append(pushcn.Targets, misc_dto.CreateNotificationTarget{
			UserId: mid,
			Emails: nil,
			Tokens: tokens,
		})
	}
	err = s.miscService.SendNotification(&pushcn)
	if err != nil {
		return err
	}

	// prepare and send notifications to tenant
	emailcn = misc_dto.CreateNotification{
		Title:   string(title),
		Content: string(tenantEmailContent),
		Data: map[string]interface{}{
			"notificationType": notificationType,
			"applicationId":    am.ID,
		},
	}
	if am.CreatorID != uuid.Nil {
		var user auth_model.UserModel
		for _, u := range users {
			if u.ID == am.CreatorID {
				user = u
				break
			}
		}
		emailcn.Targets = append(emailcn.Targets, misc_dto.CreateNotificationTarget{
			UserId: am.CreatorID,
			Emails: []string{user.Email, am.Email},
		})
	} else {
		emailcn.Targets = append(emailcn.Targets, misc_dto.CreateNotificationTarget{
			UserId: am.CreatorID,
			Emails: []string{am.Email},
		})
	}
	// send email notification
	err = s.miscService.SendNotification(&emailcn)
	if err != nil {
		return err
	}
	// send push notifications
	if am.CreatorID != uuid.Nil {
		pushcn := misc_dto.CreateNotification{
			Title:   string(title),
			Content: string(tenantPushContent),
			Data: map[string]interface{}{
				"notificationType": notificationType,
				"applicationId":    am.ID,
			},
		}
		devices, err := s.miscService.GetNotificationDevice(am.CreatorID, uuid.Nil, "", "")
		if err != nil {
			return err
		}
		tokens := make([]string, len(devices))
		for i, d := range devices {
			tokens[i] = d.Token
		}
		pushcn.Targets = append(pushcn.Targets, misc_dto.CreateNotificationTarget{
			UserId: am.CreatorID,
			Tokens: tokens,
		})
		err = s.miscService.SendNotification(&pushcn)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *service) sendNotificationOnNewApplication(am *application_model.ApplicationModel) error {
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
	data.Property, err = s.domainRepo.PropertyRepo.GetPropertyById(context.Background(), am.PropertyID)
	if err != nil {
		return err
	}
	data.Unit, err = s.domainRepo.UnitRepo.GetUnitById(context.Background(), am.UnitID)
	if err != nil {
		return err
	}

	title, err := text_util.RenderText(
		struct{}{},
		fmt.Sprintf("%s/title/create_application.txt", basePath),
		nil,
	)
	if err != nil {
		return err
	}
	ownerEmailContent, err := html_util.RenderHtml(
		data,
		fmt.Sprintf("%s/email/create_application_manager.html", basePath),
		nil,
	)
	if err != nil {
		return err
	}
	ownerPushContent, err := text_util.RenderText(
		data,
		fmt.Sprintf("%s/push/create_application_manager.txt", basePath),
		nil,
	)
	if err != nil {
		return err
	}
	tenantEmailContent, err := html_util.RenderHtml(
		data,
		fmt.Sprintf("%s/email/create_application_tenant.html", basePath),
		nil,
	)
	if err != nil {
		return err
	}
	tenantPushContent, err := text_util.RenderText(
		data,
		fmt.Sprintf("%s/push/create_application_tenant.txt", basePath),
		nil,
	)
	if err != nil {
		return err
	}
	return s._sendNotification(
		am, string(title),
		string(ownerEmailContent), string(ownerPushContent),
		string(tenantEmailContent), string(tenantPushContent),
		misc_service.NOTIFICATIONTYPE_CREATEAPPLICATION,
	)
}

func (s *service) sendNotificationOnUpdateApplication(am *application_model.ApplicationModel, status database.APPLICATIONSTATUS) error {
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
	data.Property, err = s.domainRepo.PropertyRepo.GetPropertyById(context.Background(), am.PropertyID)
	if err != nil {
		return err
	}
	data.Unit, err = s.domainRepo.UnitRepo.GetUnitById(context.Background(), am.UnitID)
	if err != nil {
		return err
	}

	title, err := text_util.RenderText(
		data,
		fmt.Sprintf("%s/title/update_application.txt", basePath),
		nil,
	)
	if err != nil {
		return err
	}
	ownerEmailContent, err := html_util.RenderHtml(
		data,
		fmt.Sprintf("%s/email/update_application_manager.html", basePath),
		nil,
	)
	if err != nil {
		return err
	}
	ownerPushContent, err := text_util.RenderText(
		data,
		fmt.Sprintf("%s/push/update_application_manager.txt", basePath),
		nil,
	)
	if err != nil {
		return err
	}
	tenantEmailContent, err := html_util.RenderHtml(
		data,
		fmt.Sprintf("%s/email/update_application_tenant.html", basePath),
		nil,
	)
	if err != nil {
		return err
	}
	tenantPushContent, err := text_util.RenderText(
		data,
		fmt.Sprintf("%s/push/update_application_tenant.txt", basePath),
		nil,
	)
	if err != nil {
		return err
	}
	return s._sendNotification(
		am, string(title),
		string(ownerEmailContent), string(ownerPushContent),
		string(tenantEmailContent), string(tenantPushContent),
		misc_service.NOTIFICATIONTYPE_UPDATEAPPLICATION,
	)
}
