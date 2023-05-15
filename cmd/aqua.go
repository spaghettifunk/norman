/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/fsnotify/fsnotify"
	"github.com/spaghettifunk/norman/aqua/agent"
	"github.com/spaghettifunk/norman/aqua/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var aquaCfg *config.Configuration

// aquaCmd represents the aqua command
var aquaCmd = &cobra.Command{
	Use:     "aqua",
	Short:   "It runs the Aqua service",
	Long:    ``,
	PreRunE: setupConfig,
	Run:     aquaRun,
}

func aquaRun(cmd *cobra.Command, args []string) {
	if aquaCfg == nil {
		panic(fmt.Errorf("configuration has not loaded correctly"))
	}

	var err error
	agent, err := agent.New(aquaCfg.Config)
	if err != nil {
		panic(err)
	}
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
	<-sigc

	if err := agent.Shutdown(); err != nil {
		panic(err)
	}
}

func setupConfig(cmd *cobra.Command, args []string) error {
	viper.SetConfigName("aqua.config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		if e, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			fmt.Println("aqua.config.toml not found")
		} else {
			// Config file was found but another error was produced
			fmt.Println(e.Error())
		}
	}

	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})

	viper.WatchConfig()

	// fetch and validate configuration file
	aquaCfg = config.Fetch()
	if err := aquaCfg.Validate(); err != nil {
		panic(err.Error())
	}

	return nil
}

func init() {
	rootCmd.AddCommand(aquaCmd)
}
