package service

import (
	"context"
	"fmt"
	"slices"

	"github.com/google/uuid"

	auth_model "github.com/user2410/rrms-backend/internal/domain/auth/model"
	property_model "github.com/user2410/rrms-backend/internal/domain/property/model"
	rental_model "github.com/user2410/rrms-backend/internal/domain/rental/model"
	unit_model "github.com/user2410/rrms-backend/internal/domain/unit/model"

	misc_dto "github.com/user2410/rrms-backend/internal/domain/misc/dto"
	rental_dto "github.com/user2410/rrms-backend/internal/domain/rental/dto"

	misc_service "github.com/user2410/rrms-backend/internal/domain/misc/service"

	"github.com/user2410/rrms-backend/internal/infrastructure/database"

	rental_util "github.com/user2410/rrms-backend/internal/domain/rental/utils"
	template_util "github.com/user2410/rrms-backend/internal/utils/template"
	html_util "github.com/user2410/rrms-backend/internal/utils/template/html"
	text_util "github.com/user2410/rrms-backend/internal/utils/template/text"
)

var (
	basePath = "internal/domain/rental/service/templates"
)

func (s *service) NotifyCreatePreRental(
	r *rental_model.RentalModel,
	secret string,
) error {
	// get target emails and push tokens
	var (
		targets    = make([]misc_dto.CreateNotificationTarget, 0)
		property   property_model.PropertyModel
		unit       unit_model.UnitModel
		accessLink string
		key        string
	)
	{
		ts, err := s.mService.GetNotificationManagersTargets(r.PropertyID)
		if err != nil {
			return err
		}
		targets = append(targets, ts...)
		t, err := s.mService.GetNotificationTenantTargets(r.TenantID, r.TenantEmail)
		if err != nil {
			return err
		}
		targets = append(targets, t)

		ps, err := s.domainRepo.PropertyRepo.GetPropertiesByIds(context.Background(), []uuid.UUID{r.PropertyID}, []string{"name"})
		if err != nil {
			return err
		}
		if len(ps) == 0 {
			return database.ErrRecordNotFound
		}
		property = ps[0]

		us, err := s.domainRepo.UnitRepo.GetUnitsByIds(context.Background(), []uuid.UUID{r.UnitID}, []string{"name"})
		if err != nil {
			return err
		}
		if len(us) == 0 {
			return database.ErrRecordNotFound
		}
		unit = us[0]

		if r.TenantID != uuid.Nil {
			accessLink = fmt.Sprintf("%s/manage/rentals/prerentals/prerental/%d", s.feSite, r.ID)
		} else {
			key, err = rental_util.CreatePreRentalKey(secret, r)
			if err != nil {
				return err
			}
			accessLink = fmt.Sprintf("%s/manage/rentals/prerentals/prerental/%d?key=%s", s.feSite, r.ID, key)
		}
	}

	data := struct {
		FESite     string
		Rental     *rental_model.RentalModel
		Property   property_model.PropertyModel
		Unit       unit_model.UnitModel
		AccessLink string
	}{
		FESite:     s.feSite,
		Property:   property,
		Unit:       unit,
		Rental:     r,
		AccessLink: accessLink,
	}

	title, err := text_util.RenderText(
		data,
		fmt.Sprintf("%s/title/create_prerental.txt", basePath),
		map[string]any{
			"Dereference": template_util.Dereference("-"),
		},
	)
	if err != nil {
		return err
	}
	emailContent, err := html_util.RenderHtml(
		data,
		fmt.Sprintf("%s/email/create_prerental.gohtml", basePath),
		map[string]any{
			"Dereference": template_util.Dereference("-"),
		},
	)
	if err != nil {
		return err
	}
	pushContent, err := text_util.RenderText(
		data,
		fmt.Sprintf("%s/push/create_prerental.txt", basePath),
		map[string]any{
			"Dereference": template_util.Dereference("-"),
		},
	)
	if err != nil {
		return err
	}

	cn := misc_dto.CreateNotification{
		Title:   string(title),
		Content: string(emailContent),
		Data: map[string]interface{}{
			"notificationType": misc_service.NOTIFICATIONTYPE_CREATEPRERENTAL,
			"prerentalId":      r.ID,
			"key":              key,
		},
		Targets: func() []misc_dto.CreateNotificationTarget {
			var ts []misc_dto.CreateNotificationTarget
			for _, t := range targets {
				if t.UserId == r.CreatorID {
					continue
				}
				ts = append(ts, misc_dto.CreateNotificationTarget{
					UserId: t.UserId,
					Emails: t.Emails,
					Tokens: []string{},
				})
			}
			return ts
		}(),
	}
	if err = s.mService.SendNotification(&cn); err != nil {
		return err
	}

	cn.Content = string(pushContent)
	cn.Targets = func() []misc_dto.CreateNotificationTarget {
		var ts []misc_dto.CreateNotificationTarget
		for _, t := range targets {
			if t.UserId == r.CreatorID {
				continue
			}
			ts = append(ts, misc_dto.CreateNotificationTarget{
				UserId: t.UserId,
				Emails: []string{},
				Tokens: t.Tokens,
			})
		}
		return ts
	}()
	if err = s.mService.SendNotification(&cn); err != nil {
		return err
	}

	return nil
}

