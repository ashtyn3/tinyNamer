package discover

import (
	"encoding/hex"
	"strings"

	"github.com/ashtyn3/tinynamer/msg"
	"github.com/ashtyn3/tinynamer/p2p"
)

func peers(p *p2p.Peer, m *msg.ProtoMessage, h *p2p.Handlers) {
	// log.Error().Msg("Unimplemented")
	dec, _ := hex.DecodeString(m.Data)
	peers_b := strings.Split(string(dec), ",")

	h.Host.Mu.Lock()
	for _, p_b := range peers_b {
		if len(p_b) > 1 {
			if h.Host.Peers.HasPeer(p_b) == true || p_b == h.Host.Address {
				continue
			}
			peer := &p2p.Peer{
				Address: strings.Split(p_b, "@")[0],
				Ip:      strings.Split(p_b, "@")[1],
				Port:    strings.Split(strings.Split(p_b, "@")[1], ":")[1],
				Hash:    "",
				Sock:    nil,
				Halt:    false,
			}

			h.Host.Peers.AddPeer(peer)
		}
	}
	p.Send(msg.Msg(h.Host.Address, "kill", []byte{}))
	p.Sock.Close()
	p.Halt = true
	h.Host.Mu.Unlock()
}

func get_peers(p *p2p.Peer, m *msg.ProtoMessage, h *p2p.Handlers) {
	h.Host.Mu.Lock()

	b := []byte(h.Host.Peers.MarshalPeers())

	p.Send(msg.Msg(h.Host.Address, "peers", b))
	h.Host.Mu.Unlock()
}

func Run() {
	n := p2p.NewNode(true)
	n.Address = "DISCOVERY@" + n.Ip
	n.Handlers.AddMethod("get_peers", get_peers)
	n.Handlers.AddMethod("peers", peers)

	n.Run("5779")
}
