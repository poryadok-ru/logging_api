package owner_handler

import (
	"net/http"

	"logging_api/internal/models"
	customerrors "logging_api/internal/utils/errors"
	"logging_api/internal/utils/validator_error_handling"

	"github.com/gin-gonic/gin"
)

type OwnerService interface {
	CreateOwner(fullName string, isActive bool) (*models.Owner, error)
	GetOwnerByID(ownerID string) (*models.Owner, error)
	GetAllOwners() ([]*models.Owner, error)
	UpdateOwner(ownerID string, fullName *string, isActive *bool) (*models.Owner, error)
	DeleteOwner(ownerID string) error
}

type OwnerHandler struct {
	ownerService OwnerService
}

func NewOwnerHandler(ownerService OwnerService) *OwnerHandler {
	return &OwnerHandler{
		ownerService: ownerService,
	}
}

// @Summary Создать владельца
// @Description Создаёт нового владельца (требуется админский токен)
// @Tags owners
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateOwnerRequest true "Данные владельца"
// @Success 201 {object} models.Owner
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /v1/owners [post]
func (h *OwnerHandler) CreateOwner(c *gin.Context) {
	var request CreateOwnerRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		if validationErrs := validator_error_handling.ValidateError(err); validationErrs != nil {
			c.JSON(http.StatusBadRequest, gin.H{"errors": validationErrs.Errors()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат данных"})
		return
	}

	owner, err := h.ownerService.CreateOwner(request.FullName, request.IsActive)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, owner)
}

// @Summary Получить владельца
// @Description Возвращает информацию о владельце по ID (требуется админский токен)
// @Tags owners
// @Produce json
// @Security BearerAuth
// @Param owner_id path string true "ID владельца (UUID)"
// @Success 200 {object} models.Owner
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /v1/owners/{owner_id} [get]
func (h *OwnerHandler) GetOwner(c *gin.Context) {
	ownerID := c.Param("owner_id")

	owner, err := h.ownerService.GetOwnerByID(ownerID)
	if err != nil {
		if customerrors.IsNotFound(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, owner)
}

// @Summary Получить всех владельцев
// @Description Возвращает список всех владельцев (требуется админский токен)
// @Tags owners
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.Owner
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /v1/owners [get]
func (h *OwnerHandler) GetAllOwners(c *gin.Context) {
	owners, err := h.ownerService.GetAllOwners()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, owners)
}

// @Summary Обновить владельца
// @Description Обновляет данные владельца (требуется админский токен)
// @Tags owners
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param owner_id path string true "ID владельца (UUID)"
// @Param request body UpdateOwnerRequest true "Обновлённые данные"
// @Success 200 {object} models.Owner
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /v1/owners/{owner_id} [put]
func (h *OwnerHandler) UpdateOwner(c *gin.Context) {
	ownerID := c.Param("owner_id")

	var request UpdateOwnerRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		if validationErrs := validator_error_handling.ValidateError(err); validationErrs != nil {
			c.JSON(http.StatusBadRequest, gin.H{"errors": validationErrs.Errors()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат данных"})
		return
	}

	owner, err := h.ownerService.UpdateOwner(ownerID, request.FullName, request.IsActive)
	if err != nil {
		if customerrors.IsNotFound(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, owner)
}

// @Summary Удалить владельца
// @Description Полностью удаляет владельца из базы данных (требуется админский токен)
// @Tags owners
// @Produce json
// @Security BearerAuth
// @Param owner_id path string true "ID владельца (UUID)"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /v1/owners/{owner_id} [delete]
func (h *OwnerHandler) DeleteOwner(c *gin.Context) {
	ownerID := c.Param("owner_id")

	err := h.ownerService.DeleteOwner(ownerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "владелец удалён"})
}