func (s *service) NotifyUpdatePreRental(
	preRental *rental_model.PreRental,
	rental *rental_model.RentalModel,
	updateData *rental_dto.UpdatePreRental,
) error {
	// get target emails and push tokens
	var (
		targets  []misc_dto.CreateNotificationTarget = make([]misc_dto.CreateNotificationTarget, 0)
		property property_model.PropertyModel
		unit     unit_model.UnitModel
	)
	{
		ts, err := s.mService.GetNotificationManagersTargets(preRental.PropertyID)
		if err != nil {
			return err
		}
		targets = append(targets, ts...)

		ps, err := s.domainRepo.PropertyRepo.GetPropertiesByIds(context.Background(), []uuid.UUID{preRental.PropertyID}, []string{"name"})
		if err != nil {
			return err
		}
		if len(ps) == 0 {
			return database.ErrRecordNotFound
		}
		property = ps[0]

		us, err := s.domainRepo.UnitRepo.GetUnitsByIds(context.Background(), []uuid.UUID{preRental.UnitID}, []string{"name"})
		if err != nil {
			return err
		}
		if len(us) == 0 {
			return database.ErrRecordNotFound
		}
		unit = us[0]
	}

	data := struct {
		FESite     string
		PreRental  *rental_model.PreRental
		Rental     *rental_model.RentalModel
		Property   property_model.PropertyModel
		Unit       unit_model.UnitModel
		UpdateData *rental_dto.UpdatePreRental
	}{
		FESite:     s.feSite,
		PreRental:  preRental,
		Rental:     rental,
		Property:   property,
		Unit:       unit,
		UpdateData: updateData,
	}

	title, err := text_util.RenderText(
		data,
		fmt.Sprintf("%s/title/update_prerental.txt", basePath),
		map[string]any{
			"Dereference": template_util.Dereference("-"),
		},
	)
	if err != nil {
		return err
	}

	emailContent, err := html_util.RenderHtml(
		data,
		fmt.Sprintf("%s/email/update_prerental.gohtml", basePath),
		map[string]any{
			"Dereference": template_util.Dereference("-"),
		},
	)
	if err != nil {
		return err
	}

	pushContent, err := text_util.RenderText(
		data,
		fmt.Sprintf("%s/push/update_prerental.txt", basePath),
		map[string]any{
			"Dereference": template_util.Dereference("-"),
		},
	)
	if err != nil {
		return err
	}

	cn := misc_dto.CreateNotification{
		Title:   string(title),
		Content: string(emailContent),
		Data: map[string]interface{}{
			"notificationType": misc_service.NOTIFICATIONTYPE_UPDATEPRERENTAL,
			"preRentalId":      preRental.ID,
		},
		Targets: func() []misc_dto.CreateNotificationTarget {
			var ts []misc_dto.CreateNotificationTarget
			for _, t := range targets {
				ts = append(ts, misc_dto.CreateNotificationTarget{
					UserId: t.UserId,
					Emails: t.Emails,
					Tokens: []string{},
				})
			}
			return ts
		}(),
	}
	if rental != nil {
		cn.Data["rentalId"] = rental.ID
	}
	if err = s.mService.SendNotification(&cn); err != nil {
		return err
	}

	cn.Content = string(pushContent)
	cn.Targets = func() []misc_dto.CreateNotificationTarget {
		var ts []misc_dto.CreateNotificationTarget
		for _, t := range targets {
			ts = append(ts, misc_dto.CreateNotificationTarget{
				UserId: t.UserId,
				Emails: []string{},
				Tokens: t.Tokens,
			})
		}
		return ts
	}()
	if err = s.mService.SendNotification(&cn); err != nil {
		return err
	}

	return nil
}

