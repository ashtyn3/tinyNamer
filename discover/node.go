package discover

import (
	"encoding/hex"
	"net"
	"strings"

	"github.com/ashtyn3/tinynamer/msg"
	"github.com/ashtyn3/tinynamer/p2p"
	"github.com/rs/zerolog/log"
)

func peers(p *p2p.Peer, m *msg.ProtoMessage, h *p2p.Handlers) {
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
<<<<<<< HEAD
			peer := &p2p.Peer{
				Address: strings.Split(p_b, "@")[0],
				Ip:      strings.Split(p_b, "@")[1],
				Port:    strings.Split(strings.Split(p_b, "@")[1], ":")[1],
				Hash:    "",
				Sock:    nil,
				Halt:    false,
			}
=======
>>>>>>> parent of dd78d42 (fix: peer discovery node closes connection)

			c, err := net.Dial("tcp", z[1])
			log.Error().Err(err)
			h.Host.Outbound(c)
		}
	}
	p.Halt = true
	p.Sock.Close()
	// h.host.Mu.Unlock()
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
