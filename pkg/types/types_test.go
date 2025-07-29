package types

import (
	"testing"
	"time"
)

func TestGenerateFileID(t *testing.T) {
	name := "test.txt"
	content := []byte("test content")

	id1 := GenerateFileID(name, content)
	id2 := GenerateFileID(name, content)

	// Same input should generate same ID
	if id1 != id2 {
		t.Errorf("Expected same ID for same input, got %s and %s", id1, id2)
	}

	// ID should be 64 characters (SHA256 hex)
	if len(id1) != 64 {
		t.Errorf("Expected ID length 64, got %d", len(id1))
	}

	// Different content should generate different ID
	id3 := GenerateFileID(name, []byte("different content"))
	if id1 == id3 {
		t.Errorf("Expected different ID for different content")
	}
}

func TestGenerateChunkID(t *testing.T) {
	fileID := "test-file-id"
	index := 0
	content := []byte("chunk content")

	id1 := GenerateChunkID(fileID, index, content)
	id2 := GenerateChunkID(fileID, index, content)

	// Same input should generate same ID
	if id1 != id2 {
		t.Errorf("Expected same ID for same input, got %s and %s", id1, id2)
	}

	// Different index should generate different ID
	id3 := GenerateChunkID(fileID, 1, content)
	if id1 == id3 {
		t.Errorf("Expected different ID for different index")
	}
}

func TestCalculateHash(t *testing.T) {
	data := []byte("test data")

	hash1 := CalculateHash(data)
	hash2 := CalculateHash(data)

	// Same data should generate same hash
	if hash1 != hash2 {
		t.Errorf("Expected same hash for same data, got %s and %s", hash1, hash2)
	}

	// Hash should be 64 characters (SHA256 hex)
	if len(hash1) != 64 {
		t.Errorf("Expected hash length 64, got %d", len(hash1))
	}

	// Different data should generate different hash
	hash3 := CalculateHash([]byte("different data"))
	if hash1 == hash3 {
		t.Errorf("Expected different hash for different data")
	}
}

func TestNodeStatus(t *testing.T) {
	tests := []struct {
		status   NodeStatus
		expected string
	}{
		{NodeStatusOnline, "online"},
		{NodeStatusOffline, "offline"},
		{NodeStatusSuspended, "suspended"},
		{NodeStatusMaintenance, "maintenance"},
		{NodeStatus(999), "unknown"},
	}

	for _, test := range tests {
		if got := test.status.String(); got != test.expected {
			t.Errorf("Expected %s, got %s", test.expected, got)
		}
	}
}

func TestMessageType(t *testing.T) {
	tests := []struct {
		msgType  MessageType
		expected string
	}{
		{MessageTypeFileRequest, "file_request"},
		{MessageTypeFileResponse, "file_response"},
		{MessageTypeChunkRequest, "chunk_request"},
		{MessageTypeChunkResponse, "chunk_response"},
		{MessageTypeNodeAnnouncement, "node_announcement"},
		{MessageTypeHeartbeat, "heartbeat"},
		{MessageTypeSyncRequest, "sync_request"},
		{MessageTypeSyncResponse, "sync_response"},
		{MessageType(999), "unknown"},
	}

	for _, test := range tests {
		if got := test.msgType.String(); got != test.expected {
			t.Errorf("Expected %s, got %s", test.expected, got)
		}
	}
}

func TestFileInfo(t *testing.T) {
	fileInfo := &FileInfo{
		ID:          "test-id",
		Name:        "test.txt",
		Size:        1024,
		Hash:        "test-hash",
		ContentType: "text/plain",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Owner:       "test-owner",
		Replicas:    3,
		IsEncrypted: true,
	}

	if fileInfo.ID != "test-id" {
		t.Errorf("Expected ID test-id, got %s", fileInfo.ID)
	}

	if fileInfo.Size != 1024 {
		t.Errorf("Expected size 1024, got %d", fileInfo.Size)
	}

	if !fileInfo.IsEncrypted {
		t.Errorf("Expected file to be encrypted")
	}
}