var (
	mapStatusToNotificationFile = map[database.RENTALPAYMENTSTATUS]struct {
		Title string
		Email string
		Push  string
	}{
		database.RENTALPAYMENTSTATUSPLAN:          {fmt.Sprintf("%s/title/update_rpstatusplan.txt", basePath), fmt.Sprintf("%s/email/update_rpstatusplan.gohtml", basePath), fmt.Sprintf("%s/push/update_rpstatusplan.txt", basePath)},
		database.RENTALPAYMENTSTATUSISSUED:        {fmt.Sprintf("%s/title/update_rpstatusissued.txt", basePath), fmt.Sprintf("%s/email/update_rpstatusissued.gohtml", basePath), fmt.Sprintf("%s/push/update_rpstatusissued.txt", basePath)},
		database.RENTALPAYMENTSTATUSPENDING:       {fmt.Sprintf("%s/title/update_rpstatuspending.txt", basePath), fmt.Sprintf("%s/email/update_rpstatuspending.gohtml", basePath), fmt.Sprintf("%s/push/update_rpstatuspending.txt", basePath)},
		database.RENTALPAYMENTSTATUSREQUEST2PAY:   {fmt.Sprintf("%s/title/update_rpstatusrequest2pay.txt", basePath), fmt.Sprintf("%s/email/update_rpstatusrequest2pay.gohtml", basePath), fmt.Sprintf("%s/push/update_rpstatusrequest2pay.txt", basePath)},
		database.RENTALPAYMENTSTATUSPARTIALLYPAID: {fmt.Sprintf("%s/title/update_rpstatuspending.txt", basePath), fmt.Sprintf("%s/email/update_rpstatuspending.gohtml", basePath), fmt.Sprintf("%s/push/update_rpstatuspending.txt", basePath)},
		database.RENTALPAYMENTSTATUSPAYFINE:       {fmt.Sprintf("%s/title/update_rpstatuspayfine.txt", basePath), fmt.Sprintf("%s/email/update_rpstatuspayfine.gohtml", basePath), fmt.Sprintf("%s/push/update_rpstatuspayfine.txt", basePath)},
	}
)

func (s *service) NotifyUpdatePayments(
	r *rental_model.RentalModel,
	rp *rental_model.RentalPayment, // old rental payment data before update
	u *rental_dto.UpdateRentalPayment,
) error {
	// get target emails and push tokens
	var (
		targets  []misc_dto.CreateNotificationTarget = make([]misc_dto.CreateNotificationTarget, 0)
		property property_model.PropertyModel
		unit     unit_model.UnitModel
	)
	if slices.Contains([]database.RENTALPAYMENTSTATUS{
		database.RENTALPAYMENTSTATUSPLAN,
		database.RENTALPAYMENTSTATUSREQUEST2PAY,
		database.RENTALPAYMENTSTATUSPAYFINE,
	}, rp.Status) {
		// notify tenant
		t, err := s.mService.GetNotificationTenantTargets(r.TenantID, r.TenantEmail)
		if err != nil {
			return err
		}
		targets = append(targets, t)
	} else if slices.Contains([]database.RENTALPAYMENTSTATUS{
		database.RENTALPAYMENTSTATUSISSUED,
		database.RENTALPAYMENTSTATUSPENDING,
		database.RENTALPAYMENTSTATUSPARTIALLYPAID,
	}, rp.Status) {
		// notify managers
		ts, err := s.mService.GetNotificationManagersTargets(r.PropertyID)
		if err != nil {
			return err
		}
		targets = append(targets, ts...)
	}

	{
		ps, err := s.domainRepo.PropertyRepo.GetPropertiesByIds(context.Background(), []uuid.UUID{r.PropertyID}, []string{"name"})
		if err != nil {
			return err
		}
		if len(ps) == 0 {
			return database.ErrRecordNotFound
		}
		property = ps[0]

		us, err := s.domainRepo.UnitRepo.GetUnitsByIds(context.Background(), []uuid.UUID{r.UnitID}, []string{"name"})
		if err != nil {
			return err
		}
		if len(us) == 0 {
			return database.ErrRecordNotFound
		}
		unit = us[0]
	}

	rpServiceName, _ := rental_util.GetServiceName(rp.Code, r.Services)
	data := struct {
		FESite         string
		Property       property_model.PropertyModel
		Unit           unit_model.UnitModel
		Rental         *rental_model.RentalModel
		Payment        *rental_model.RentalPayment
		PaymentService string
		UpdateData     *rental_dto.UpdateRentalPayment
	}{
		FESite:         s.feSite,
		Property:       property,
		Unit:           unit,
		Rental:         r,
		Payment:        rp,
		PaymentService: rpServiceName,
		UpdateData:     u,
	}

	title, err := text_util.RenderText(
		data,
		mapStatusToNotificationFile[rp.Status].Title,
		map[string]any{
			"Dereference": template_util.Dereference("-"),
		},
	)
	if err != nil {
		return err
	}
	emailContent, err := html_util.RenderHtml(
		data,
		mapStatusToNotificationFile[rp.Status].Email,
		map[string]any{
			"Dereference": template_util.Dereference("-"),
		},
	)
	if err != nil {
		return err
	}
	pushContent, err := text_util.RenderText(
		data,
		mapStatusToNotificationFile[rp.Status].Push,
		map[string]any{
			"Dereference": template_util.Dereference("0"),
		},
	)
	if err != nil {
		return err
	}

	cn := misc_dto.CreateNotification{
		Title:   string(title),
		Content: string(emailContent),
		Data: map[string]interface{}{
			"notificationType": misc_service.NOTIFICATIONTYPE_UPDATERENTALPAYMENT,
			"rentalId":         r.ID,
			"rentalPaymentId":  rp.ID,
		},
		Targets: func() []misc_dto.CreateNotificationTarget {
			var ts []misc_dto.CreateNotificationTarget
			for _, t := range targets {
				ts = append(ts, misc_dto.CreateNotificationTarget{
					UserId: t.UserId,
					Emails: t.Emails,
					Tokens: []string{},
				})
			}
			return ts
		}(),
	}
	if err = s.mService.SendNotification(&cn); err != nil {
		return err
	}

	cn.Content = string(pushContent)
	cn.Targets = func() []misc_dto.CreateNotificationTarget {
		var ts []misc_dto.CreateNotificationTarget
		for _, t := range targets {
			ts = append(ts, misc_dto.CreateNotificationTarget{
				UserId: t.UserId,
				Emails: []string{},
				Tokens: t.Tokens,
			})
		}
		return ts
	}()
	if err = s.mService.SendNotification(&cn); err != nil {
		return err
	}

	return nil
}

