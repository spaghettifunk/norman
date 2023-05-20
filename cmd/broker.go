/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"
	"github.com/spaghettifunk/norman/internal/broker"
	"github.com/spaghettifunk/norman/pkg/logger"
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
	if normanCfg == nil {
		panic(fmt.Errorf("configuration has not loaded correctly"))
	}

	// initialize global logging
	logger.InitLogger(*normanCfg)

	// initialize service
	c, err := broker.New(*normanCfg)
	if err != nil {
		panic(err.Error())
	}

	// signal channel to capture system calls
	done := make(chan bool, 1)
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	// start shutdown goroutine
	go func() {
		// capture sigterm and other system call here
		<-sigCh
		// shutdown http server
		if err := c.ShutdownServer(); err != nil {
			log.Fatal().Err(err)
			return
		}
		log.Info().Msg("Broker is down. Bye Bye!")
		close(done)
	}()

	// start http server goroutine
	go func() {
		addr := fmt.Sprintf("%s:%d", normanCfg.Broker.Address, normanCfg.Broker.Port)
		if err := c.StartServer(addr); err != nil {
			log.Fatal().Err(err)
			close(done)
		}
	}()

	// time to say goodbye
	<-done
}

func init() {
	rootCmd.AddCommand(brokerCmd)
}
