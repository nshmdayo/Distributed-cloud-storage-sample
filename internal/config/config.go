// Package config handles application configuration
package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config holds the application configuration
type Config struct {
	Node       NodeConfig       `mapstructure:"node"`
	API        APIConfig        `mapstructure:"api"`
	Storage    StorageConfig    `mapstructure:"storage"`
	P2P        P2PConfig        `mapstructure:"p2p"`
	Crypto     CryptoConfig     `mapstructure:"crypto"`
	Blockchain BlockchainConfig `mapstructure:"blockchain"`
	Logging    LoggingConfig    `mapstructure:"logging"`
}

// NodeConfig contains node-specific configuration
type NodeConfig struct {
	ID         string `mapstructure:"id"`
	DataDir    string `mapstructure:"data_dir"`
	StorageDir string `mapstructure:"storage_dir"`
	MaxStorage int64  `mapstructure:"max_storage"`
	Replicas   int    `mapstructure:"replicas"`
	ChunkSize  int    `mapstructure:"chunk_size"`
}

// APIConfig contains API server configuration
type APIConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	TLS      bool   `mapstructure:"tls"`
	CertFile string `mapstructure:"cert_file"`
	KeyFile  string `mapstructure:"key_file"`
}

// StorageConfig contains storage-related configuration
type StorageConfig struct {
	Backend     string `mapstructure:"backend"`
	Path        string `mapstructure:"path"`
	MaxFileSize int64  `mapstructure:"max_file_size"`
	Compression bool   `mapstructure:"compression"`
}

// P2PConfig contains P2P network configuration
type P2PConfig struct {
	ListenAddr     string   `mapstructure:"listen_addr"`
	BootstrapPeers []string `mapstructure:"bootstrap_peers"`
	MaxPeers       int      `mapstructure:"max_peers"`
	PrivateKey     string   `mapstructure:"private_key"`
}

// CryptoConfig contains cryptographic configuration
type CryptoConfig struct {
	Algorithm   string `mapstructure:"algorithm"`
	KeySize     int    `mapstructure:"key_size"`
	EnableTLS   bool   `mapstructure:"enable_tls"`
	TLSCertPath string `mapstructure:"tls_cert_path"`
	TLSKeyPath  string `mapstructure:"tls_key_path"`
}

// BlockchainConfig contains blockchain-related configuration
type BlockchainConfig struct {
	Network         string `mapstructure:"network"`
	RPCEndpoint     string `mapstructure:"rpc_endpoint"`
	ContractAddress string `mapstructure:"contract_address"`
	PrivateKey      string `mapstructure:"private_key"`
	GasLimit        uint64 `mapstructure:"gas_limit"`
}

// LoggingConfig contains logging configuration
type LoggingConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
	Output string `mapstructure:"output"`
}

// DefaultConfig returns a configuration with default values
func DefaultConfig() *Config {
	homeDir, _ := os.UserHomeDir()
	dataDir := filepath.Join(homeDir, ".dcs")

	return &Config{
		Node: NodeConfig{
			DataDir:    dataDir,
			StorageDir: filepath.Join(dataDir, "storage"),
			MaxStorage: 10 * 1024 * 1024 * 1024, // 10GB
			Replicas:   3,
			ChunkSize:  1024 * 1024, // 1MB
		},
		API: APIConfig{
			Host: "localhost",
			Port: 8080,
			TLS:  false,
		},
		Storage: StorageConfig{
			Backend:     "filesystem",
			Path:        filepath.Join(dataDir, "files"),
			MaxFileSize: 100 * 1024 * 1024, // 100MB
			Compression: true,
		},
		P2P: P2PConfig{
			ListenAddr: "/ip4/0.0.0.0/tcp/4001",
			MaxPeers:   100,
		},
		Crypto: CryptoConfig{
			Algorithm: "AES-256-GCM",
			KeySize:   32,
			EnableTLS: true,
		},
		Blockchain: BlockchainConfig{
			Network:  "polygon-mumbai",
			GasLimit: 500000,
		},
		Logging: LoggingConfig{
			Level:  "info",
			Format: "json",
			Output: "stdout",
		},
	}
}

// LoadConfig loads configuration from file and environment variables
func LoadConfig(configPath string) (*Config, error) {
	config := DefaultConfig()

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	if configPath != "" {
		viper.SetConfigFile(configPath)
	} else {
		viper.AddConfigPath(".")
		viper.AddConfigPath("./config")
		viper.AddConfigPath("$HOME/.dcs")
	}

	// Environment variables
	viper.SetEnvPrefix("DCS")
	viper.AutomaticEnv()

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	// Unmarshal config
	if err := viper.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return config, nil
}

// Save saves the configuration to file
func (c *Config) Save(filepath string) error {
	viper.Set("node", c.Node)
	viper.Set("api", c.API)
	viper.Set("storage", c.Storage)
	viper.Set("p2p", c.P2P)
	viper.Set("crypto", c.Crypto)
	viper.Set("blockchain", c.Blockchain)
	viper.Set("logging", c.Logging)

	return viper.WriteConfigAs(filepath)
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.Node.ChunkSize <= 0 {
		return fmt.Errorf("invalid chunk size: %d", c.Node.ChunkSize)
	}

	if c.Node.Replicas <= 0 {
		return fmt.Errorf("invalid replicas count: %d", c.Node.Replicas)
	}

	if c.API.Port <= 0 || c.API.Port > 65535 {
		return fmt.Errorf("invalid API port: %d", c.API.Port)
	}

	if c.Storage.MaxFileSize <= 0 {
		return fmt.Errorf("invalid max file size: %d", c.Storage.MaxFileSize)
	}

	return nil
}