// NotifyCreateRentalPayment notifies tenant about newly issued rental payment
func (s *service) NotifyCreateRentalPayment(
	r *rental_model.RentalModel,
	rp *rental_model.RentalPayment,
) error {
	// get target emails and push tokens
	target, err := s.mService.GetNotificationTenantTargets(r.TenantID, r.TenantEmail)
	if err != nil {
		return err
	}

	var (
		property property_model.PropertyModel
		unit     unit_model.UnitModel
	)
	{
		ps, err := s.domainRepo.PropertyRepo.GetPropertiesByIds(context.Background(), []uuid.UUID{r.PropertyID}, []string{"name"})
		if err != nil {
			return err
		}
		if len(ps) == 0 {
			return database.ErrRecordNotFound
		}
		property = ps[0]

		us, err := s.domainRepo.UnitRepo.GetUnitsByIds(context.Background(), []uuid.UUID{r.UnitID}, []string{"name"})
		if err != nil {
			return err
		}
		if len(us) == 0 {
			return database.ErrRecordNotFound
		}
		unit = us[0]
	}

	rpServiceName, _ := rental_util.GetServiceName(rp.Code, r.Services)
	data := struct {
		FESite         string
		Property       property_model.PropertyModel
		Unit           unit_model.UnitModel
		Rental         *rental_model.RentalModel
		Payment        *rental_model.RentalPayment
		PaymentService string
	}{
		FESite:         s.feSite,
		Property:       property,
		Unit:           unit,
		Rental:         r,
		Payment:        rp,
		PaymentService: rpServiceName,
	}

	title, err := text_util.RenderText(
		data,
		fmt.Sprintf("%s/title/create_rentalpayment.txt", basePath),
		map[string]any{
			"Dereference": template_util.Dereference("-"),
		},
	)
	if err != nil {
		return err
	}
	emailContent, err := html_util.RenderHtml(
		data,
		fmt.Sprintf("%s/email/create_rentalpayment.gohtml", basePath),
		map[string]any{
			"Dereference": template_util.Dereference("-"),
		},
	)
	if err != nil {
		return err
	}
	pushContent, err := text_util.RenderText(
		data,
		fmt.Sprintf("%s/push/create_rentalpayment.txt", basePath),
		map[string]any{
			"Dereference": template_util.Dereference("-"),
		},
	)
	if err != nil {
		return err
	}

	cn := misc_dto.CreateNotification{
		Title:   string(title),
		Content: string(emailContent),
		Data: map[string]interface{}{
			"notificationType": misc_service.NOTIFICATIONTYPE_CREATERENTALPAYMENT,
			"rentalId":         r.ID,
			"rentalPaymentId":  rp.ID,
		},
		Targets: []misc_dto.CreateNotificationTarget{{
			UserId: target.UserId,
			Emails: target.Emails,
		}},
	}
	if err = s.mService.SendNotification(&cn); err != nil {
		return err
	}

	cn.Content = string(pushContent)
	cn.Targets = []misc_dto.CreateNotificationTarget{{
		UserId: target.UserId,
		Tokens: target.Tokens,
	}}
	return s.mService.SendNotification(&cn)
}

