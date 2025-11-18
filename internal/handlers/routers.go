package handlers

import (
	"net/http"

	"logging_api/internal/handlers/auth_handler"
	"logging_api/internal/handlers/bot_handler"
	"logging_api/internal/handlers/eff_run_handler"
	"logging_api/internal/handlers/log_handler"
	"logging_api/internal/handlers/owner_handler"
	"logging_api/internal/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRoutes(
	authHandler *auth_handler.AuthHandler,
	botHandler *bot_handler.BotHandler,
	ownerHandler *owner_handler.OwnerHandler,
	logHandler *log_handler.LogHandler,
	effRunHandler *eff_run_handler.EffRunHandler,
	authMiddleware *middleware.AuthMiddleware,
) *gin.Engine {
	router := gin.Default()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/swagger/doc.json")))

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	api := router.Group("/api/v1")
	{
		auth := api.Group("/auth")
		auth.Use(authMiddleware.AuthRequired())
		{
			auth.GET("/me", authHandler.GetMe)
		}

		tokens := api.Group("/tokens")
		tokens.Use(authMiddleware.AdminRequired())
		{
			tokens.POST("", authHandler.CreateToken)
			tokens.PUT("/:token_id", authHandler.UpdateToken)
			tokens.PATCH("/:token_id/deactivate", authHandler.DeactivateToken)
			tokens.DELETE("/:token_id", authHandler.DeleteToken)
		}

		owners := api.Group("/owners")
		owners.Use(authMiddleware.AdminRequired())
		{
			owners.POST("", ownerHandler.CreateOwner)
			owners.GET("", ownerHandler.GetAllOwners)
			owners.GET("/:owner_id", ownerHandler.GetOwner)
			owners.PUT("/:owner_id", ownerHandler.UpdateOwner)
			owners.DELETE("/:owner_id", ownerHandler.DeleteOwner)
		}

		bots := api.Group("/bots")
		bots.Use(authMiddleware.AdminRequired())
		{
			bots.POST("", botHandler.CreateBot)
			bots.GET("", botHandler.GetAllBots)
			bots.GET("/:bot_id", botHandler.GetBot)
			bots.PUT("/:bot_id", botHandler.UpdateBot)
			bots.DELETE("/:bot_id", botHandler.DeleteBot)
		}

		logs := api.Group("/logs")
		logs.Use(authMiddleware.AuthRequired())
		{
			logs.POST("", logHandler.CreateLog)
		}

		effRuns := api.Group("/eff-runs")
		effRuns.Use(authMiddleware.AuthRequired())
		{
			effRuns.POST("", effRunHandler.CreateEffRun)
		}
	}

	return router
}
