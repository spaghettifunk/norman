package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"
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
	logger.InitLogger("storage", *normanCfg)

	// initialize service
	st, err := storageserver.New(*normanCfg)
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
		if err := st.ShutdownServer(); err != nil {
			log.Fatal().Err(err)
			return
		}
		log.Info().Msg("Storage Server is down. Bye Bye!")
		close(done)
	}()

	go func() {
		addr := fmt.Sprintf("%s:%d", normanCfg.Storage.Address, normanCfg.Storage.Port)
		if err := st.StartServer(addr); err != nil {
			log.Fatal().Err(err)
			close(done)
		}
	}()

	// time to say goodbye
	<-done
}

func init() {
	rootCmd.AddCommand(storageCmd)
}