func (s *service) NotifyCreateContract(
	c *rental_model.ContractModel,
	r *rental_model.RentalModel,
) error {
	// get target emails and push tokens
	target, err := s.mService.GetNotificationTenantTargets(r.TenantID, r.TenantEmail)
	if err != nil {
		return err
	}

	var (
		creator  auth_model.UserModel
		property property_model.PropertyModel
		unit     unit_model.UnitModel
	)
	{
		cs, err := s.domainRepo.AuthRepo.GetUsersByIds(context.Background(), []uuid.UUID{c.CreatedBy}, []string{"first_name", "last_name", "email"})
		if err != nil {
			return err
		}
		if len(cs) == 0 {
			return database.ErrRecordNotFound
		}
		creator = cs[0]

		ps, err := s.domainRepo.PropertyRepo.GetPropertiesByIds(context.Background(), []uuid.UUID{r.PropertyID}, []string{"name"})
		if err != nil {
			return err
		}
		if len(ps) == 0 {
			return database.ErrRecordNotFound
		}
		property = ps[0]

		us, err := s.domainRepo.UnitRepo.GetUnitsByIds(context.Background(), []uuid.UUID{r.UnitID}, []string{"name"})
		if err != nil {
			return err
		}
		if len(us) == 0 {
			return database.ErrRecordNotFound
		}
		unit = us[0]
	}

	data := struct {
		FESite   string
		Contract *rental_model.ContractModel
		Rental   *rental_model.RentalModel
		Property property_model.PropertyModel
		Unit     unit_model.UnitModel
		Creator  auth_model.UserModel
	}{
		FESite:   s.feSite,
		Contract: c,
		Rental:   r,
		Property: property,
		Unit:     unit,
		Creator:  creator,
	}

	title, err := text_util.RenderText(
		data,
		fmt.Sprintf("%s/title/create_rentalcontract.txt", basePath),
		map[string]any{
			"Dereference": template_util.Dereference("-"),
		},
	)
	if err != nil {
		return err
	}
	emailContent, err := html_util.RenderHtml(
		data,
		fmt.Sprintf("%s/email/create_rentalcontract.gohtml", basePath),
		map[string]any{
			"Dereference": template_util.Dereference("-"),
		},
	)
	if err != nil {
		return err
	}
	pushContent, err := text_util.RenderText(
		data,
		fmt.Sprintf("%s/push/create_rentalcontract.txt", basePath),
		map[string]any{
			"Dereference": template_util.Dereference("-"),
		},
	)
	if err != nil {
		return err
	}

	cn := misc_dto.CreateNotification{
		Title:   string(title),
		Content: string(emailContent),
		Data: map[string]interface{}{
			"notificationType": misc_service.NOTIFICATIONTYPE_CREATECONTRACT,
			"rentalId":         r.ID,
			"contractId":       c.ID,
		},
		Targets: []misc_dto.CreateNotificationTarget{{
			UserId: target.UserId,
			Emails: target.Emails,
		}},
	}
	if err = s.mService.SendNotification(&cn); err != nil {
		return err
	}

	cn.Content = string(pushContent)
	cn.Targets = []misc_dto.CreateNotificationTarget{{
		UserId: target.UserId,
		Tokens: target.Tokens,
	}}
	if err = s.mService.SendNotification(&cn); err != nil {
		return err
	}

	return nil
}

