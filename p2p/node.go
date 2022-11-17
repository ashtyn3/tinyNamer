package p2p

import (
	"bufio"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"io"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ashtyn3/tinynamer/msg"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog/log"
)

type Node struct {
	keypair *ecdsa.PrivateKey
	address string
	peers   *PeerStore
	mu      sync.Mutex
}

func NewNode() *Node {
	n := &Node{}
	prv, _ := crypto.GenerateKey()
	n.keypair = prv
	n.mu = sync.Mutex{}
	pub := crypto.FromECDSAPub(&n.keypair.PublicKey)
	n.address = hex.EncodeToString(pub)
	n.peers = NewStore()

	return n
}
func (n *Node) listener() {
	rand.Seed(time.Now().UnixNano())
	port := 5770 + rand.Intn(7)
	port_s := strconv.Itoa(port)
	log.Info().Msg(fmt.Sprint("listening on ", port))
	l, err := net.Listen("tcp", ":"+port_s)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			log.Fatal().Err(err).Send()
		}
		p := NewPeer(c)

		go n.handle(p)
		p.Send(msg.Msg(n.address, []byte("get_peers")))
	}
}

func (n *Node) handle(peer *Peer) {
	reader := bufio.NewReader(peer.Sock)
	for peer.Halt == false {
		data, err := reader.ReadString('$')
		switch err {
		case nil:
			{
				frags := strings.Split(data, "\n")
				if !peer.developed {
					pub_key := frags[0]
					peer.Address = pub_key

					n.mu.Lock()
					n.peers.Add(peer)
					n.mu.Unlock()
				} else {
					fmt.Printf("(%s:%s): %s", peer.Ip, peer.Port, data)
				}
			}
		case io.EOF:
			{
				log.Info().Str("Address", peer.Address).Msgf("closed connection")
				peer.Halt = true
			}
		}
	}
	peer.Sock.Close()
}

func (n *Node) outbound(con net.Conn) {
	p := NewPeer(con)

	go n.handle(p)
	p.Send(msg.Msg(n.address, []byte("get_peers")))
	// n.peers.Add(p)
}

func (n *Node) Discover() {
	for i := 0; i < 7; i++ {
		port := strconv.Itoa(5770 + i)
		con, err := net.Dial("tcp", ":"+port)
		if err != nil {
			continue
		}
		n.outbound(con)
		break
	}
	// for _, p := range n.peers.peers {
	// 	con, err := net.Dial("tcp", p.Ip+":"+p.Port)
	// 	if err != nil {
	// 		continue
	// 	}
	// 	n.outbound(con)
	// }
}

func (n *Node) Run() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	go func() {
		select {
		case <-c:
			n.peers.Marshal()
			os.Exit(0)
		}
	}()

	n.peers.Unmarshal()
	n.Discover()
	n.listener()
}
