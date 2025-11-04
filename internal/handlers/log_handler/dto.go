package log_handler

type CreateLogRequest struct {
	Status string `json:"status" binding:"required,oneof=Debug Info Warning Error Critical" example:"Info" enums:"Debug,Info,Warning,Error,Critical"`
	Msg    string `json:"msg" binding:"required,min=1" example:"Операция выполнена успешно"`
}

