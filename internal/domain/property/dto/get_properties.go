package dto

type GetPropertiesQuery struct {
	Fields []string `query:"fields" validate:"omitempty,dive,oneof=name building project area number_of_floors year_built orientation entrance_width facade full_address city district ward lat lng place_url description type is_public created_at updated_at features tags media"`
}
