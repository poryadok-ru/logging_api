package owner_handler

type CreateOwnerRequest struct {
	FullName string `json:"full_name" binding:"required,min=2,max=255" example:"Иван Иванов"`
	IsActive bool   `json:"is_active" example:"true"`
}

type UpdateOwnerRequest struct {
	FullName *string `json:"full_name,omitempty" binding:"omitempty,min=2,max=255" example:"Иван Петров"`
	IsActive *bool   `json:"is_active,omitempty" example:"true"`
}
