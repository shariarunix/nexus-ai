package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"backend/internal/config"
	"backend/internal/embedding"
	"backend/internal/ingestion"
	"backend/internal/repository"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	_ "github.com/pgvector/pgvector-go"
)

func main() {
	// Flags
	filePath := flag.String("file", "", "Path to PDF file")
	subject := flag.String("subject", "", "Subject of the book")
	chapter := flag.Int("chapter", 0, "Chapter number")
	language := flag.String("lang", "", "Language (en or bn)")
	flag.Parse()

	if *filePath == "" || *subject == "" || *chapter == 0 || *language == "" {
		fmt.Println("Usage: cli -file <path> -subject <subject> -chapter <num> -lang <en|bn>")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// 1. Load Config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 2. Initialize DB
	db, err := sqlx.Connect("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer db.Close()

	// 3. Initialize Embedder
	ctx := context.Background()
	embedder, err := embedding.NewGeminiClient(ctx, cfg.GeminiAPIKey)
	if err != nil {
		log.Fatalf("Failed to init embedding client: %v", err)
	}

	// 4. Initialize Components
	vectorRepo := repository.NewPostgresVectorRepo(db)
	pdfParser := ingestion.NewPDFParser()
	chunker := ingestion.NewChunker(1000, 200)
	ingestionService := ingestion.NewIngestionService(pdfParser, chunker, embedder, vectorRepo)

	// 5. Open File
	file, err := os.Open(*filePath)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		log.Fatalf("Failed to stat file: %v", err)
	}

	// 6. Ingest
	fmt.Printf("Ingesting %s (Size: %d bytes)...\n", *filePath, fileInfo.Size())
	start := time.Now()

	err = ingestionService.Ingest(ctx, file, fileInfo.Size(), *subject, *chapter, *language)
	if err != nil {
		log.Fatalf("Ingestion failed: %v", err)
	}

	fmt.Printf("Successfully ingested document in %v\n", time.Since(start))
}
