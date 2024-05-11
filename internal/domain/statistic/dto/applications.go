package dto

type ApplicationStatisticResponse struct {
	NewApplicationsThisMonth []int64 `json:"newApplicationsThisMonth"`
	NewApplicationsLastMonth []int64 `json:"newApplicationsLastMonth"`
}
