package p2p

import (
	"bytes"
	"encoding/hex"
	"net"
	"strings"

	"github.com/ashtyn3/tinynamer/msg"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/proto"
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
	// log.Error().Msg("Unimplemented")
	str, _ := hex.DecodeString(m.Data)
	peers_b := bytes.Split(str, []byte("\r\r\n"))

	// h.host.Mu.Lock()
	for _, p := range peers_b {
		if len(p) > 1 {
			temp_p := &ProtoPeer{}
			proto.Unmarshal(p, temp_p)
			z := strings.SplitN(temp_p.Address, ":", 2)

			if h.host.Peers.peers[temp_p.Address] != nil || temp_p.Address == h.host.Address {
				continue
			}

			c, err := net.Dial("tcp", z[1])
			log.Error().Err(err)
			h.host.outbound(c)
		}
	}
	// h.host.Mu.Unlock()
}

func (h *Handlers) get_peers(p *Peer, m *msg.ProtoMessage) {
	h.host.Mu.Lock()

	b := bytes.Join(h.host.Peers.PartialMarshal(), []byte("\r\r\n"))

	p.Send(msg.Msg(h.host.Address, "peers", b))
	h.host.Mu.Unlock()
}
