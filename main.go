package main

import (
	"SimpleTally/vmix"
	"fmt"
	"os"
	"time"

	"github.com/andreykaipov/goobs"
	"github.com/andreykaipov/goobs/api/events"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

var (
	// OBS Websocket connection settings
	address  string
	port     int64
	password string
	source   string

	// vMix Listener bind address and port
	listenAddress = "0.0.0.0"
	listenPort    = 8099

	debug bool

	oLog zerolog.Logger
	vLog zerolog.Logger
)

func init() {
	// Init CLI handler
	app := &cli.App{
		Name:        "SimpleTally",
		Description: "A very simple tally for OBS Studio",
		Usage:       "Tally... simple tally",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "address",
				Value:    "127.0.0.1",
				Usage:    "obs-websocket address",
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
				Name:     "password",
				Usage:    "obs-websocket password",
				Aliases:  []string{"pw"},
				Required: true,
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

		// Set globals to values from flags
		Action: func(ctx *cli.Context) error {
			address = ctx.String("address")
			port = ctx.Int64("port")
			source = ctx.String("source")
			password = ctx.String("password")

			debug = ctx.Bool("debug")

			return nil
		},
	}

	// Parse flags
	if err := app.Run(os.Args); err != nil {
		panic(err)
	}

	// Init logger
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC822})

	oLog = log.With().Str("component", "OBS").Logger()
	vLog = log.With().Str("component", "vMix").Logger()
}

func main() {
	log.Info().Msg("Init finished, starting listeners...")

	// Listen for vMix TCP connections
	vmix := vmix.New(listenAddress, listenPort, vLog)
	go vmix.Listen()

	vLog.Info().
		Msgf("vMix TCP listener started on '%s:%d'", listenAddress, listenPort)

	// Start OBS websocket connection
	obs, err := goobs.New(fmt.Sprintf("%s:%d", address, port), goobs.WithPassword(password))
	if err != nil {
		oLog.Fatal().
			AnErr("error", err).
			Msgf("Unable to connect to OBS at '%s:%d' with password '%s'", address, port, password)
	}
	defer obs.Disconnect()

	// Get info from OBS to verify the connection is working
	info, err := obs.General.GetVersion()
	if err != nil {
		oLog.Error().
			AnErr("error", err).
			Msg("Unable to get OBS version information, this could mean there's a problem with the connection")
	}
	oLog.Info().Msgf("Connected to OBS (v%s) at '%s:%d'", info.ObsVersion, address, port)

	// Listen and handle OBS events
	oLog.Info().Msg("Listening for OBS websocket events")
	obs.Listen(obsEvent)
}

func obsEvent(event any) {
	switch e := event.(type) {

	// Is the source currently in program?
	case *events.InputActiveStateChanged:
		if e.InputName == source {

		}

		break

	// Is the source currently in preview?
	case *events.InputShowStateChanged:
		if e.InputName == source {

		}

		break
	}
}
