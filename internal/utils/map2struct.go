package utils

import (
	"encoding/json"
	"fmt"
	"github.com/user2410/rrms-backend/internal/utils/validation"
)

// Map2JSONStruct converts an interface{} to a json-tagged struct and validate it
func Map2JSONStruct(s interface{}, d interface{}) error {
	b, err := json.Marshal(s)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, d)
	if err != nil {
		return err
	}

	if errs := validation.ValidateStruct(nil, d); len(errs) > 0 {
		return fmt.Errorf("%s", validation.GetValidationError(errs))
	}

	return nil
}
