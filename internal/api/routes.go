package api

import (
	"backend/internal/api/handlers"
	"backend/internal/middleware"

	_ "backend/docs" // Import generated docs

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRoutes(router *gin.Engine,
	docHandler *handlers.DocumentHandler,
	questionHandler *handlers.QuestionHandler,
	authHandler *handlers.AuthHandler,
	authMiddleware *middleware.AuthMiddleware,
) {
	// Public Routes
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Swagger UI
	router.GET("/swagger", func(c *gin.Context) {
		c.Redirect(301, "/swagger/index.html")
	})
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/api/v1")

	// Auth Routes
	auth := api.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}

	// Protected Routes (if auth enabled)
	// For now, assume authMiddleware handles token validation if keys are present
	if authMiddleware != nil {
		api.Use(authMiddleware.ValidateToken())
	}

	api.POST("/questions/generate", questionHandler.Generate)
	api.POST("/documents/upload", docHandler.Upload)
}
