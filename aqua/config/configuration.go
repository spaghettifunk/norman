package config

import (
	"fmt"
	"os"
	"path"

	"github.com/spaghettifunk/norman/aqua/agent"
	"github.com/spf13/viper"
)

type Configuration struct {
	agent.Config
	ConfigDir string `toml:"config_dir" description:"Configuration directory where the certificates are"`
	Bootstrap bool   `toml:"bootstrap" description:"Bootstrapping the cluster (true) or if the node is joining (false)"`

	ServerTLSConfig TLSConfig `toml:"server_tls"`
	PeerTLSConfig   TLSConfig `toml:"peer_tls"`
}

type TLSConfig struct {
	CertFile      string `toml:"cert_file" description:"Path to server tls cert."`
	KeyFile       string `toml:"key_file" description:"Path to server tls key."`
	CAFile        string `toml:"ca_file" description:"Path to server certificate authority."`
	ServerAddress string `description:"Not exposed in configuration -- used for test"`
	Server        bool   `description:"Not exposed in configuration -- used for test"`
}

func Fetch() *Configuration {
	// TODO: change me
	dataDir := path.Join(os.TempDir(), "aqua/")
	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		err := os.Mkdir(dataDir, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}

	// get current hostname from the node itself
	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	cfg := &Configuration{
		ConfigDir: getStringOrDefault("config_dir", "/Users/davideberdin/Documents/github/norman/test/"),
		Config: agent.Config{
			DataDir:        getStringOrDefault("data_dir", dataDir),
			NodeName:       getStringOrDefault("node_name", hostname),
			BindAddr:       getStringOrDefault("bind_addr", "127.0.0.1:8401"),
			RPCPort:        getIntOrDefault("rpc_port", 8400),
			StartJoinAddrs: getStringArrayOrDefault("start_join_addrs", nil),
		},
		Bootstrap: getBoolOrDefault("bootstrap", false),
	}

	// set up defaults for TLS server
	cfg.ServerTLSConfig = TLSConfig{
		CertFile: fmt.Sprintf("%s%s", cfg.ConfigDir, getStringOrDefault("server_tls.cert_file", "certs/server.pem")),
		KeyFile:  fmt.Sprintf("%s%s", cfg.ConfigDir, getStringOrDefault("server_tls.key_file", "certs/server-key.pem")),
		CAFile:   fmt.Sprintf("%s%s", cfg.ConfigDir, getStringOrDefault("server_tls.ca_file", "certs/ca.pem")),
	}
	// set up defaults for TLS peer
	cfg.PeerTLSConfig = TLSConfig{
		CertFile: fmt.Sprintf("%s%s", cfg.ConfigDir, getStringOrDefault("peer_tls.cert_file", "certs/client.pem")),
		KeyFile:  fmt.Sprintf("%s%s", cfg.ConfigDir, getStringOrDefault("peer_tls.key_file", "certs/client-key.pem")),
		CAFile:   fmt.Sprintf("%s%s", cfg.ConfigDir, getStringOrDefault("peer_tls.ca_file", "certs/ca.pem")),
	}

	// setup TLS server
	if cfg.ServerTLSConfig.CertFile != "" &&
		cfg.ServerTLSConfig.KeyFile != "" {
		cfg.ServerTLSConfig.Server = true
		cfg.Config.ServerTLSConfig, err = SetupTLSConfig(cfg.ServerTLSConfig)
		if err != nil {
			panic(err)
		}
	}

	// setup TLS peer
	if cfg.PeerTLSConfig.CertFile != "" &&
		cfg.PeerTLSConfig.KeyFile != "" {
		cfg.Config.PeerTLSConfig, err = SetupTLSConfig(cfg.PeerTLSConfig)
		if err != nil {
			panic(err)
		}
	}

	return cfg
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

func getStringArrayOrDefault(key string, defaultValue []string) []string {
	value := viper.GetStringSlice(key)
	if len(value) > 0 {
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