func (s *service) NotifyUpdateContract(
	c *rental_model.ContractModel,
	r *rental_model.RentalModel,
	side string,
) error {
	var (
		targets []misc_dto.CreateNotificationTarget
		err     error
	)
	if side == "A" {
		target, err := s.mService.GetNotificationTenantTargets(r.TenantID, r.TenantEmail)
		if err != nil {
			return err
		}
		targets = []misc_dto.CreateNotificationTarget{target}
	} else if side == "B" {
		targets, err = s.mService.GetNotificationManagersTargets(r.PropertyID)
		if err != nil {
			return err
		}
	}

	var (
		updater  auth_model.UserModel
		property *property_model.PropertyModel
		unit     *unit_model.UnitModel
	)
	{
		us, err := s.domainRepo.AuthRepo.GetUsersByIds(context.Background(), []uuid.UUID{c.UpdatedBy}, []string{"first_name", "last_name", "email"})
		if err != nil {
			return err
		}
		if len(us) > 0 {
			updater = us[0]
		}

		ps, err := s.domainRepo.PropertyRepo.GetPropertiesByIds(context.Background(), []uuid.UUID{r.PropertyID}, []string{"name"})
		if err != nil {
			return err
		}
		if len(ps) > 0 {
			property = &ps[0]
		}

		units, err := s.domainRepo.UnitRepo.GetUnitsByIds(context.Background(), []uuid.UUID{r.UnitID}, []string{"name"})
		if err != nil {
			return err
		}
		if len(us) > 0 {
			unit = &units[0]
		}
	}

	data := struct {
		FESite   string
		Contract *rental_model.ContractModel
		Rental   *rental_model.RentalModel
		Property *property_model.PropertyModel
		Unit     *unit_model.UnitModel
		Updater  auth_model.UserModel
	}{
		FESite:   s.feSite,
		Contract: c,
		Rental:   r,
		Property: property,
		Unit:     unit,
		Updater:  updater,
	}

	title, err := text_util.RenderText(
		data,
		fmt.Sprintf("%s/title/update_contract.txt", basePath),
		map[string]any{
			"Dereference": template_util.Dereference("-"),
		},
	)
	if err != nil {
		return err
	}
	emailContent, err := html_util.RenderHtml(
		data,
		fmt.Sprintf("%s/email/update_contract.gohtml", basePath),
		map[string]any{
			"Dereference": template_util.Dereference("-"),
		},
	)
	if err != nil {
		return err
	}
	pushContent, err := text_util.RenderText(
		data,
		fmt.Sprintf("%s/push/update_contract.txt", basePath),
		map[string]any{
			"Dereference": template_util.Dereference("-"),
		},
	)
	if err != nil {
		return err
	}

	cn := misc_dto.CreateNotification{
		Title:   string(title),
		Content: string(emailContent),
		Data: map[string]interface{}{
			"notificationType": misc_service.NOTIFICATIONTYPE_UPDATECONTRACT,
			"rentalId":         r.ID,
			"contractId":       c.ID,
		},
		Targets: func() []misc_dto.CreateNotificationTarget {
			var ts []misc_dto.CreateNotificationTarget
			for _, t := range targets {
				ts = append(ts, misc_dto.CreateNotificationTarget{
					UserId: t.UserId,
					Emails: t.Emails,
					Tokens: []string{},
				})
			}
			return ts
		}(),
	}
	if err = s.mService.SendNotification(&cn); err != nil {
		return err
	}

	cn.Content = string(pushContent)
	cn.Targets = func() []misc_dto.CreateNotificationTarget {
		var ts []misc_dto.CreateNotificationTarget
		for _, t := range targets {
			ts = append(ts, misc_dto.CreateNotificationTarget{
				UserId: t.UserId,
				Emails: []string{},
				Tokens: t.Tokens,
			})
		}
		return ts
	}()
	if err = s.mService.SendNotification(&cn); err != nil {
		return err
	}

	return nil
}

