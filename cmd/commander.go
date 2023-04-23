package cmd

import (
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
		if err := c.ShutdownServer(); err != nil {
			log.Fatal().Err(err)
			return
		}
		log.Info().Msg("Commander Server is down. Bye Bye!")
		close(done)
	}()

	// start http server goroutine
	go func() {
		if err := c.StartServer(); err != nil {
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
