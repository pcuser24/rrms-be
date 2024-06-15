package service

import (
	"context"
	"encoding/json"

	"github.com/elastic/go-elasticsearch/v8/typedapi/core/update"
	"github.com/google/uuid"
	"github.com/user2410/rrms-backend/internal/domain/listing/dto"
	"github.com/user2410/rrms-backend/internal/infrastructure/es"
)

func (s *service) UpdateListing(id uuid.UUID, data *dto.UpdateListing) error {
	err := s.domainRepo.ListingRepo.UpdateListing(context.Background(), id, data)
	if err != nil {
		return err
	}

	// update es document
	client := s.esClient.GetTypedClient()
	doc := make(map[string]interface{})
	if data.Title != nil {
		doc["title"] = *data.Title
	}
	if data.Description != nil {
		doc["description"] = *data.Description
	}
	if data.Price != nil {
		doc["price"] = *data.Price
	}
	if data.SecurityDeposit != nil {
		doc["security_deposit"] = *data.SecurityDeposit
	}
	if data.LeaseTerm != nil {
		doc["lease_term"] = *data.LeaseTerm
	}
	if data.PetsAllowed != nil {
		doc["pets_allowed"] = *data.PetsAllowed
	}
	if data.NumberOfResidents != nil {
		doc["number_of_residents"] = *data.NumberOfResidents
	}
	if data.Tags != nil {
		tags := make([]map[string]string, 0, len(data.Tags))
		for _, tag := range data.Tags {
			tags = append(tags, map[string]string{"name": tag})
		}
		doc["tags"] = tags
	}
	docByte, err := json.Marshal(doc)
	if err != nil {
		return err
	}
	_, err = client.Update(string(es.LISTINGINDEX), id.String()).
		Request(&update.Request{
			Doc: json.RawMessage(docByte),
		}).
		Do(context.Background())

	return err
}

func (s *service) UpdateListingStatus(id uuid.UUID, active bool) error {
	err := s.domainRepo.ListingRepo.UpdateListingStatus(context.Background(), id, active)
	if err != nil {
		return err
	}

	// update es document
	doc := map[string]interface{}{
		"active": active,
	}
	docByte, err := json.Marshal(doc)
	if err != nil {
		return err
	}
	client := s.esClient.GetTypedClient()
	_, err = client.Update(string(es.LISTINGINDEX), id.String()).
		Request(&update.Request{
			Doc: json.RawMessage(docByte),
		}).
		Do(context.Background())

	return err
}

func (s *service) UpdateListingExpiration(id uuid.UUID, duration int64) error {
	expiredAt, err := s.domainRepo.ListingRepo.UpdateListingExpiration(context.Background(), id, duration)
	if err != nil {
		return err
	}

	// update es document
	doc := map[string]interface{}{
		"expired_at": expiredAt,
	}
	docByte, err := json.Marshal(doc)
	if err != nil {
		return err
	}
	client := s.esClient.GetTypedClient()
	_, err = client.Update(string(es.LISTINGINDEX), id.String()).
		Request(&update.Request{
			Doc: json.RawMessage(docByte),
		}).
		Do(context.Background())

	return err
}

func (s *service) UpdateListingPriority(id uuid.UUID, priority int) error {
	err := s.domainRepo.ListingRepo.UpdateListingPriority(context.Background(), id, priority)
	if err != nil {
		return err
	}

	// update es document
	doc := map[string]interface{}{
		"priority": priority,
	}
	docByte, err := json.Marshal(doc)
	if err != nil {
		return err
	}
	client := s.esClient.GetTypedClient()
	_, err = client.Update(string(es.LISTINGINDEX), id.String()).
		Request(&update.Request{
			Doc: json.RawMessage(docByte),
		}).
		Do(context.Background())

	return err
}