func (s *service) NotifyCreateRentalComplaint(
	c *rental_model.RentalComplaint,
	r *rental_model.RentalModel,
) error {
	var (
		targets  []misc_dto.CreateNotificationTarget
		err      error
		property property_model.PropertyModel
		unit     unit_model.UnitModel
		side     string
	)
	{
		// get side of creator of the complaint
		side, err = s.domainRepo.RentalRepo.GetRentalSide(context.Background(), r.ID, c.CreatorID)
		if err != nil {
			return err
		}
		if side == "A" {
			target, err := s.mService.GetNotificationTenantTargets(r.TenantID, r.TenantEmail)
			if err != nil {
				return err
			}
			targets = []misc_dto.CreateNotificationTarget{target}
		} else if side == "B" {
			targets, err = s.mService.GetNotificationManagersTargets(r.PropertyID)
			if err != nil {
				return err
			}
		}

		ps, err := s.domainRepo.PropertyRepo.GetPropertiesByIds(context.Background(), []uuid.UUID{r.PropertyID}, []string{"name"})
		if err != nil {
			return err
		}
		if len(ps) == 0 {
			return database.ErrRecordNotFound
		}
		property = ps[0]

		units, err := s.domainRepo.UnitRepo.GetUnitsByIds(context.Background(), []uuid.UUID{r.UnitID}, []string{"name"})
		if err != nil {
			return err
		}
		if len(units) == 0 {
			return database.ErrRecordNotFound
		}
		unit = units[0]
	}

	data := struct {
		FESite    string
		Complaint *rental_model.RentalComplaint
		Rental    *rental_model.RentalModel
		Property  property_model.PropertyModel
		Unit      unit_model.UnitModel
		Side      string
	}{
		FESite:    s.feSite,
		Complaint: c,
		Rental:    r,
		Property:  property,
		Unit:      unit,
		Side:      side,
	}

	title, err := text_util.RenderText(
		data,
		fmt.Sprintf("%s/title/create_rentalcomplaint.txt", basePath),
		map[string]any{
			"Dereference": template_util.Dereference("-"),
		},
	)
	if err != nil {
		return err
	}
	emailContent, err := html_util.RenderHtml(
		data,
		fmt.Sprintf("%s/email/create_rentalcomplaint.gohtml", basePath),
		map[string]any{
			"Dereference": template_util.Dereference("-"),
		},
	)
	if err != nil {
		return err
	}
	pushContent, err := text_util.RenderText(
		data,
		fmt.Sprintf("%s/push/create_rentalcomplaint.txt", basePath),
		map[string]any{
			"Dereference": template_util.Dereference("-"),
		},
	)
	if err != nil {
		return err
	}

	cn := misc_dto.CreateNotification{
		Title:   string(title),
		Content: string(emailContent),
		Data: map[string]interface{}{
			"notificationType": misc_service.NOTIFICATIONTYPE_CREATERENTALCOMPLAINT,
			"rentalId":         r.ID,
			"complaintId":      c.ID,
		},
		Targets: func() []misc_dto.CreateNotificationTarget {
			var ts []misc_dto.CreateNotificationTarget
			for _, t := range targets {
				ts = append(ts, misc_dto.CreateNotificationTarget{
					UserId: t.UserId,
					Emails: t.Emails,
					Tokens: []string{},
				})
			}
			return ts
		}(),
	}
	if err = s.mService.SendNotification(&cn); err != nil {
		return err
	}

	cn.Content = string(pushContent)
	cn.Targets = func() []misc_dto.CreateNotificationTarget {
		var ts []misc_dto.CreateNotificationTarget
		for _, t := range targets {
			ts = append(ts, misc_dto.CreateNotificationTarget{
				UserId: t.UserId,
				Emails: []string{},
				Tokens: t.Tokens,
			})
		}
		return ts
	}()
	if err = s.mService.SendNotification(&cn); err != nil {
		return err
	}

	return nil
}

func (s *service) NotifyCreateComplaintReply(
	c *rental_model.RentalComplaint,
	cr *rental_model.RentalComplaintReply,
	r *rental_model.RentalModel,
) error {
	var (
		targets []misc_dto.CreateNotificationTarget
		err     error
		side    string
	)
	{
		// get side of creator of the complaint
		side, err = s.domainRepo.RentalRepo.GetRentalSide(context.Background(), r.ID, cr.ReplierID)
		if err != nil {
			return err
		}
		if side == "A" {
			target, err := s.mService.GetNotificationTenantTargets(r.TenantID, r.TenantEmail)
			if err != nil {
				return err
			}
			targets = []misc_dto.CreateNotificationTarget{target}
		} else if side == "B" {
			targets, err = s.mService.GetNotificationManagersTargets(r.PropertyID)
			if err != nil {
				return err
			}
		}
	}

	data := struct {
		FESite         string
		Complaint      *rental_model.RentalComplaint
		ComplaintReply *rental_model.RentalComplaintReply
		Rental         *rental_model.RentalModel
		Side           string
	}{
		FESite:         s.feSite,
		Complaint:      c,
		ComplaintReply: cr,
		Rental:         r,
		Side:           side,
	}

	title, err := text_util.RenderText(
		data,
		fmt.Sprintf("%s/title/create_complaintreply.txt", basePath),
		map[string]any{
			"Dereference": template_util.Dereference("-"),
		},
	)
	if err != nil {
		return err
	}
	emailContent, err := html_util.RenderHtml(
		data,
		fmt.Sprintf("%s/email/create_complaintreply.gohtml", basePath),
		map[string]any{
			"Dereference": template_util.Dereference("-"),
		},
	)
	if err != nil {
		return err
	}
	pushContent, err := text_util.RenderText(
		data,
		fmt.Sprintf("%s/push/create_complaintreply.txt", basePath),
		map[string]any{
			"Dereference": template_util.Dereference("-"),
		},
	)
	if err != nil {
		return err
	}

	cn := misc_dto.CreateNotification{
		Title:   string(title),
		Content: string(emailContent),
		Data: map[string]interface{}{
			"notificationType": misc_service.NOTIFICATIONTYPE_CREATERENTALCOMPLAINTREPLY,
			"rentalId":         r.ID,
			"complaintId":      c.ID,
		},
		Targets: func() []misc_dto.CreateNotificationTarget {
			var ts []misc_dto.CreateNotificationTarget
			for _, t := range targets {
				ts = append(ts, misc_dto.CreateNotificationTarget{
					UserId: t.UserId,
					Emails: t.Emails,
					Tokens: []string{},
				})
			}
			return ts
		}(),
	}
	if err = s.mService.SendNotification(&cn); err != nil {
		return err
	}

	cn.Content = string(pushContent)
	cn.Targets = func() []misc_dto.CreateNotificationTarget {
		var ts []misc_dto.CreateNotificationTarget
		for _, t := range targets {
			ts = append(ts, misc_dto.CreateNotificationTarget{
				UserId: t.UserId,
				Emails: []string{},
				Tokens: t.Tokens,
			})
		}
		return ts
	}()
	if err = s.mService.SendNotification(&cn); err != nil {
		return err
	}

	return nil
}

