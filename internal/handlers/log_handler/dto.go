package log_handler

type CreateLogRequest struct {
	Status string `json:"status" binding:"required,oneof=success warning error" example:"success"`
	Msg    string `json:"msg" binding:"required,min=1" example:"Операция выполнена успешно"`
}

