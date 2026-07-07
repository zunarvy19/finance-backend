package account

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

func (h *Handler) Create(c fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return common.ErrUnauthorized
	}

	var req CreateAccountRequest
	if err := c.Bind().Body(&req); err != nil {
		return common.SendError(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}

	if errors := validator.ValidateStruct(req); errors != nil {
		return common.SendError(c, fiber.StatusBadRequest, "Validation failed", errors)
	}

	res, err := h.service.Create(c.Context(), userID, &req)
	if err != nil {
		return err
	}

	return common.SendSuccess(c, fiber.StatusCreated, "Account created successfully", res)
}

func (h *Handler) Update(c fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return common.ErrUnauthorized
	}

	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return common.SendError(c, fiber.StatusBadRequest, "Invalid account ID", nil)
	}

	var req UpdateAccountRequest
	if err := c.Bind().Body(&req); err != nil {
		return common.SendError(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}

	if errors := validator.ValidateStruct(req); errors != nil {
		return common.SendError(c, fiber.StatusBadRequest, "Validation failed", errors)
	}

	res, err := h.service.Update(c.Context(), id, userID, &req)
	if err != nil {
		return err
	}

	return common.SendSuccess(c, fiber.StatusOK, "Account updated successfully", res)
}

func (h *Handler) Delete(c fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return common.ErrUnauthorized
	}

	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return common.SendError(c, fiber.StatusBadRequest, "Invalid account ID", nil)
	}

	if err := h.service.Delete(c.Context(), id, userID); err != nil {
		return err
	}

	return common.SendSuccess(c, fiber.StatusOK, "Account deleted successfully", nil)
}

func (h *Handler) Get(c fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return common.ErrUnauthorized
	}

	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return common.SendError(c, fiber.StatusBadRequest, "Invalid account ID", nil)
	}

	res, err := h.service.GetByID(c.Context(), id, userID)
	if err != nil {
		return err
	}

	return common.SendSuccess(c, fiber.StatusOK, "Account retrieved successfully", res)
}

func (h *Handler) List(c fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return common.ErrUnauthorized
	}

	res, err := h.service.List(c.Context(), userID)
	if err != nil {
		return err
	}

	return common.SendSuccess(c, fiber.StatusOK, "Accounts retrieved successfully", res)
}
