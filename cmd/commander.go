package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"
	"github.com/spaghettifunk/norman/internal/commander"
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
	if normanCfg == nil {
		panic(fmt.Errorf("configuration has not loaded correctly"))
	}

	// initialize global logging
	logger.InitLogger(*normanCfg)

	// initialize service
	c, err := commander.New(*normanCfg)
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
		log.Info().Msg("Commander Server is down. Bye Bye!")
		close(done)
	}()

	// start http server goroutine
	go func() {
		addr := fmt.Sprintf("%s:%d", normanCfg.Commander.Address, normanCfg.Commander.Port)
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