func (s *service) NotifyUpdateComplaintStatus(
	c *rental_model.RentalComplaint,
	r *rental_model.RentalModel,
	status database.RENTALCOMPLAINTSTATUS,
	updatedBy uuid.UUID,
) error {
	var (
		targets []misc_dto.CreateNotificationTarget
		err     error
	)
	{
		// get side of creator of the complaint
		side, err := s.domainRepo.RentalRepo.GetRentalSide(context.Background(), r.ID, updatedBy)
		if err != nil {
			return err
		}
		if side == "A" {
			target, err := s.mService.GetNotificationTenantTargets(r.TenantID, r.TenantEmail)
			if err != nil {
				return err
			}
			targets = []misc_dto.CreateNotificationTarget{target}
		} else if side == "B" {
			targets, err = s.mService.GetNotificationManagersTargets(r.PropertyID)
			if err != nil {
				return err
			}
		}
	}

	data := struct {
		FESite    string
		Complaint *rental_model.RentalComplaint
		Status    database.RENTALCOMPLAINTSTATUS
		Rental    *rental_model.RentalModel
	}{
		FESite:    s.feSite,
		Complaint: c,
		Rental:    r,
		Status:    status,
	}

	title, err := text_util.RenderText(
		data,
		fmt.Sprintf("%s/title/update_complaintstatus.txt", basePath),
		map[string]any{
			"Dereference": template_util.Dereference("-"),
		},
	)
	if err != nil {
		return err
	}
	emailContent, err := html_util.RenderHtml(
		data,
		fmt.Sprintf("%s/email/update_complaintstatus.gohtml", basePath),
		map[string]any{
			"Dereference": template_util.Dereference("-"),
		},
	)
	if err != nil {
		return err
	}
	pushContent, err := text_util.RenderText(
		data,
		fmt.Sprintf("%s/push/update_complaintstatus.txt", basePath),
		map[string]any{
			"Dereference": template_util.Dereference("-"),
		},
	)
	if err != nil {
		return err
	}

	cn := misc_dto.CreateNotification{
		Title:   string(title),
		Content: string(emailContent),
		Data: map[string]interface{}{
			"notificationType": misc_service.NOTIFICATIONTYPE_UPDATERENTALCOMPLAINTSTATUS,
			"rentalId":         r.ID,
			"complaintId":      c.ID,
		},
		Targets: func() []misc_dto.CreateNotificationTarget {
			var ts []misc_dto.CreateNotificationTarget
			for _, t := range targets {
				ts = append(ts, misc_dto.CreateNotificationTarget{
					UserId: t.UserId,
					Emails: t.Emails,
					Tokens: []string{},
				})
			}
			return ts
		}(),
	}
	if err = s.mService.SendNotification(&cn); err != nil {
		return err
	}

	cn.Content = string(pushContent)
	cn.Targets = func() []misc_dto.CreateNotificationTarget {
		var ts []misc_dto.CreateNotificationTarget
		for _, t := range targets {
			ts = append(ts, misc_dto.CreateNotificationTarget{
				UserId: t.UserId,
				Emails: []string{},
				Tokens: t.Tokens,
			})
		}
		return ts
	}()
	if err = s.mService.SendNotification(&cn); err != nil {
		return err
	}

	return nil
}
