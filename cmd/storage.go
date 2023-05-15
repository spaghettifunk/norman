/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

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
	if normanCfg == nil {
		panic(fmt.Errorf("configuration has not loaded correctly"))
	}

	// initialize global logging
	logger.InitLogger(*normanCfg)

	// initialize service
	_ = storageserver.New(*normanCfg)
}

func init() {
	rootCmd.AddCommand(storageCmd)
}
