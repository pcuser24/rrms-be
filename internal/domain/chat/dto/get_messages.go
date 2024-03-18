package dto

type GetMessagesOfGroupQuery struct {
	Offset int32 `query:"offset" validate:"omitempty,gte=0"`
	Limit  int32 `query:"limit" validate:"omitempty,gte=0"`
}
