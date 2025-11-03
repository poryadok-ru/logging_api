package auth_handler

type CreateTokenRequest struct {
	BotID     *string `json:"bot_id,omitempty" binding:"omitempty,uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	TokenName string  `json:"token_name" binding:"required,min=3,max=100" example:"Production Server"`
	IsAdmin   bool    `json:"is_admin" example:"false"`
}

type UpdateTokenRequest struct {
	TokenName string `json:"token_name" binding:"required,min=3,max=100" example:"Production Server v2"`
}

type TokenResponse struct {
	Token string `json:"token" example:"550e8400-e29b-41d4-a716-446655440000"`
}
