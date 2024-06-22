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

	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	"github.com/user2410/rrms-backend/pkg/ds/set"

	rental_util "github.com/user2410/rrms-backend/internal/domain/rental/utils"
	template_util "github.com/user2410/rrms-backend/internal/utils/template"
	html_util "github.com/user2410/rrms-backend/internal/utils/template/html"
	text_util "github.com/user2410/rrms-backend/internal/utils/template/text"
)

var (
	basePath = "internal/domain/rental/service/templates"
)

func (s *service) _getNotificationManagersTargets(propertyID uuid.UUID) ([]misc_dto.CreateNotificationTarget, error) {
	managers, err := s.domainRepo.PropertyRepo.GetPropertyManagers(context.Background(), propertyID)
	if err != nil {
		return nil, err
	}
	managerIds := make([]uuid.UUID, 0)
	for _, m := range managers {
		managerIds = append(managerIds, m.ManagerID)
	}

	us, err := s.domainRepo.AuthRepo.GetUsersByIds(context.Background(), managerIds, []string{"email"})
	if err != nil {
		return nil, err
	}

	targets := make([]misc_dto.CreateNotificationTarget, 0)
	for _, u := range us {
		target := misc_dto.CreateNotificationTarget{
			UserId: u.ID,
			Emails: []string{u.Email},
		}
		ds, err := s.mService.GetNotificationDevice(u.ID, uuid.Nil, "", "")
		if err != nil {
			return nil, err
		}
		for _, d := range ds {
			target.Tokens = append(target.Tokens, d.Token)
		}
		targets = append(targets, target)
	}

	return targets, nil
}

func (s *service) _getNotificationTenantTargets(tenantID uuid.UUID, tenantEmail string) (misc_dto.CreateNotificationTarget, error) {
	emailTargets := set.NewSet[string]()
	pushTargets := set.NewSet[string]()
	emailTargets.Add(tenantEmail)
	if tenantID == uuid.Nil {
		ts, err := s.domainRepo.AuthRepo.GetUsersByIds(context.Background(), []uuid.UUID{tenantID}, []string{"email"})
		if err != nil {
			return misc_dto.CreateNotificationTarget{}, err
		}
		if len(ts) > 0 {
			emailTargets.Add(ts[0].Email)
			ds, err := s.mService.GetNotificationDevice(tenantID, uuid.Nil, "", "")
			if err != nil {
				return misc_dto.CreateNotificationTarget{}, err
			}
			for _, d := range ds {
				pushTargets.Add(d.Token)
			}
		}
	}
	return misc_dto.CreateNotificationTarget{
		UserId: tenantID,
		Emails: emailTargets.ToSlice(),
		Tokens: pushTargets.ToSlice(),
	}, nil
}

func (s *service) notifyCreatePreRental(
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
		ts, err := s._getNotificationManagersTargets(r.PropertyID)
		if err != nil {
			return err
		}
		targets = append(targets, ts...)
		t, err := s._getNotificationTenantTargets(r.TenantID, r.TenantEmail)
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

	if err = s.mService.SendNotification(&misc_dto.CreateNotification{
		Title:   string(title),
		Content: string(emailContent),
		Data: map[string]interface{}{
			"notificationType": "CREATE_PRERENTAL",
			"rentalId":         r.ID,
			"key":              key,
		},
		Targets: targets,
	}); err != nil {
		return err
	}

	if err = s.mService.SendNotification(&misc_dto.CreateNotification{
		Title:   string(title),
		Content: string(pushContent),
		Data: map[string]interface{}{
			"notificationType": "CREATE_PRERENTAL",
			"rentalId":         r.ID,
			"key":              key,
		},
		Targets: targets,
	}); err != nil {
		return err
	}

	return nil
}

