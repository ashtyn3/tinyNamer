package p2p

import (
	"bytes"

	"github.com/ashtyn3/tinynamer/msg"
	"github.com/rs/zerolog/log"
)

type Handlers struct {
	List map[string]func(*Peer, *msg.ProtoMessage)
	host *Node
}

func InitHandlers(n *Node) *Handlers {
	h := &Handlers{
		List: make(map[string]func(*Peer, *msg.ProtoMessage)),
		host: n,
	}
	h.List["get_peers"] = h.get_peers
	h.List["peers"] = h.peers

	return h
}

func (h *Handlers) peers(p *Peer, m *msg.ProtoMessage) {
	log.Error().Msg("Unimplemented")
}

func (h *Handlers) get_peers(p *Peer, m *msg.ProtoMessage) {
	h.host.Peers.PartialMarshal()
	b := bytes.Join(h.host.Peers.Buffer, []byte("\n"))
	p.Send(msg.Msg(h.host.Address, "peers", b))
}
