// Package main provides the API server entrypoint
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/nshmdayo/distributed-cloud-storage/internal/api"
	"github.com/nshmdayo/distributed-cloud-storage/internal/config"
	"github.com/nshmdayo/distributed-cloud-storage/internal/crypto"
	"github.com/nshmdayo/distributed-cloud-storage/internal/storage"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	configFile string
	logLevel   string
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "api",
		Short: "Distributed Cloud Storage API Server",
		Long:  "API server for the distributed cloud storage system",
		Run:   runAPIServer,
	}

	rootCmd.Flags().StringVarP(&configFile, "config", "c", "", "Config file path")
	rootCmd.Flags().StringVarP(&logLevel, "log-level", "l", "info", "Log level (debug, info, warn, error)")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func runAPIServer(cmd *cobra.Command, args []string) {
	// Setup logger
	logger := logrus.New()
	if level, err := logrus.ParseLevel(logLevel); err == nil {
		logger.SetLevel(level)
	}
	logger.SetFormatter(&logrus.JSONFormatter{})

	// Load configuration
	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if err := cfg.Validate(); err != nil {
		log.Fatalf("Invalid config: %v", err)
	}

	logger.WithFields(logrus.Fields{
		"storage_dir": cfg.Storage.Path,
		"api_port":    cfg.API.Port,
		"max_storage": cfg.Node.MaxStorage,
		"chunk_size":  cfg.Node.ChunkSize,
	}).Info("Starting API server with configuration")

	// Initialize storage
	fileStorage, err := storage.NewFileStorage(cfg.Storage.Path, logger)
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}

	// Generate or load encryption key
	encKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatalf("Failed to generate encryption key: %v", err)
	}

	// Initialize chunk manager
	chunkManager := storage.NewChunkManager(fileStorage, encKey, cfg.Node.ChunkSize, logger)

	// Initialize API server
	server := api.NewServer(fileStorage, chunkManager, logger)

	// Start server
	addr := fmt.Sprintf("%s:%d", cfg.API.Host, cfg.API.Port)
	logger.WithField("address", addr).Info("API server starting")

	if err := server.Start(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
