package configuration

import (
	"github.com/spf13/viper"
)

type Configuration struct {
	ConfigDir string     `toml:"config_dir"`
	Commander *commander `toml:"commander"`
	Broker    *broker    `toml:"broker"`
	Storage   *storage   `toml:"storage"`
	Logger    *logger    `toml:"logger"`
}

type commander struct {
	Address string `toml:"address"`
	Port    int    `toml:"port"`
	Aqua    *aqua  `toml:"aqua"`
}

type aqua struct {
	CAFile               string `toml:"cafile"`
	ServerCertFile       string `toml:"server_cert_file"`
	ServerKeyFile        string `toml:"server_key_file"`
	RootClientCertFile   string `toml:"root_client_cert_file"`
	RootClientKeyFile    string `toml:"root_client_key_file"`
	NobodyClientCertFile string `toml:"nobody_client_cert_file"`
	NobodyClientKeyFile  string `toml:"nobody_client_key_file"`
}

type broker struct {
	Address string `toml:"address"`
	Port    int    `toml:"port"`
}

type storage struct {
	Address     string       `toml:"address"`
	Port        int          `toml:"port"`
	DeepStorage *deepStorage `toml:"deep_storage"`
}

type deepStorage struct {
	Type   string `toml:"type"`
	Bucket string `toml:"bucket"`
}

type logger struct {
	Level  string `toml:"level"`
	Pretty bool   `toml:"pretty"`
}

func Fetch() *Configuration {
	return &Configuration{
		ConfigDir: getStringOrDefault("config_dir", "/Users/davideberdin/Documents/github/norman/test/"),
		Commander: &commander{
			Address: getStringOrDefault("commander.address", "127.0.0.1"),
			Port:    getIntOrDefault("commander.port", 8080),
			Aqua: &aqua{
				CAFile:               getStringOrDefault("commander.aqua.cafile", "certs/ca.pem"),
				ServerCertFile:       getStringOrDefault("commander.aqua.server_cert_file", "certs/server.pem"),
				ServerKeyFile:        getStringOrDefault("commander.aqua.server_key_file", "certs/server-key.pem"),
				RootClientCertFile:   getStringOrDefault("commander.aqua.root_client_cert_file", "certs/root-client.pem"),
				RootClientKeyFile:    getStringOrDefault("commander.aqua.root_client_key_file", "certs/root-client-key.pem"),
				NobodyClientCertFile: getStringOrDefault("commander.aqua.nobody_client_cert_file", "certs/nobody-client.pem"),
				NobodyClientKeyFile:  getStringOrDefault("commander.aqua.nobody_client_key_file", "certs/nobody-client-key.pem"),
			},
		},
		Broker: &broker{
			Address: getStringOrDefault("broker.address", "127.0.0.1"),
			Port:    getIntOrDefault("broker.port", 8081),
		},
		Storage: &storage{
			Address: getStringOrDefault("storage.address", "127.0.0.1"),
			Port:    getIntOrDefault("storage.port", 8082),
			DeepStorage: &deepStorage{
				Type:   getStringOrDefault("storage.deep_storage", "s3"),
				Bucket: getStringOrDefault("storage.deep_storage", "test-deep-storage"),
			},
		},
		Logger: &logger{
			Level:  getStringOrDefault("logger.level", "info"),
			Pretty: getBoolOrDefault("logger.pretty", false),
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
