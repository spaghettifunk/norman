package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"
	"github.com/spaghettifunk/norman/internal/commander"
	configuration "github.com/spaghettifunk/norman/internal/common"
	"github.com/spaghettifunk/norman/pkg/logger"
	"github.com/spf13/cobra"
)

// commanderCmd represents the commander command
var commanderCmd = &cobra.Command{
	Use:   "commander",
	Short: "",
	Long:  ``,
	Run:   commanderRun,
}

func commanderRun(cmd *cobra.Command, args []string) {
	// fetch and validate configuration file
	config := configuration.Fetch()
	if err := config.Validate(); err != nil {
		panic(err.Error())
	}

	// initialize global logging
	logger.InitLogger(*config)

	// initialize service
	c := commander.New(*config)

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

		// shutdown aqua service

		log.Info().Msg("Commander Server is down. Bye Bye!")
		close(done)
	}()

	// start aqua service goroutine
	go func() {

	}()

	// start http server goroutine
	go func() {
		addr := fmt.Sprintf("%s:%d", config.Commander.Address, config.Commander.Port)
		if err := c.StartServer(addr); err != nil {
			log.Fatal().Err(err)
			close(done)
		}
	}()

	// time to say goodbye
	<-done
}

func init() {
	rootCmd.AddCommand(commanderCmd)
}
