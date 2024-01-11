package property

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/user2410/rrms-backend/internal/domain/property/dto"
	"github.com/user2410/rrms-backend/internal/infrastructure/database"
	sqlbuilders "github.com/user2410/rrms-backend/internal/infrastructure/database/sql_builders"
	"github.com/user2410/rrms-backend/pkg/utils/types"
)

func TestGetProperties(t *testing.T) {
	dao, err := database.NewDAO("postgresql://root:mysecret@localhost:32755/rrms?sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}
	defer dao.Close()

	repo := NewRepo(dao)
	res, err := repo.GetProperties(
		context.Background(),
		[]string{
			("8d8ec157-a6bc-4793-9a27-989386ef7d07"),
			("936f1a14-728d-4e60-9dee-6254a0547278"),
		},
		[]string{"name", "lat", "lng", "media"},
	)
	if err != nil {
		t.Fatal(err)
	}

	for _, p := range res {
		t.Log(p.Name, *p.Lat, *p.Lng, p.Features, p.Media)
	}
}

func TestSearchProperty(t *testing.T) {
	sql, args := sqlbuilders.SearchPropertyBuilder(
		[]string{"properties.id", "properties.name"},
		&dto.SearchPropertyQuery{
			PTypes:          []string{"APARTMENT", "HOUSE", "TEST"},
			PName:           types.Ptr[string]("test name"),
			PBuilding:       types.Ptr[string]("test building"),
			PProject:        types.Ptr[string]("test project"),
			PFullAddress:    types.Ptr[string]("test address"),
			PCity:           types.Ptr[string]("test city"),
			PDistrict:       types.Ptr[string]("test district"),
			PWard:           types.Ptr[string]("test ward"),
			PMinArea:        types.Ptr[float32](3000),
			PMaxArea:        types.Ptr[float32](2000),
			PNumberOfFloors: types.Ptr[int32](12),
			PYearBuilt:      types.Ptr[int32](2023),
			POrientation:    types.Ptr[string]("nw"),
			PMinFacade:      types.Ptr[int32](12),
			PIsPublic:       types.Ptr[bool](true),
			PFeatures:       []int32{1, 2, 3},
			PTags:           []string{"tag 1", "tag 2"},
			PMinCreatedAt:   types.Ptr[time.Time](time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)),
			PMaxCreatedAt:   types.Ptr[time.Time](time.Date(2010, time.November, 10, 23, 0, 0, 0, time.UTC)),
		}, "", "")
	log.Println(sql)
	log.Println(args)
}
