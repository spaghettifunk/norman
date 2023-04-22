/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// brokerCmd represents the broker command
var brokerCmd = &cobra.Command{
	Use:   "broker",
	Short: "A brief description of your command",
	Long:  ``,
	Run:   brokerRun,
}

func brokerRun(cmd *cobra.Command, args []string) {
	fmt.Println("broker called")
}

func init() {
	rootCmd.AddCommand(brokerCmd)
}
