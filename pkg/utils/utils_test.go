package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGenerateRandomID(t *testing.T) {
	id1, err := GenerateRandomID(32)
	if err != nil {
		t.Fatalf("Failed to generate ID: %v", err)
	}

	id2, err := GenerateRandomID(32)
	if err != nil {
		t.Fatalf("Failed to generate ID: %v", err)
	}

	// IDs should be different
	if id1 == id2 {
		t.Errorf("Expected different IDs, got %s and %s", id1, id2)
	}

	// ID should be correct length
	if len(id1) != 32 {
		t.Errorf("Expected length 32, got %d", len(id1))
	}
}

func TestEnsureDir(t *testing.T) {
	tmpDir := filepath.Join(os.TempDir(), "test-ensure-dir")
	defer os.RemoveAll(tmpDir)

	// Directory shouldn't exist initially
	if FileExists(tmpDir) {
		t.Errorf("Directory already exists: %s", tmpDir)
	}

	// Create directory
	if err := EnsureDir(tmpDir); err != nil {
		t.Fatalf("Failed to create directory: %v", err)
	}

	// Directory should exist now
	if !FileExists(tmpDir) {
		t.Errorf("Directory was not created: %s", tmpDir)
	}

	// Should not error if directory already exists
	if err := EnsureDir(tmpDir); err != nil {
		t.Errorf("EnsureDir failed on existing directory: %v", err)
	}
}

func TestSplitData(t *testing.T) {
	data := []byte("0123456789")

	// Test normal split
	chunks := SplitData(data, 3)
	expected := [][]byte{
		[]byte("012"),
		[]byte("345"),
		[]byte("678"),
		[]byte("9"),
	}

	if len(chunks) != len(expected) {
		t.Errorf("Expected %d chunks, got %d", len(expected), len(chunks))
	}

	for i, chunk := range chunks {
		if string(chunk) != string(expected[i]) {
			t.Errorf("Chunk %d: expected %s, got %s", i, expected[i], chunk)
		}
	}

	// Test edge cases
	singleChunk := SplitData(data, 20)
	if len(singleChunk) != 1 {
		t.Errorf("Expected 1 chunk for large chunk size, got %d", len(singleChunk))
	}

	zeroChunk := SplitData(data, 0)
	if len(zeroChunk) != 1 {
		t.Errorf("Expected 1 chunk for zero chunk size, got %d", len(zeroChunk))
	}
}

func TestJoinChunks(t *testing.T) {
	chunks := [][]byte{
		[]byte("Hello"),
		[]byte(" "),
		[]byte("World"),
		[]byte("!"),
	}

	result := JoinChunks(chunks)
	expected := "Hello World!"

	if string(result) != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func TestGetStoragePath(t *testing.T) {
	baseDir := "/storage"
	fileID := "abcdef1234567890"

	path := GetStoragePath(baseDir, fileID)
	expected := filepath.Join(baseDir, "ab", fileID)

	if path != expected {
		t.Errorf("Expected %s, got %s", expected, path)
	}
}

func TestFormatBytes(t *testing.T) {
	tests := []struct {
		bytes    int64
		expected string
	}{
		{0, "0 B"},
		{512, "512 B"},
		{1024, "1.0 KB"},
		{1536, "1.5 KB"},
		{1048576, "1.0 MB"},
		{1073741824, "1.0 GB"},
	}

	for _, test := range tests {
		result := FormatBytes(test.bytes)
		if result != test.expected {
			t.Errorf("FormatBytes(%d): expected %s, got %s", test.bytes, test.expected, result)
		}
	}
}

func TestValidateFileID(t *testing.T) {
	// Valid file ID (64 hex characters)
	validID := "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"
	if !ValidateFileID(validID) {
		t.Errorf("Expected valid ID to pass validation: %s", validID)
	}

	// Invalid IDs
	invalidIDs := []string{
		"",      // Empty
		"short", // Too short
		"1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdefg", // Too long
		"1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdeg",  // Invalid hex
	}

	for _, id := range invalidIDs {
		if ValidateFileID(id) {
			t.Errorf("Expected invalid ID to fail validation: %s", id)
		}
	}
}
