package dto

type PutObjectPresignRequest struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
	Type string `json:"type"`
}
