package configuration

import (
	"crypto/tls"

	"github.com/spf13/viper"
)

type Configuration struct {
	Commander *commander `toml:"commander"`
	Broker    *broker    `toml:"broker"`
	Storage   *storage   `toml:"storage"`
	Logger    *logger    `toml:"logger"`
}

type commander struct {
	Address string `toml:"address"`
	Port    int    `toml:"port"`
}

type broker struct {
	Address string `toml:"address"`
	Port    int    `toml:"port"`
}

type storage struct {
	Address         string `toml:"address"`
	Port            int    `toml:"port"`
	BindAddr        string `toml:"bind_addr" description:"Address where binding the gRPC server to"`
	RPCPort         int    `toml:"rpc_port" description:"Port for RPC clients connections."`
	ServerTLSConfig *tls.Config
	PeerTLSConfig   *tls.Config
	DeepStorage     *deepStorage `toml:"deep_storage"`
}

type deepStorage struct {
	Type   string `toml:"type"`
	Bucket string `toml:"bucket"`
}

type logger struct {
	Level string `toml:"level"`
}

func Fetch() *Configuration {
	return &Configuration{
		Commander: &commander{
			Address: getStringOrDefault("commander.address", "127.0.0.1"),
			Port:    getIntOrDefault("commander.port", 8080),
		},
		Broker: &broker{
			Address: getStringOrDefault("broker.address", "127.0.0.1"),
			Port:    getIntOrDefault("broker.port", 8081),
		},
		Storage: &storage{
			Address:  getStringOrDefault("storage.address", "127.0.0.1"),
			Port:     getIntOrDefault("storage.port", 8082),
			BindAddr: getStringOrDefault("bind_addr", "127.0.0.1:8401"),
			RPCPort:  getIntOrDefault("rpc_port", 8400),
		},
		Logger: &logger{
			Level: getStringOrDefault("logger.level", "info"),
		},
	}
}

func (c *Configuration) Validate() error {
	return nil
}

func getStringOrDefault(key string, defaultValue string) string {
	value := viper.GetString(key)
	if value != "" {
		return value
	}
	return defaultValue
}

func getIntOrDefault(key string, defaultValue int) int {
	value := viper.GetInt(key)
	if value != 0 {
		return value
	}
	return defaultValue
}

func getBoolOrDefault(key string, defaultValue bool) bool {
	value := viper.GetBool(key)
	if !value {
		return value
	}
	return defaultValue
}
