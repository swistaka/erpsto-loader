// main.go
package main

import (
	"arm/erpsto-loader/internal/config"
	"arm/erpsto-loader/internal/logger"
	"arm/erpsto-loader/internal/processor"
	"arm/erpsto-loader/internal/storage"

	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	startTime := time.Now()

	logger := logger.NewFileLogger("error.log")

	cfg, err := config.ParseCommandLineArgs()
	if err != nil {
		logger.LogError(fmt.Errorf("config error: %w", err))
		log.Fatalf("config error: %v", err)
		os.Exit(1)
	}

	db, err := storage.NewSQLiteStorage("erpsto.db")
	if err != nil {
		logger.LogError(fmt.Errorf("database init error: %w", err))
		log.Fatalf("database init error: %v", err)
		os.Exit(1)
	}
	defer db.Close()

	if err := processor.Process(cfg, db, logger); err != nil {
		logger.LogError(err)
		log.Fatalf("Processing error: %v", err)
	}

	log.Printf("data loaded successfully in %v\n", time.Since(startTime))
}
