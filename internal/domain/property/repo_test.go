package property

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/huandu/go-sqlbuilder"
)

func TestXxx(t *testing.T) {
	aid := []int64{1, 2, 3}
	// transform aid to []interface{}
	aid_i := make([]interface{}, len(aid))
	for i, v := range aid {
		aid_i[i] = v
	}
	uid, _ := uuid.Parse("d01bfb0b-dfbf-442f-8674-b0823b5eac60")

	ib := sqlbuilder.PostgreSQL.NewDeleteBuilder()
	ib.DeleteFrom("property_amenity")
	ib.Where(
		ib.Equal("property_id", uid),
		ib.In("amenity_id", aid_i...),
	)
	sql, args := ib.Build()
	fmt.Println(sql)
	fmt.Println(args)

}
