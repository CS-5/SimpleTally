package vmix

import (
	"fmt"

	"github.com/rs/zerolog"
	"github.com/tidwall/evio"
)

type VMix struct {
	address string
	port    int

	log zerolog.Logger

	listeners map[string]evio.Conn
}

func New(listenAddress string, listenPort int, logger zerolog.Logger) *VMix {
	return &VMix{
		address: listenAddress,
		port:    listenPort,
		log:     logger,
	}
}

// Listen is a blocking function that starts listening for TCP connections from
// clients on the configured interface
func (vm *VMix) Listen() {
	err := evio.Serve(evio.Events{
		Opened: vm.opened,
		Closed: vm.closed,
		Data:   vm.data,
	}, fmt.Sprintf("tcp://%s:%d", vm.address, vm.port))
	if err != nil {
		vm.log.Fatal().
			AnErr("error", err).
			Msgf("Unable to start vMix TCP listener on '%s:%d'", vm.address, vm.port)
	}
}

func (vm *VMix) opened(c evio.Conn) (out []byte, opts evio.Options, action evio.Action) {
	vm.log.Info().Msgf("New connection from '%s' accepted", c.RemoteAddr().String())

	// Add connection as listener
	vm.listeners[c.RemoteAddr().String()] = c

	return
}

func (vm *VMix) closed(c evio.Conn, err error) (acton evio.Action) {
	vm.log.Info().Msgf("Connection from '%s' closed", c.RemoteAddr().String())

	// Remove connection from listeners
	delete(vm.listeners, c.RemoteAddr().String())

	return
}

func (vm *VMix) data(c evio.Conn, in []byte) (out []byte, action evio.Action) {
	vm.log.Debug().Str("payload", string(in)).Msg("Payload received")

	return
}
