package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	gw "github.com/orus-io/nats-websocket-gw"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = cobra.Command{
	Run: rootCmdRun,
}

func init() {
	flags := rootCmd.Flags()

	flags.Int("port", 8910, "Port to run http server on")
	flags.String("host", "localhost", "host/IP to run http server on")
	flags.String("path", "/nats", "webpath of the websocket")
	flags.String("nats", "localhost:4222", "nats server address:port")
	flags.Bool("no-origin-check", false, "Disable websocket origin check")
	flags.Bool("trace", false, "Enable trace logs")

	viper.BindPFlag("port", flags.Lookup("port"))
	viper.BindPFlag("host", flags.Lookup("host"))
	viper.BindPFlag("path", flags.Lookup("path"))
	viper.BindPFlag("nats", flags.Lookup("nats"))
	viper.BindPFlag("no-origin-check", flags.Lookup("no-origin-check"))
	viper.BindPFlag("trace", flags.Lookup("trace"))
}

func rootCmdRun(cmd *cobra.Command, args []string) {

	settings := gw.Settings{
		NatsAddr: viper.GetString("nats"),
	}

	if viper.GetBool("no-origin-check") {
		settings.WSUpgrader = &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     func(r *http.Request) bool { return true },
		}
	}
	if viper.GetBool("trace") {
		settings.Trace = true
	}

	listenOn := viper.GetString("host") + ":" + viper.GetString("port")

	gateway := gw.NewGateway(settings)
	http.HandleFunc(viper.GetString("path"), gateway.Handler)
	http.ListenAndServe(listenOn, nil)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
