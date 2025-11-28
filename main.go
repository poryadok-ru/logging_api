package main

import (
	"fmt"
	"log"

	"logging_api/configs"
	_ "logging_api/docs"
	"logging_api/internal/handlers"
	"logging_api/internal/handlers/auth_handler"
	"logging_api/internal/handlers/bot_handler"
	"logging_api/internal/handlers/eff_run_handler"
	"logging_api/internal/handlers/log_handler"
	"logging_api/internal/handlers/owner_handler"
	"logging_api/internal/middleware"
	authservice "logging_api/internal/service/auth_service"
	botservice "logging_api/internal/service/bot_service"
	effrunservice "logging_api/internal/service/eff_run_service"
	logservice "logging_api/internal/service/log_service"
	ownerservice "logging_api/internal/service/owner_service"
	authrepo "logging_api/internal/storage/auth_repo"
	botrepo "logging_api/internal/storage/bot_repo"
	effrunrepo "logging_api/internal/storage/eff_run_repo"
	logrepo "logging_api/internal/storage/log_repo"
	ownerrepo "logging_api/internal/storage/owner_repo"
	"logging_api/pkg/postgres"
)

// @title Logging API
// @version 1.0
// @description API для управления логами ботов. Для авторизации используйте заголовок Authorization: Bearer {token}.
// @host api.automation.poryadok.ru
// @BasePath /logging
// @schemes https
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	config, err := configs.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := postgres.Connect(
		config.Database.Host,
		config.Database.Port,
		config.Database.User,
		config.Database.Password,
		config.Database.DBName,
		config.Database.SSLMode,
	)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Println("Successfully connected to database!")

	authRepo := authrepo.NewAuthRepo(db)
	botRepo := botrepo.NewBotRepo(db)
	ownerRepo := ownerrepo.NewOwnerRepo(db)
	logRepo := logrepo.NewLogRepo(db)
	effRunRepo := effrunrepo.NewEffRunRepo(db)

	authService := authservice.NewAuthService(authRepo, botRepo)
	botService := botservice.NewBotService(botRepo)
	ownerService := ownerservice.NewOwnerService(ownerRepo)
	logService := logservice.NewLogService(logRepo)
	effRunService := effrunservice.NewEffRunService(effRunRepo)

	authMiddleware := middleware.NewAuthMiddleware(authService)

	authHandler := auth_handler.NewAuthHandler(authService)
	botHandler := bot_handler.NewBotHandler(botService)
	ownerHandler := owner_handler.NewOwnerHandler(ownerService)
	logHandler := log_handler.NewLogHandler(logService)
	effRunHandler := eff_run_handler.NewEffRunHandler(effRunService)

	router := handlers.SetupRoutes(authHandler, botHandler, ownerHandler, logHandler, effRunHandler, authMiddleware)

	addr := fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port)
	log.Printf("Starting server on %s", addr)
	log.Printf("Swagger доступен по адресу: http://%s/swagger/index.html", addr)

	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
