package p2p

import (
	"crypto/sha256"
	"encoding/hex"
	"net"
	"strings"

	"github.com/ashtyn3/tinynamer/msg"
	"github.com/rs/zerolog/log"
)

type Peer struct {
	Address   string
	Ip        string
	Port      string
	Hash      string
	Sock      net.Conn
	developed bool
	Halt      bool
}

func NewPeer(sock net.Conn) *Peer {
	p := &Peer{}
	addr := sock.RemoteAddr().String()
	p.Sock = sock
	p.Halt = false
	p.Ip = strings.Split(addr, ":")[0]
	p.developed = false
	sha := sha256.New()
	sha.Write([]byte(p.Address + p.Port))
	p.Hash = hex.EncodeToString(sha.Sum(nil))

	return p
}

func (p *Peer) Send(msg *msg.Message) {
	_, err := p.Sock.Write(msg.Marshal())
	if err != nil {
		log.Fatal().Err(err).Send()
	}
}
