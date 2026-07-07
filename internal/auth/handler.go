package auth

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/zunarvy19/finance-backend/internal/common"
	"github.com/zunarvy19/finance-backend/pkg/validator"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Register(c fiber.Ctx) error {
	var req RegisterRequest
	if err := c.Bind().Body(&req); err != nil {
		return common.SendError(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}

	if errors := validator.ValidateStruct(req); errors != nil {
		return common.SendError(c, fiber.StatusBadRequest, "Validation failed", errors)
	}

	if err := h.service.Register(c.Context(), &req); err != nil {
		return err // handled by error middleware
	}

	return common.SendSuccess(c, fiber.StatusCreated, "User registered successfully", nil)
}

func (h *Handler) Login(c fiber.Ctx) error {
	var req LoginRequest
	if err := c.Bind().Body(&req); err != nil {
		return common.SendError(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}

	if errors := validator.ValidateStruct(req); errors != nil {
		return common.SendError(c, fiber.StatusBadRequest, "Validation failed", errors)
	}

	res, err := h.service.Login(c.Context(), &req)
	if err != nil {
		return err
	}

	return common.SendSuccess(c, fiber.StatusOK, "Login successful", res)
}

func (h *Handler) Refresh(c fiber.Ctx) error {
	var req RefreshRequest
	if err := c.Bind().Body(&req); err != nil {
		return common.SendError(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}

	if errors := validator.ValidateStruct(req); errors != nil {
		return common.SendError(c, fiber.StatusBadRequest, "Validation failed", errors)
	}

	res, err := h.service.Refresh(c.Context(), &req)
	if err != nil {
		return err
	}

	return common.SendSuccess(c, fiber.StatusOK, "Token refreshed successfully", res)
}

func (h *Handler) Logout(c fiber.Ctx) error {
	userIDStr, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return common.ErrUnauthorized
	}

	if err := h.service.Logout(c.Context(), userIDStr); err != nil {
		return err
	}

	return common.SendSuccess(c, fiber.StatusOK, "Logout successful", nil)
}
