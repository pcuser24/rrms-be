package utils

import (
	"encoding/json"
	"fmt"
	"reflect"

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

func RedisStructToMap(v interface{}) (map[string]interface{}, error) {
	out := make(map[string]interface{})
	val := reflect.ValueOf(v).Elem()
	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		tag := field.Tag.Get("redis")
		if tag == "" || tag == "-" {
			tag = field.Name
		}
		value := val.Field(i).Interface()

		switch reflect.TypeOf(value).Kind() {
		case reflect.Slice, reflect.Array, reflect.Struct:
			bytes, err := json.Marshal(value)
			if err != nil {
				return nil, err
			}
			out[tag] = string(bytes)
		default:
			out[tag] = value
		}
	}
	return out, nil
}
