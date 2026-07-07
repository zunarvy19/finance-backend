package dashboard

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/zunarvy19/finance-backend/internal/common"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetDashboard(c fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return common.ErrUnauthorized
	}

	res, err := h.service.GetDashboard(c.Context(), userID)
	if err != nil {
		return err
	}

	return common.SendSuccess(c, fiber.StatusOK, "Dashboard retrieved successfully", res)
}
