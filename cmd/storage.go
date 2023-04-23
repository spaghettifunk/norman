/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	configuration "github.com/spaghettifunk/norman/internal/common"
	storageserver "github.com/spaghettifunk/norman/internal/storage"
	"github.com/spaghettifunk/norman/pkg/logger"
	"github.com/spf13/cobra"
)

// storageCmd represents the storage command
var storageCmd = &cobra.Command{
	Use:   "storage",
	Short: "A brief description of your command",
	Long:  ``,
	Run:   storageRun,
}

func storageRun(cmd *cobra.Command, args []string) {
	// fetch and validate configuration file
	config := configuration.Fetch()
	if err := config.Validate(); err != nil {
		panic(err.Error())
	}

	// initialize global logging
	logger.InitLogger(*config)

	// initialize service
	_ = storageserver.New(*config)
}

func init() {
	rootCmd.AddCommand(storageCmd)
}
