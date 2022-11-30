package p2p

import (
	"os"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
	"google.golang.org/protobuf/proto"
)

type Store struct {
	db   *leveldb.DB
	path string
	self string
}

const PERM = 0777

func NewStore(basePath string, self string) *Store {
	os.MkdirAll(basePath+"/store", PERM)

	store := &Store{
		path: basePath + "/store",
		self: self,
	}
	d, err := leveldb.OpenFile(store.path, nil)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
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

func (ps *Store) AddPeer(p *Peer) {
	if p.Address == ps.self || ps.HasPeer(p.Address) {
		return
	}
	log.Info().Str("Address", p.Address).Msg("connected new peer")
	b, _ := proto.Marshal(p.ToProtoPeer())
	err := ps.db.Put([]byte("tn://"+p.Address), b, nil)
	if err != nil {
		log.Error().Err(err).Send()
	}
}
func (ps *Store) MarshalPeers() string {
	iter := ps.db.NewIterator(util.BytesPrefix([]byte("tn://")), nil)
	buffer := []string{}
	for iter.Next() {
		value := iter.Value()
		pb := &ProtoPeer{}
		proto.Unmarshal(value, pb)
		buffer = append(buffer, pb.Address+"@"+pb.Ip+":"+pb.Port)
	}
	iter.Release()
	err := iter.Error()
	log.Error().Err(err)

	return strings.Join(buffer, ",")
}

func (ps *Store) HasPeer(addr string) bool {
	state, _ := ps.db.Has([]byte("tn://"+addr), nil)

	return state
}
