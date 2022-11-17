package main

import (
	"os"

	"github.com/ashtyn3/tinynamer/p2p"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	n := p2p.NewNode()
	n.Run()
}
