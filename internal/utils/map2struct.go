package utils

import (
	"encoding/json"
	"fmt"
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

	if errs := ValidateStruct(d); len(errs) > 0 && errs[0].Error {
		return fmt.Errorf("%s", GetValidationError(errs))
	}

	return nil
}
