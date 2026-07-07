package main

import (
	"log/slog"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/recover"
	fiberLogger "github.com/gofiber/fiber/v3/middleware/logger"

	"github.com/zunarvy19/finance-backend/configs"
	"github.com/zunarvy19/finance-backend/pkg/database"
	"github.com/zunarvy19/finance-backend/pkg/validator"
	"github.com/zunarvy19/finance-backend/internal/middleware"

	"github.com/zunarvy19/finance-backend/internal/auth"
	"github.com/zunarvy19/finance-backend/internal/account"
	"github.com/zunarvy19/finance-backend/internal/category"
	"github.com/zunarvy19/finance-backend/internal/transaction"
	"github.com/zunarvy19/finance-backend/internal/ledger"
	"github.com/zunarvy19/finance-backend/internal/dashboard"
)

func main() {
	cfg := configs.LoadConfig()
	db := database.NewPostgresDB(cfg)
	validator.Init()

	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler,
	})

	app.Use(recover.New())
	app.Use(fiberLogger.New())

	authMiddleware := middleware.AuthMiddleware(cfg)

	// Repositories
	authRepo := auth.NewRepository(db)
	accountRepo := account.NewRepository(db)
	categoryRepo := category.NewRepository(db)
	txRepo := transaction.NewRepository(db)
	dashboardRepo := dashboard.NewRepository(db)

	// Services
	ledgerSvc := ledger.NewService(db, accountRepo, txRepo)
	authSvc := auth.NewService(authRepo, cfg)
	accountSvc := account.NewService(accountRepo)
	categorySvc := category.NewService(categoryRepo)
	txSvc := transaction.NewService(txRepo, ledgerSvc)
	dashboardSvc := dashboard.NewService(dashboardRepo)

	// Handlers
	authHandler := auth.NewHandler(authSvc)
	accountHandler := account.NewHandler(accountSvc)
	categoryHandler := category.NewHandler(categorySvc)
	txHandler := transaction.NewHandler(txSvc)
	dashboardHandler := dashboard.NewHandler(dashboardSvc)

	// Routes
	api := app.Group("/api")

	authGrp := api.Group("/auth")
	authGrp.Post("/register", authHandler.Register)
	authGrp.Post("/login", authHandler.Login)
	authGrp.Post("/refresh", authHandler.Refresh)
	authGrp.Post("/logout", authHandler.Logout, authMiddleware)

	accountGrp := api.Group("/accounts", authMiddleware)
	accountGrp.Post("/", accountHandler.Create)
	accountGrp.Put("/:id", accountHandler.Update)
	accountGrp.Delete("/:id", accountHandler.Delete)
	accountGrp.Get("/:id", accountHandler.Get)
	accountGrp.Get("/", accountHandler.List)

	categoryGrp := api.Group("/categories", authMiddleware)
	categoryGrp.Post("/", categoryHandler.Create)
	categoryGrp.Put("/:id", categoryHandler.Update)
	categoryGrp.Delete("/:id", categoryHandler.Delete)
	categoryGrp.Get("/:id", categoryHandler.Get)
	categoryGrp.Get("/", categoryHandler.List)

	txGrp := api.Group("/transactions", authMiddleware)
	txGrp.Post("/", txHandler.Create)
	txGrp.Put("/:id", txHandler.Update)
	txGrp.Delete("/:id", txHandler.Delete)
	txGrp.Get("/:id", txHandler.Get)
	txGrp.Get("/", txHandler.List)

	dashGrp := api.Group("/dashboard", authMiddleware)
	dashGrp.Get("/", dashboardHandler.GetDashboard)

	slog.Info("Starting server on port " + cfg.AppPort)
	if err := app.Listen(":" + cfg.AppPort); err != nil {
		slog.Error("Server failed to start", "error", err)
	}
}
