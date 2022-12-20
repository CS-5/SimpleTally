package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/urfave/cli/v2"
)

var (
	address string
	port    int64
	source  string
	debug   bool

	listenAddress = "localhost"
	listenPort    = 8099
)

func init() {
	app := &cli.App{
		Name:        "SimpleTally",
		Description: "A very simple tally for OBS Studio",
		Usage:       "Tally... simple tally",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "address",
				Value:    "127.0.0.1",
				Usage:    "OBS Studio IP address",
				Aliases:  []string{"a"},
				Required: false,
			},
			&cli.Int64Flag{
				Name:     "port",
				Value:    4455,
				Usage:    "obs-websocket port",
				Aliases:  []string{"p"},
				Required: false,
			},
			&cli.StringFlag{
				Name:     "source",
				Usage:    "the name of the source to monitor",
				Aliases:  []string{"s"},
				Required: true,
			},
			&cli.BoolFlag{
				Name:    "debug",
				Aliases: []string{"d"},
				Value:   false,
			},
		},
		Action: func(ctx *cli.Context) error {
			address = ctx.String("address")
			port = ctx.Int64("port")
			source = ctx.String("source")
			debug = ctx.Bool("debug")

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}

func main() {
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", listenAddress, listenPort))
	if err != nil {
		fmt.Printf("Unable to start TCP listener on port: %d: %s\n", listenPort, err.Error())
	}

	defer l.Close()

	log.Default().Printf("Listening on %s:%d", listenAddress, listenPort)
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
		}

		go handle(conn)
	}
}

func handle(conn net.Conn) {

}
