package settings

import (
	"github.com/spf13/viper"
)

// App settings models
type StorageConfig struct {
	Type      string
	MaxSizeMB int64
	IPFS      struct {
		APIURL  string
		Gateway string
	}
}

type CryptoConfig struct {
	Enabled    bool
	KeyVersion string
	KeyEnv     string
}

type CompressionConfig struct {
	Enabled      bool   // Whether compression evidence is enabled.
	Algorithm    string // Compression algorithm: "gzip" (default), "btae".
	ModelPath    string // Reserved for BTAE model path.
	MinSizeMB    int64  // Minimum file size to trigger compression evidence.
	MaxSizeMB    int64  // Maximum file size to trigger compression evidence.
	MinSizeBytes int64  // Computed minimum bytes.
	MaxSizeBytes int64  // Computed maximum bytes.
}

type Config struct {
	Storage     StorageConfig
	Crypto      CryptoConfig
	Compression CompressionConfig
}

var Cfg Config

func Init() (err error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./settings")
	if err = viper.ReadInConfig(); err != nil {
		return err
	}

	Cfg.Storage.Type = viper.GetString("storage.type")
	Cfg.Storage.MaxSizeMB = viper.GetInt64("storage.max_size_mb")
	Cfg.Storage.IPFS.APIURL = viper.GetString("storage.ipfs.api_url")
	Cfg.Storage.IPFS.Gateway = viper.GetString("storage.ipfs.gateway_url")

	Cfg.Crypto.Enabled = viper.GetBool("crypto.enabled")
	Cfg.Crypto.KeyVersion = viper.GetString("crypto.key_version")
	Cfg.Crypto.KeyEnv = viper.GetString("crypto.key_env")

	Cfg.Compression.Enabled = viper.GetBool("compression.enabled")
	Cfg.Compression.Algorithm = viper.GetString("compression.algorithm")
	if Cfg.Compression.Algorithm == "" {
		Cfg.Compression.Algorithm = "gzip"
	}
	Cfg.Compression.ModelPath = viper.GetString("compression.model_path")
	Cfg.Compression.MinSizeMB = viper.GetInt64("compression.min_size_mb")
	Cfg.Compression.MaxSizeMB = viper.GetInt64("compression.max_size_mb")
	if Cfg.Compression.MinSizeMB <= 0 {
		Cfg.Compression.MinSizeMB = 5
	}
	if Cfg.Compression.MaxSizeMB <= 0 {
		Cfg.Compression.MaxSizeMB = 5 * 1024
	}
	if Cfg.Compression.MaxSizeMB < Cfg.Compression.MinSizeMB {
		Cfg.Compression.MaxSizeMB = Cfg.Compression.MinSizeMB
	}
	Cfg.Compression.MinSizeBytes = Cfg.Compression.MinSizeMB * 1024 * 1024
	Cfg.Compression.MaxSizeBytes = Cfg.Compression.MaxSizeMB * 1024 * 1024

	return nil
}
