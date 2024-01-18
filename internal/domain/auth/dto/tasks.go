package dto

type TaskSendVerifyEmailPayload struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	CreatedAt int64  `json:"created_at"`
}
