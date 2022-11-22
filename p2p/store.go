package p2p

import (
	"bytes"
	"encoding/hex"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/syndtr/goleveldb/leveldb"
	"google.golang.org/protobuf/proto"
)

type PeerStore struct {
	db   *leveldb.DB
	path string
	self string
}

const PERM = 0777

func NewStore(basePath string, self string) *PeerStore {
	os.MkdirAll(basePath+"/peers.store", PERM)

	store := &PeerStore{
		path: basePath + "/peers.store",
		self: self,
	}
	d, err := leveldb.OpenFile(store.path, nil)
	log.Error().Err(err)
	store.db = d

	return store
}

func (p *Peer) ToProtoPeer() *ProtoPeer {
	return &ProtoPeer{
		Address:   p.Address,
		Ip:        p.Ip,
		Port:      p.Port,
		Hash:      p.Hash,
		Developed: false,
	}
}

func (p *ProtoPeer) ToPeer() *Peer {
	return &Peer{
		Address:   p.Address,
		Ip:        p.Ip,
		Port:      p.Port,
		Hash:      p.Hash,
		Sock:      nil,
		developed: p.Developed,
		Halt:      false,
	}
}

func (ps *PeerStore) Add(p *Peer) {
	b, _ := proto.Marshal(p.ToProtoPeer())
	ps.db.Put([]byte(p.Ip), b, nil)
	// if err != nil {
	// 	log.Error().Err(err)
	// }
}
func (ps *PeerStore) Marshal() string {
	iter := ps.db.NewIterator(nil, nil)
	buffer := [][]byte{}
	for iter.Next() {
		// Remember that the contents of the returned slice should not be modified, and
		// only valid until the next call to Next.
		value := iter.Value()
		buffer = append(buffer, value)
	}
	iter.Release()
	err := iter.Error()
	log.Error().Err(err)

	return hex.EncodeToString(bytes.Join(buffer, []byte("\r\r\n")))
}

func (ps *PeerStore) HasPeer(addr string) bool {
	state, _ := ps.db.Has([]byte(addr), nil)

	return state
}
