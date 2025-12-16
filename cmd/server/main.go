package main

import (
	"context"
	"log"
	"time"

	"backend/internal/api"
	"backend/internal/api/handlers"
	"backend/internal/auth"
	"backend/internal/config"
	"backend/internal/embedding"
	"backend/internal/generation"
	"backend/internal/ingestion"
	"backend/internal/middleware"
	"backend/internal/rag"
	"backend/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	_ "github.com/pgvector/pgvector-go"
)

// @title           Nexus RAG API
// @version         1.0
// @description     API for RAG-based Physics Question Generation and Document Ingestion.
// @termsOfService  http://swagger.io/terms/

// @contact.name    API Support
// @contact.url     http://www.swagger.io/support
// @contact.email   support@swagger.io

// @license.name    Apache 2.0
// @license.url     http://www.apache.org/licenses/LICENSE-2.0.html

// @host            localhost:8080
// @BasePath        /api/v1
// @schemes         http

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	// 1. Load Config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 2. Initialize DB
	db, err := sqlx.Connect("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Printf("Warning: Failed to connect to DB: %v (Proceeding might fail if DB needed)", err)
	} else {
		db.SetMaxOpenConns(25)
		db.SetMaxIdleConns(25)
		db.SetConnMaxLifetime(5 * time.Minute)
	}

	// 3. Initialize AI Clients
	ctx := context.Background()
	embedder, err := embedding.NewGeminiClient(ctx, cfg.GeminiAPIKey)
	if err != nil {
		log.Fatalf("Failed to init embedding client: %v", err)
	}

	generator, err := generation.NewGeminiClient(ctx, cfg.GeminiAPIKey)
	if err != nil {
		log.Fatalf("Failed to init generation client: %v", err)
	}

	// 4. Initialize Core Components
	vectorRepo := repository.NewPostgresVectorRepo(db)

	// Ingestion
	pdfParser := ingestion.NewPDFParser()
	chunker := ingestion.NewChunker(1000, 200) // 1000 chars, 200 overlap
	ingestionService := ingestion.NewIngestionService(pdfParser, chunker, embedder, vectorRepo)

	// RAG
	retriever := rag.NewRetriever(embedder, vectorRepo)
	generatorService := rag.NewGeneratorService(generator, retriever)

	// Auth
	authService := auth.NewAuthService(cfg.SupabaseURL, cfg.SupabaseAnonKey)

	// 5. Initialize Handlers
	docHandler := handlers.NewDocumentHandler(ingestionService)
	questionHandler := handlers.NewQuestionHandler(generatorService)
	authHandler := handlers.NewAuthHandler(authService)

	// 6. Middleware
	authMiddleware := middleware.NewAuthMiddleware(cfg.SupabaseJWTSecret)

	// 7. Setup Router
	gin.SetMode(gin.ReleaseMode) // Set to release for production
	if cfg.Environment == "development" {
		gin.SetMode(gin.DebugMode)
	}

	router := gin.Default()
	api.SetupRoutes(router, docHandler, questionHandler, authHandler, authMiddleware)

	// 8. Run
	port := cfg.Port
	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Server failed to allow: %v", err)
	}
}