func (s *service) notifyUpdatePreRental(
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
		ts, err := s._getNotificationManagersTargets(preRental.PropertyID)
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

	if err = s.mService.SendNotification(&misc_dto.CreateNotification{
		Title:   string(title),
		Content: string(emailContent),
		Data: map[string]interface{}{
			"notificationType": "UPDATE_PRERENTAL",
			"preRentalId":      preRental.ID,
		},
		Targets: targets,
	}); err != nil {
		return err
	}

	if err = s.mService.SendNotification(&misc_dto.CreateNotification{
		Title:   string(title),
		Content: string(pushContent),
		Data: map[string]interface{}{
			"notificationType": "UPDATE_PRERENTAL",
			"preRentalId":      preRental.ID,
		},
		Targets: targets,
	}); err != nil {
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

func (s *service) notifyUpdatePayments(
	r *rental_model.RentalModel,
	rp *rental_model.RentalPayment,
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
		t, err := s._getNotificationTenantTargets(r.TenantID, r.TenantEmail)
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
		ts, err := s._getNotificationManagersTargets(r.PropertyID)
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
			"Dereference": template_util.Dereference("-"),
		},
	)
	if err != nil {
		return err
	}

	if err = s.mService.SendNotification(&misc_dto.CreateNotification{
		Title:   string(title),
		Content: string(emailContent),
		Data: map[string]interface{}{
			"notificationType": "UPDATE_RENTALPAYMENT",
			"rentalId":         r.ID,
			"rentalPaymentId":  rp.ID,
		},
		Targets: targets,
	}); err != nil {
		return err
	}

	if err = s.mService.SendNotification(&misc_dto.CreateNotification{
		Title:   string(title),
		Content: string(pushContent),
		Data: map[string]interface{}{
			"notificationType": "UPDATE_RENTALPAYMENT",
			"rentalId":         r.ID,
			"rentalPaymentId":  rp.ID,
		},
		Targets: targets,
	}); err != nil {
		return err
	}

	return nil
}

