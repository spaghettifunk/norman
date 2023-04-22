/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// commanderCmd represents the commander command
var commanderCmd = &cobra.Command{
	Use:   "commander",
	Short: "A brief description of your command",
	Long:  ``,
	Run:   commanderRun,
}

func commanderRun(cmd *cobra.Command, args []string) {
	fmt.Println("commander called")
}

func init() {
	rootCmd.AddCommand(commanderCmd)
}
