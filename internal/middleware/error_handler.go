package middleware

import (
	"log/slog"

	"github.com/gofiber/fiber/v3"
	"github.com/zunarvy19/finance-backend/internal/common"
)

func ErrorHandler(c fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"
	var errorsData interface{}

	// Check if it's a fiber Error
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	// Check if it's our custom AppError
	if e, ok := err.(*common.AppError); ok {
		code = e.Code
		message = e.Message
		if e.Err != nil {
			slog.Error("AppError", "error", e.Err)
		}
	} else if err == common.ErrNotFound {
		code = fiber.StatusNotFound
		message = "Resource not found"
	} else if err == common.ErrConflict {
		code = fiber.StatusConflict
		message = "Conflict / Version mismatch"
	} else if err == common.ErrUnauthorized {
		code = fiber.StatusUnauthorized
		message = "Unauthorized"
	}

	if code == fiber.StatusInternalServerError {
		slog.Error("Unhandled Error", "error", err)
	}

	return common.SendError(c, code, message, errorsData)
}
