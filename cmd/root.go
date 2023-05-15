/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/fsnotify/fsnotify"
	configuration "github.com/spaghettifunk/norman/internal/common"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var normanCfg *configuration.Configuration

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:               "norman",
	Short:             "A brief description of your application",
	PersistentPreRunE: setupConfiguration,
	Long:              ``,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func setupConfiguration(cmd *cobra.Command, args []string) error {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		if e, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			fmt.Println("config.toml not found")
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
	normanCfg = configuration.Fetch()
	if err := normanCfg.Validate(); err != nil {
		panic(err.Error())
	}
	return nil
}
