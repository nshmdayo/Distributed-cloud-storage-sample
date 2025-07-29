// Package utils provides utility functions
package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
)

// GenerateRandomID generates a random hexadecimal ID
func GenerateRandomID(length int) (string, error) {
	bytes := make([]byte, length/2)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// EnsureDir ensures that a directory exists, creating it if necessary
func EnsureDir(dirPath string) error {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		return os.MkdirAll(dirPath, 0755)
	}
	return nil
}

// GetFileSize returns the size of a file
func GetFileSize(filePath string) (int64, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return 0, err
	}
	return fileInfo.Size(), nil
}

// FileExists checks if a file exists
func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

// SplitData splits data into chunks of specified size
func SplitData(data []byte, chunkSize int) [][]byte {
	if chunkSize <= 0 {
		return [][]byte{data}
	}

	var chunks [][]byte
	for i := 0; i < len(data); i += chunkSize {
		end := i + chunkSize
		if end > len(data) {
			end = len(data)
		}
		chunks = append(chunks, data[i:end])
	}
	return chunks
}

// JoinChunks combines chunks back into original data
func JoinChunks(chunks [][]byte) []byte {
	var result []byte
	for _, chunk := range chunks {
		result = append(result, chunk...)
	}
	return result
}

// GetStoragePath returns the full path for storing a file
func GetStoragePath(baseDir, fileID string) string {
	// Create subdirectories based on first 2 characters of file ID
	// to avoid too many files in a single directory
	subDir := fileID[:2]
	return filepath.Join(baseDir, subDir, fileID)
}

// FormatBytes formats bytes into human readable format
func FormatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// ValidateFileID validates if a file ID is properly formatted
func ValidateFileID(fileID string) bool {
	if len(fileID) != 64 {
		return false
	}
	_, err := hex.DecodeString(fileID)
	return err == nil
}
