package p2p

import (
	"encoding/hex"
	"net"
	"strings"

	"github.com/ashtyn3/tinynamer/msg"
	"github.com/rs/zerolog/log"
)

type Handlers struct {
	List map[string]func(*Peer, *msg.ProtoMessage, *Handlers)
	Host *Node
}

func InitHandlers(n *Node) *Handlers {
	h := &Handlers{
		List: make(map[string]func(*Peer, *msg.ProtoMessage, *Handlers)),
		Host: n,
	}

	if !n.Discovery {
		// none discovery node specific
		h.List["get_peers"] = h.get_peers
		h.List["peers"] = h.peers
	} else {
		log.Info().Msg("Booting discovery node")
	}

	return h
}

func (h *Handlers) AddMethod(cmd string, fn func(*Peer, *msg.ProtoMessage, *Handlers)) *Handlers {
	h.List[cmd] = fn
	return h
}

func (h *Handlers) peers(p *Peer, m *msg.ProtoMessage, _ *Handlers) {
	// log.Error().Msg("Unimplemented")
	dec, _ := hex.DecodeString(m.Data)
	peers_b := strings.Split(string(dec), ",")

	// h.host.Mu.Lock()
	for _, p := range peers_b {
		if len(p) > 1 {
			z := strings.SplitN(p, ":", 2)

			if h.Host.Peers.HasPeer(p) == true || p == h.Host.Address {
				continue
			}

			c, err := net.Dial("tcp", z[1])
			if err != nil {
				continue
			}
			h.Host.Outbound(c)
		}
	}
	// h.host.Mu.Unlock()
}

func (h *Handlers) get_peers(p *Peer, m *msg.ProtoMessage, _ *Handlers) {
	h.Host.Mu.Lock()

	b := []byte(h.Host.Peers.MarshalPeers())

	p.Send(msg.Msg(h.Host.Address, "peers", b))
	h.Host.Mu.Unlock()
}
