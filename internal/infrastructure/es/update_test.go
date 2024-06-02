package es

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/elastic/go-elasticsearch/v8/typedapi/core/update"
	"github.com/stretchr/testify/require"
	"github.com/user2410/rrms-backend/internal/domain/listing/dto"
	"github.com/user2410/rrms-backend/internal/utils/types"
)

func TestUpdateListing(t *testing.T) {
	data := &dto.UpdateListing{
		LeaseTerm: types.Ptr[int32](36),
		Tags:      []string{"tag1", "tag2"},
	}

	// update es document
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
	require.NoError(t, err)

	_, err = typedClient.Update(string(LISTINGINDEX), "00808905-a3c4-4c5c-b989-15e6341bceca").
		Request(&update.Request{
			Doc: json.RawMessage(docByte),
		}).
		Do(context.Background())

	require.NoError(t, err)
}

func TestUpdateListingStatus(t *testing.T) {
	// update es document
	doc := map[string]interface{}{
		"active": false,
	}
	docByte, err := json.Marshal(doc)
	require.NoError(t, err)
	_, err = typedClient.Update(
		string(LISTINGINDEX),
		"00808905-a3c4-4c5c-b989-15e6341bceca").
		Request(&update.Request{
			Doc: json.RawMessage(docByte),
		}).
		Do(context.Background())

	require.NoError(t, err)
}

func TestUpdateListingExpiration(t *testing.T) {
	// 2024-08-15 09:05:50.585+00
	// expiredAt, err := time.Parse(time.RFC3339, "2024-08-15T09:05:50.585Z")
	// require.NoError(t, err)

	// update es document
	doc := map[string]interface{}{
		"expired_at": "2024-08-15T16:05:50.585+07:00",
	}
	docByte, err := json.Marshal(doc)
	require.NoError(t, err)
	_, err = typedClient.Update(string(LISTINGINDEX), "00808905-a3c4-4c5c-b989-15e6341bceca").
		Request(&update.Request{
			Doc: json.RawMessage(docByte),
		}).
		Do(context.Background())

	require.NoError(t, err)
}

func TestUpdateListingPriority(t *testing.T) {
	// update es document
	doc := map[string]interface{}{
		"priority": 3,
	}
	docByte, err := json.Marshal(doc)
	require.NoError(t, err)
	_, err = typedClient.Update(string(LISTINGINDEX), "00808905-a3c4-4c5c-b989-15e6341bceca").
		Request(&update.Request{
			Doc: json.RawMessage(docByte),
		}).
		Do(context.Background())

	require.NoError(t, err)
}
