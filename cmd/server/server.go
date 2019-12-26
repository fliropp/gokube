package server

import (
	"os"
	"os/signal"

	"github.com/amedia/cloudordino/pkg/logging"
	"github.com/fliropp/aresworld/pkg/web"
	"github.com/spf13/cobra"
	//"github.com/spf13/viper"
)

var ServerCmd = &cobra.Command{
	Use:   "server",
	Short: "the server",
	Run: func(cmd *cobra.Command, args []string) {
		log := logging.GetLogger()
		log.Info("Start server . . . ")
		webserver := web.NewWebServer(log)
		webserver.Start()
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, os.Interrupt)
	Serverloop:
		for {
			select {
			case <-signals:
				break Serverloop
			}
		}
	},
}
