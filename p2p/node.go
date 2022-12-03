package p2p

import (
	"bufio"
	"bytes"
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
	"google.golang.org/protobuf/proto"
)

type Node struct {
	Keypair    *ecdsa.PrivateKey
	Address    string
	PublicKey  string
	Ip         string
	host_ip    string
	Peers      *Store
	Mu         sync.Mutex
	Handlers   *Handlers
	Listen_net net.Listener
	BasePath   string
	Discovery  bool
}

func NewNode(disc bool) *Node {
	n := &Node{}
	prv, _ := crypto.GenerateKey()
	n.Keypair = prv
	n.Mu = sync.Mutex{}
	pub := crypto.FromECDSAPub(&n.Keypair.PublicKey)
	n.Address = hex.EncodeToString(pub)
	n.PublicKey = n.Address
	home, _ := os.UserHomeDir()
	n.BasePath = home + "/.tinyNamer"
	n.Peers = NewStore(n.BasePath, n.Address)
	n.Discovery = disc
	n.Handlers = InitHandlers(n)

	return n
}
func (n *Node) listener() {

	for {
		c, err := n.Listen_net.Accept()
		if err != nil {
			log.Fatal().Err(err).Send()
		}
		p := NewPeer(c)

		go n.handle(p)
		p.Send(msg.Msg(n.Address, "get_peers", []byte(n.Ip)))
	}
}

func (n *Node) handle(peer *Peer) {
	reader := bufio.NewReader(peer.Sock)
	lh := ""
	for peer.Halt == false {
		data, err := reader.ReadBytes('$')
		data = bytes.ReplaceAll(data, []byte("$"), []byte(""))
		m := &msg.ProtoMessage{}
		MarshalErr := proto.Unmarshal(data, m)
		log.Error().Err(MarshalErr)
		if m.Hash != "" && m.Hash != lh {
			lh = m.Hash
		} else if m.Hash == "" {
		} else {
			continue
		}

		switch err {
		case nil:
			{
				if !peer.developed {
					temp := strings.Split(m.Address, "@")
					peer.Address = temp[0]
					peer.Port = strings.Split(temp[1], ":")[1]

<<<<<<< HEAD
					if peer.Address != "DISCOVERY" {
						n.Mu.Lock()
						n.Peers.AddPeer(peer)
						n.Mu.Unlock()
					}
=======
					n.Mu.Lock()
					n.Peers.AddPeer(peer)
					n.Mu.Unlock()
>>>>>>> parent of dd78d42 (fix: peer discovery node closes connection)
					n.Mu.Lock()
					peer.developed = true
					n.Mu.Unlock()
					n.Handlers.List[m.Command](peer, m, n.Handlers)
				} else {
					if n.Handlers.List[m.Command] == nil {
						peer.Send(msg.Msg(n.Address, "unknown", nil))
					}
					n.Handlers.List[m.Command](peer, m, n.Handlers)
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

func (n *Node) Outbound(con net.Conn) {
	p := NewPeer(con)

	go n.handle(p)
	p.Send(msg.Msg(n.Address, "get_peers", nil))
	// n.peers.Add(p)
}

func (n *Node) Discover() {
<<<<<<< HEAD
	// for i := 0; i < 7; i++ {
	// 	port := strconv.Itoa(5770 + i)
	// 	if strings.Contains(n.Ip, port) {
	// 		continue
	// 	}
	// 	con, err := net.Dial("tcp", ":"+port)
	// 	if err != nil {
	// 		continue
	// 	}
	// 	n.Outbound(con)
	// 	break
	// }

	// discovery node finder
	disc_nodes := []string{"0.0.0.0:5779", "8.9.36.237:5779"}
	for _, p := range disc_nodes {
		con, err := net.Dial("tcp", p)
=======
	for i := 0; i < 7; i++ {
		port := strconv.Itoa(5770 + i)
		if strings.Contains(n.Ip, port) {
			continue
		}
		con, err := net.Dial("tcp", ":"+port)
>>>>>>> parent of dd78d42 (fix: peer discovery node closes connection)
		if err != nil {
			continue
		}
		n.Outbound(con)
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

func (n *Node) Run(port string) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	go func() {
		select {
		case <-c:
			n.Peers.db.Close()
			os.Exit(0)
		}
	}()

	rand.Seed(time.Now().UnixNano())
	// port := 5770 + rand.Intn(7)
	// port_s := strconv.Itoa(port)
	log.Info().Msg(fmt.Sprint("listening on ", port))
	l, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	n.Listen_net = l
	defer n.Listen_net.Close()
	n.Ip = l.Addr().String()
<<<<<<< HEAD
	if n.Discovery {
		n.Address = "DISCOVERY@127.0.0.1:" + port
	} else {
		n.Address += "@127.0.0.1:" + port
	}
	n.Peers = NewStore(n.BasePath, n.Address)
=======
	n.Address += ":" + n.Ip
>>>>>>> parent of dd78d42 (fix: peer discovery node closes connection)

	n.Discover()
	n.listener()
}
