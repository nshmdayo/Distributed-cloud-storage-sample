// Package types contains common data structures used across the application
package types

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
)

// FileInfo represents metadata about a stored file
type FileInfo struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Size        int64       `json:"size"`
	Hash        string      `json:"hash"`
	ContentType string      `json:"content_type"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
	Owner       string      `json:"owner"`
	Chunks      []ChunkInfo `json:"chunks"`
	Replicas    int         `json:"replicas"`
	IsEncrypted bool        `json:"is_encrypted"`
}

// ChunkInfo represents a chunk of a file
type ChunkInfo struct {
	ID       string   `json:"id"`
	Index    int      `json:"index"`
	Size     int64    `json:"size"`
	Hash     string   `json:"hash"`
	NodeIDs  []string `json:"node_ids"`
	Checksum string   `json:"checksum"`
}

// NodeInfo represents information about a storage node
type NodeInfo struct {
	ID           string     `json:"id"`
	Address      string     `json:"address"`
	Port         int        `json:"port"`
	PublicKey    string     `json:"public_key"`
	StorageUsed  int64      `json:"storage_used"`
	StorageTotal int64      `json:"storage_total"`
	Status       NodeStatus `json:"status"`
	LastSeen     time.Time  `json:"last_seen"`
	Reputation   float64    `json:"reputation"`
}

// NodeStatus represents the status of a node
type NodeStatus int

const (
	NodeStatusOnline NodeStatus = iota
	NodeStatusOffline
	NodeStatusSuspended
	NodeStatusMaintenance
)

func (s NodeStatus) String() string {
	switch s {
	case NodeStatusOnline:
		return "online"
	case NodeStatusOffline:
		return "offline"
	case NodeStatusSuspended:
		return "suspended"
	case NodeStatusMaintenance:
		return "maintenance"
	default:
		return "unknown"
	}
}

// NetworkMessage represents a message in the P2P network
type NetworkMessage struct {
	Type      MessageType `json:"type"`
	From      string      `json:"from"`
	To        string      `json:"to"`
	Timestamp time.Time   `json:"timestamp"`
	Data      interface{} `json:"data"`
	Signature string      `json:"signature"`
}

// MessageType defines the type of network message
type MessageType int

const (
	MessageTypeFileRequest MessageType = iota
	MessageTypeFileResponse
	MessageTypeChunkRequest
	MessageTypeChunkResponse
	MessageTypeNodeAnnouncement
	MessageTypeHeartbeat
	MessageTypeSyncRequest
	MessageTypeSyncResponse
)

func (m MessageType) String() string {
	switch m {
	case MessageTypeFileRequest:
		return "file_request"
	case MessageTypeFileResponse:
		return "file_response"
	case MessageTypeChunkRequest:
		return "chunk_request"
	case MessageTypeChunkResponse:
		return "chunk_response"
	case MessageTypeNodeAnnouncement:
		return "node_announcement"
	case MessageTypeHeartbeat:
		return "heartbeat"
	case MessageTypeSyncRequest:
		return "sync_request"
	case MessageTypeSyncResponse:
		return "sync_response"
	default:
		return "unknown"
	}
}

// GenerateFileID generates a unique ID for a file
func GenerateFileID(name string, content []byte) string {
	hasher := sha256.New()
	hasher.Write([]byte(name))
	hasher.Write(content)
	return hex.EncodeToString(hasher.Sum(nil))
}

// GenerateChunkID generates a unique ID for a chunk
func GenerateChunkID(fileID string, index int, content []byte) string {
	hasher := sha256.New()
	hasher.Write([]byte(fileID))
	hasher.Write([]byte{byte(index)})
	hasher.Write(content)
	return hex.EncodeToString(hasher.Sum(nil))
}

// CalculateHash calculates SHA256 hash of data
func CalculateHash(data []byte) string {
	hasher := sha256.New()
	hasher.Write(data)
	return hex.EncodeToString(hasher.Sum(nil))
}
