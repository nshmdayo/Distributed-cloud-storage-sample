// Package main provides the storage node entrypoint
package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

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
		Use:   "node",
		Short: "Distributed Cloud Storage Node",
		Long:  "Storage node for the distributed cloud storage system",
		Run:   runStorageNode,
	}

	rootCmd.Flags().StringVarP(&configFile, "config", "c", "", "Config file path")
	rootCmd.Flags().StringVarP(&logLevel, "log-level", "l", "info", "Log level (debug, info, warn, error)")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func runStorageNode(cmd *cobra.Command, args []string) {
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
		"max_storage": cfg.Node.MaxStorage,
		"chunk_size":  cfg.Node.ChunkSize,
		"replicas":    cfg.Node.Replicas,
	}).Info("Starting storage node with configuration")

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
	_ = storage.NewChunkManager(fileStorage, encKey, cfg.Node.ChunkSize, logger)

	logger.Info("Storage node initialized successfully")

	// Setup graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start node services (P2P network, sync, etc.)
	// This is where we would initialize P2P networking, blockchain connectivity, etc.

	logger.Info("Storage node started, waiting for shutdown signal...")

	// Wait for shutdown signal
	<-sigChan
	logger.Info("Received shutdown signal, stopping storage node...")

	// Cleanup and graceful shutdown
	logger.Info("Storage node stopped")
}
