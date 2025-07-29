// Package main provides the client CLI entrypoint
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	serverURL string
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "client",
		Short: "Distributed Cloud Storage Client",
		Long:  "Client CLI for the distributed cloud storage system",
	}

	rootCmd.PersistentFlags().StringVarP(&serverURL, "server", "s", "http://localhost:8080", "Server URL")

	// Upload command
	var uploadCmd = &cobra.Command{
		Use:   "upload [file]",
		Short: "Upload a file",
		Args:  cobra.ExactArgs(1),
		Run:   uploadFile,
	}

	// Download command
	var downloadCmd = &cobra.Command{
		Use:   "download [file-id] [output-path]",
		Short: "Download a file",
		Args:  cobra.ExactArgs(2),
		Run:   downloadFile,
	}

	// List command
	var listCmd = &cobra.Command{
		Use:   "list",
		Short: "List all files",
		Run:   listFiles,
	}

	// Delete command
	var deleteCmd = &cobra.Command{
		Use:   "delete [file-id]",
		Short: "Delete a file",
		Args:  cobra.ExactArgs(1),
		Run:   deleteFile,
	}

	// Info command
	var infoCmd = &cobra.Command{
		Use:   "info [file-id]",
		Short: "Get file information",
		Args:  cobra.ExactArgs(1),
		Run:   getFileInfo,
	}

	rootCmd.AddCommand(uploadCmd, downloadCmd, listCmd, deleteCmd, infoCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func uploadFile(cmd *cobra.Command, args []string) {
	filePath := args[0]

	// Open file
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	// Create multipart form
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	// Add file
	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		log.Fatalf("Failed to create form file: %v", err)
	}

	if _, err := io.Copy(part, file); err != nil {
		log.Fatalf("Failed to copy file: %v", err)
	}

	writer.Close()

	// Make request
	req, err := http.NewRequest("POST", serverURL+"/api/v1/files", &body)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to upload file: %v", err)
	}
	defer resp.Body.Close()

	// Parse response
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Fatalf("Failed to parse response: %v", err)
	}

	if resp.StatusCode == http.StatusOK {
		fmt.Printf("File uploaded successfully!\n")
		fmt.Printf("File ID: %s\n", result["file_id"])
		fmt.Printf("File Name: %s\n", result["file_name"])
		fmt.Printf("Size: %v bytes\n", result["size"])
		fmt.Printf("Hash: %s\n", result["hash"])
	} else {
		fmt.Printf("Upload failed: %v\n", result["error"])
	}
}

func downloadFile(cmd *cobra.Command, args []string) {
	fileID := args[0]
	outputPath := args[1]

	// Make request
	resp, err := http.Get(serverURL + "/api/v1/files/" + fileID)
	if err != nil {
		log.Fatalf("Failed to download file: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)
		log.Fatalf("Download failed: %v", result["error"])
	}

	// Create output file
	outFile, err := os.Create(outputPath)
	if err != nil {
		log.Fatalf("Failed to create output file: %v", err)
	}
	defer outFile.Close()

	// Copy data
	if _, err := io.Copy(outFile, resp.Body); err != nil {
		log.Fatalf("Failed to write file: %v", err)
	}

	fmt.Printf("File downloaded successfully to: %s\n", outputPath)
}

func listFiles(cmd *cobra.Command, args []string) {
	// Make request
	resp, err := http.Get(serverURL + "/api/v1/files")
	if err != nil {
		log.Fatalf("Failed to list files: %v", err)
	}
	defer resp.Body.Close()

	// Parse response
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Fatalf("Failed to parse response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("List failed: %v", result["error"])
	}

	files := result["files"].([]interface{})
	fmt.Printf("Found %v files:\n\n", result["count"])

	for _, fileInterface := range files {
		file := fileInterface.(map[string]interface{})
		fmt.Printf("ID: %s\n", file["id"])
		fmt.Printf("Name: %s\n", file["name"])
		fmt.Printf("Size: %v bytes\n", file["size"])
		fmt.Printf("Content Type: %s\n", file["content_type"])
		fmt.Printf("Created: %s\n", file["created_at"])
		fmt.Printf("Owner: %s\n", file["owner"])
		fmt.Println("---")
	}
}

func deleteFile(cmd *cobra.Command, args []string) {
	fileID := args[0]

	// Make request
	req, err := http.NewRequest("DELETE", serverURL+"/api/v1/files/"+fileID, nil)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to delete file: %v", err)
	}
	defer resp.Body.Close()

	// Parse response
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Fatalf("Failed to parse response: %v", err)
	}

	if resp.StatusCode == http.StatusOK {
		fmt.Printf("File deleted successfully: %s\n", result["message"])
	} else {
		fmt.Printf("Delete failed: %v\n", result["error"])
	}
}

func getFileInfo(cmd *cobra.Command, args []string) {
	fileID := args[0]

	// Make request
	resp, err := http.Get(serverURL + "/api/v1/files/" + fileID + "/info")
	if err != nil {
		log.Fatalf("Failed to get file info: %v", err)
	}
	defer resp.Body.Close()

	// Parse response
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Fatalf("Failed to parse response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Info failed: %v", result["error"])
	}

	// Pretty print JSON
	prettyJSON, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Fatalf("Failed to format JSON: %v", err)
	}

	fmt.Println(string(prettyJSON))
}
