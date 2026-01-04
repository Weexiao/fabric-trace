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

type Config struct {
	Storage StorageConfig
	Crypto  CryptoConfig
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

	return nil
}
