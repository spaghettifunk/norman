/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

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
	fmt.Println("storage called")
}

func init() {
	rootCmd.AddCommand(storageCmd)
}
