package cmd

import (
	"Trading-Engine/internal/config"
	"Trading-Engine/internal/handler"
	"Trading-Engine/internal/server"
	"Trading-Engine/internal/storage/mysql"
	"Trading-Engine/pkg/logger"
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "",
	Long:  ``,
	Run:   RunServer,
}

func init() {
	rootCmd.AddCommand(serverCmd)
}

func RunServer(cmd *cobra.Command, args []string) {
	defer cmdRecover()
	logger.SetLogger(local)
	config := config.GetConfig(configFile)
	db := mysql.NewDatabase(config.Mysql)
	handler := handler.NewHandler(db)
	svr := server.NewServer(config.Http, handler)

	svr.Run()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	svr.Shutdown(ctx)
	db.Shutdown(ctx)
	log.Info().Msg("shutting down")
	os.Exit(0)

}
