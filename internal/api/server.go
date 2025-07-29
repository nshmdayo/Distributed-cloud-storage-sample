// Package api provides HTTP API endpoints for the distributed storage system
package api

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nshmdayo/distributed-cloud-storage/internal/storage"
	"github.com/nshmdayo/distributed-cloud-storage/pkg/types"
	"github.com/sirupsen/logrus"
)

// Server represents the API server
type Server struct {
	router       *gin.Engine
	storage      storage.Storage
	chunkManager *storage.ChunkManager
	logger       *logrus.Logger
	files        map[string]*types.FileInfo // In-memory metadata store (should be replaced with proper DB)
}

// NewServer creates a new API server
func NewServer(storage storage.Storage, chunkManager *storage.ChunkManager, logger *logrus.Logger) *Server {
	server := &Server{
		storage:      storage,
		chunkManager: chunkManager,
		logger:       logger,
		files:        make(map[string]*types.FileInfo),
	}

	server.setupRoutes()
	return server
}

// setupRoutes configures the API routes
func (s *Server) setupRoutes() {
	s.router = gin.New()

	// Middleware
	s.router.Use(gin.Logger())
	s.router.Use(gin.Recovery())
	s.router.Use(s.corsMiddleware())

	// API routes
	api := s.router.Group("/api/v1")
	{
		// File operations
		api.POST("/files", s.uploadFile)
		api.GET("/files/:id", s.downloadFile)
		api.DELETE("/files/:id", s.deleteFile)
		api.GET("/files", s.listFiles)
		api.GET("/files/:id/info", s.getFileInfo)

		// Node operations
		api.GET("/node/info", s.getNodeInfo)
		api.GET("/node/stats", s.getNodeStats)

		// Health check
		api.GET("/health", s.healthCheck)
	}
}

// corsMiddleware adds CORS headers
func (s *Server) corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		c.Next()
	}
}

// uploadFile handles file upload
func (s *Server) uploadFile(c *gin.Context) {
	// Parse multipart form
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		s.logger.WithError(err).Error("Failed to parse form file")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file"})
		return
	}
	defer file.Close()

	// Read file data
	data, err := io.ReadAll(file)
	if err != nil {
		s.logger.WithError(err).Error("Failed to read file data")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
		return
	}

	// Create file info
	fileInfo := &types.FileInfo{
		ID:          types.GenerateFileID(header.Filename, data),
		Name:        header.Filename,
		ContentType: header.Header.Get("Content-Type"),
		Owner:       c.GetHeader("X-Owner"), // Simple owner identification
	}

	// Store file
	if err := s.chunkManager.StoreFile(fileInfo, data); err != nil {
		s.logger.WithError(err).Error("Failed to store file")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store file"})
		return
	}

	// Store metadata (in production, this should be in a proper database)
	s.files[fileInfo.ID] = fileInfo

	s.logger.WithFields(logrus.Fields{
		"file_id":   fileInfo.ID,
		"file_name": fileInfo.Name,
		"size":      fileInfo.Size,
	}).Info("File uploaded successfully")

	c.JSON(http.StatusOK, gin.H{
		"file_id":   fileInfo.ID,
		"file_name": fileInfo.Name,
		"size":      fileInfo.Size,
		"hash":      fileInfo.Hash,
	})
}

// downloadFile handles file download
func (s *Server) downloadFile(c *gin.Context) {
	fileID := c.Param("id")

	// Get file info
	fileInfo, exists := s.files[fileID]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	// Retrieve file data
	data, err := s.chunkManager.RetrieveFile(fileInfo)
	if err != nil {
		s.logger.WithError(err).Error("Failed to retrieve file")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve file"})
		return
	}

	// Set response headers
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileInfo.Name))
	c.Header("Content-Type", fileInfo.ContentType)
	c.Header("Content-Length", strconv.FormatInt(fileInfo.Size, 10))

	// Send file data
	c.DataFromReader(http.StatusOK, fileInfo.Size, fileInfo.ContentType, bytes.NewReader(data), nil)

	s.logger.WithFields(logrus.Fields{
		"file_id":   fileInfo.ID,
		"file_name": fileInfo.Name,
	}).Info("File downloaded successfully")
}

// deleteFile handles file deletion
func (s *Server) deleteFile(c *gin.Context) {
	fileID := c.Param("id")

	// Get file info
	fileInfo, exists := s.files[fileID]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	// Delete file chunks
	if err := s.chunkManager.DeleteFile(fileInfo); err != nil {
		s.logger.WithError(err).Error("Failed to delete file")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file"})
		return
	}

	// Remove metadata
	delete(s.files, fileID)

	s.logger.WithFields(logrus.Fields{
		"file_id":   fileInfo.ID,
		"file_name": fileInfo.Name,
	}).Info("File deleted successfully")

	c.JSON(http.StatusOK, gin.H{"message": "File deleted successfully"})
}

// listFiles handles file listing
func (s *Server) listFiles(c *gin.Context) {
	var files []gin.H

	for _, fileInfo := range s.files {
		files = append(files, gin.H{
			"id":           fileInfo.ID,
			"name":         fileInfo.Name,
			"size":         fileInfo.Size,
			"content_type": fileInfo.ContentType,
			"created_at":   fileInfo.CreatedAt,
			"owner":        fileInfo.Owner,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"files": files,
		"count": len(files),
	})
}

// getFileInfo handles file information retrieval
func (s *Server) getFileInfo(c *gin.Context) {
	fileID := c.Param("id")

	fileInfo, exists := s.files[fileID]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	c.JSON(http.StatusOK, fileInfo)
}

// getNodeInfo handles node information retrieval
func (s *Server) getNodeInfo(c *gin.Context) {
	usage, err := s.storage.GetUsage()
	if err != nil {
		s.logger.WithError(err).Error("Failed to get storage usage")
		usage = 0
	}

	c.JSON(http.StatusOK, gin.H{
		"node_id":      "node-001", // Should be dynamic
		"status":       "online",
		"storage_used": usage,
		"files_count":  len(s.files),
		"last_seen":    time.Now(),
	})
}

// getNodeStats handles node statistics retrieval
func (s *Server) getNodeStats(c *gin.Context) {
	usage, err := s.storage.GetUsage()
	if err != nil {
		s.logger.WithError(err).Error("Failed to get storage usage")
		usage = 0
	}

	filesList, err := s.storage.List()
	if err != nil {
		s.logger.WithError(err).Error("Failed to list files")
		filesList = []string{}
	}

	c.JSON(http.StatusOK, gin.H{
		"storage_usage":  usage,
		"file_count":     len(filesList),
		"metadata_count": len(s.files),
		"uptime":         time.Since(time.Now()), // Should track actual uptime
	})
}

// healthCheck handles health check requests
func (s *Server) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"timestamp": time.Now(),
		"version":   "1.0.0",
	})
}

// Start starts the HTTP server
func (s *Server) Start(addr string) error {
	s.logger.WithField("address", addr).Info("Starting API server")
	return s.router.Run(addr)
}

// GetRouter returns the gin router for testing
func (s *Server) GetRouter() *gin.Engine {
	return s.router
}
