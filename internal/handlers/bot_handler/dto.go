package bot_handler

type CreateBotRequest struct {
	Code        string   `json:"code" binding:"required,min=2,max=50" example:"BOT_001"`
	Name        string   `json:"name" binding:"required,min=2,max=255" example:"Telegram Bot"`
	BotType     string   `json:"bot_type" binding:"required,oneof=AI Backend Frontend Robot" example:"Backend"`
	Language    string   `json:"language" binding:"required,oneof=Python Go N8N PIX JS C Other" example:"Python"`
	Description *string  `json:"description,omitempty" example:"Бот для обработки сообщений"`
	Tags        []string `json:"tags,omitempty" example:"telegram,bot"`
	OwnerID     *string  `json:"owner_id,omitempty" binding:"omitempty,uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	IsActive    bool     `json:"is_active" example:"true"`
}

type UpdateBotRequest struct {
	Code        *string  `json:"code,omitempty" binding:"omitempty,min=2,max=50" example:"BOT_002"`
	Name        *string  `json:"name,omitempty" binding:"omitempty,min=2,max=255" example:"Discord Bot"`
	BotType     *string  `json:"bot_type,omitempty" binding:"omitempty,oneof=AI Backend Frontend Robot" example:"AI"`
	Language    *string  `json:"language,omitempty" binding:"omitempty,oneof=Python Go N8N PIX JS C Other" example:"Go"`
	Description *string  `json:"description,omitempty" example:"Обновлённое описание"`
	Tags        []string `json:"tags,omitempty" example:"discord,ai"`
	OwnerID     *string  `json:"owner_id,omitempty" binding:"omitempty,uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	IsActive    *bool    `json:"is_active,omitempty" example:"false"`
}