// notifyCreateRentalPayment notifies tenant about newly issued rental payment
func (s *service) notifyCreateRentalPayment(
	r *rental_model.RentalModel,
	rp *rental_model.RentalPayment,
) error {
	// get target emails and push tokens
	target, err := s._getNotificationTenantTargets(r.TenantID, r.TenantEmail)
	if err != nil {
		return err
	}

	data := struct {
		FESite        string
		Rental        *rental_model.RentalModel
		RentalPayment *rental_model.RentalPayment
	}{
		FESite:        s.feSite,
		Rental:        r,
		RentalPayment: rp,
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

	if err = s.mService.SendNotification(&misc_dto.CreateNotification{
		Title:   string(title),
		Content: string(emailContent),
		Data: map[string]interface{}{
			"notificationType": "CREATE_RENTAL",
			"rentalId":         r.ID,
		},
		Targets: []misc_dto.CreateNotificationTarget{target},
	}); err != nil {
		return err
	}

	if err = s.mService.SendNotification(&misc_dto.CreateNotification{
		Title:   string(title),
		Content: string(pushContent),
		Data: map[string]interface{}{
			"notificationType": "CREATE_RENTAL",
			"rentalId":         r.ID,
		},
		Targets: []misc_dto.CreateNotificationTarget{target},
	}); err != nil {
		return err
	}

	return nil
}

func (s *service) notifyCreateContract(
	c *rental_model.ContractModel,
	r *rental_model.RentalModel,
) error {
	// get target emails and push tokens
	target, err := s._getNotificationTenantTargets(r.TenantID, r.TenantEmail)
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

	if err = s.mService.SendNotification(&misc_dto.CreateNotification{
		Title:   string(title),
		Content: string(emailContent),
		Data: map[string]interface{}{
			"notificationType": "CREATE_CONTRACT",
			"rentalId":         r.ID,
			"contractId":       c.ID,
		},
		Targets: []misc_dto.CreateNotificationTarget{target},
	}); err != nil {
		return err
	}

	if err = s.mService.SendNotification(&misc_dto.CreateNotification{
		Title:   string(title),
		Content: string(pushContent),
		Data: map[string]interface{}{
			"notificationType": "CREATE_CONTRACT",
			"rentalId":         r.ID,
			"contractId":       c.ID,
		},
		Targets: []misc_dto.CreateNotificationTarget{target},
	}); err != nil {
		return err
	}

	return nil
}

func (s *service) notifyUpdateContract(
	c *rental_model.ContractModel,
	r *rental_model.RentalModel,
	side string,
) error {
	var (
		targets []misc_dto.CreateNotificationTarget
		err     error
	)
	if side == "A" {
		target, err := s._getNotificationTenantTargets(r.TenantID, r.TenantEmail)
		if err != nil {
			return err
		}
		targets = []misc_dto.CreateNotificationTarget{target}
	} else if side == "B" {
		targets, err = s._getNotificationManagersTargets(r.PropertyID)
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

	if err = s.mService.SendNotification(&misc_dto.CreateNotification{
		Title:   string(title),
		Content: string(emailContent),
		Data: map[string]interface{}{
			"notificationType": "UPDATE_CONTRACT",
			"rentalId":         r.ID,
			"contractId":       c.ID,
		},
		Targets: targets,
	}); err != nil {
		return err
	}

	if err = s.mService.SendNotification(&misc_dto.CreateNotification{
		Title:   string(title),
		Content: string(pushContent),
		Data: map[string]interface{}{
			"notificationType": "UPDATE_CONTRACT",
			"rentalId":         r.ID,
			"contractId":       c.ID,
		},
		Targets: targets,
	}); err != nil {
		return err
	}

	return nil
}

func (s *service) notifyCreateRentalComplaint(
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
			target, err := s._getNotificationTenantTargets(r.TenantID, r.TenantEmail)
			if err != nil {
				return err
			}
			targets = []misc_dto.CreateNotificationTarget{target}
		} else if side == "B" {
			targets, err = s._getNotificationManagersTargets(r.PropertyID)
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

	if err = s.mService.SendNotification(&misc_dto.CreateNotification{
		Title:   string(title),
		Content: string(emailContent),
		Data: map[string]interface{}{
			"notificationType": "CREATE_RENTALCOMPLAINT",
			"rentalId":         r.ID,
			"complaintId":      c.ID,
		},
		Targets: targets,
	}); err != nil {
		return err
	}

	if err = s.mService.SendNotification(&misc_dto.CreateNotification{
		Title:   string(title),
		Content: string(pushContent),
		Data: map[string]interface{}{
			"notificationType": "CREATE_RENTALCOMPLAINT",
			"rentalId":         r.ID,
			"complaintId":      c.ID,
		},
		Targets: targets,
	}); err != nil {
		return err
	}

	return nil
}

func (s *service) notifyCreateComplaintReply(
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
			target, err := s._getNotificationTenantTargets(r.TenantID, r.TenantEmail)
			if err != nil {
				return err
			}
			targets = []misc_dto.CreateNotificationTarget{target}
		} else if side == "B" {
			targets, err = s._getNotificationManagersTargets(r.PropertyID)
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

	if err = s.mService.SendNotification(&misc_dto.CreateNotification{
		Title:   string(title),
		Content: string(emailContent),
		Data: map[string]interface{}{
			"notificationType": "CREATE_COMPLAINTREPLY",
			"rentalId":         r.ID,
			"complaintId":      c.ID,
		},
		Targets: targets,
	}); err != nil {
		return err
	}

	if err = s.mService.SendNotification(&misc_dto.CreateNotification{
		Title:   string(title),
		Content: string(pushContent),
		Data: map[string]interface{}{
			"notificationType": "CREATE_COMPLAINTREPLY",
			"rentalId":         r.ID,
			"complaintId":      c.ID,
		},
		Targets: targets,
	}); err != nil {
		return err
	}

	return nil
}

func (s *service) notifyUpdateComplaintStatus(
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
			target, err := s._getNotificationTenantTargets(r.TenantID, r.TenantEmail)
			if err != nil {
				return err
			}
			targets = []misc_dto.CreateNotificationTarget{target}
		} else if side == "B" {
			targets, err = s._getNotificationManagersTargets(r.PropertyID)
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

	if err = s.mService.SendNotification(&misc_dto.CreateNotification{
		Title:   string(title),
		Content: string(emailContent),
		Data: map[string]interface{}{
			"notificationType": "UPDATE_COMPLAINTSTATUS",
			"rentalId":         r.ID,
			"complaintId":      c.ID,
		},
		Targets: targets,
	}); err != nil {
		return err
	}

	if err = s.mService.SendNotification(&misc_dto.CreateNotification{
		Title:   string(title),
		Content: string(pushContent),
		Data: map[string]interface{}{
			"notificationType": "UPDATE_COMPLAINTSTATUS",
			"rentalId":         r.ID,
			"complaintId":      c.ID,
		},
		Targets: targets,
	}); err != nil {
		return err
	}

	return nil
}
