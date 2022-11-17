package p2p

import (
	"bytes"
	"os"

	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/proto"
)

type PeerStore struct {
	peers  map[string]*Peer
	Buffer [][]byte
}

func NewStore() *PeerStore {
	ps := &PeerStore{}
	ps.peers = make(map[string]*Peer)
	return ps
}

func (ps *PeerStore) Add(p *Peer) {
	// pb := &ProtoPeer{
	// 	Address:   p.Address,
	// 	Ip:        p.Ip,
	// 	Port:      p.Port,
	// 	Hash:      p.Hash,
	// 	Developed: p.developed,
	// }
	// b, _ := proto.Marshal(pb)
	// ps.buffer = append(ps.buffer, b)
	if ps.peers[p.Address] == nil {
		ps.peers[p.Address] = p
		log.Info().Str("Address", p.Address).Msgf("connected new peer")
	}
}

func (ps *PeerStore) PartialMarshal() {
	ps.Buffer = [][]byte{}

	for _, p := range ps.peers {
		pb := &ProtoPeer{
			Address:   p.Address,
			Ip:        p.Ip,
			Port:      p.Port,
			Hash:      p.Hash,
			Developed: p.developed,
		}
		b, _ := proto.Marshal(pb)
		ps.Buffer = append(ps.Buffer, b)
	}
}
func (ps *PeerStore) Marshal() {
	home, _ := os.UserHomeDir()
	os.Mkdir(home+"/.tinyNamer", os.ModePerm)

	ps.PartialMarshal()

	os.WriteFile(home+"/.tinyNamer/peers.store", bytes.Join(ps.Buffer, []byte("\n")), 0600)
}

func (ps *PeerStore) Unmarshal() {
	home, _ := os.UserHomeDir()
	b, err := os.ReadFile(home + "/.tinyNamer/peers.store")
	if os.IsNotExist(err) {
		return
	}

	peers_bytes := bytes.Split(b, []byte("\n"))

	for _, bytes := range peers_bytes {
		ppeer := &ProtoPeer{}
		proto.Unmarshal(bytes, ppeer)
		peer := &Peer{
			Address:   ppeer.GetAddress(),
			Ip:        ppeer.GetIp(),
			Port:      ppeer.GetPort(),
			Hash:      ppeer.GetHash(),
			Sock:      nil,
			developed: ppeer.GetDeveloped(),
			Halt:      false,
		}
		ps.peers[peer.Address] = peer
	}
}
