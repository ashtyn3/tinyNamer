package p2p

import (
	"bufio"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ashtyn3/tinynamer/msg"
	"github.com/ethereum/go-ethereum/crypto"
)

type Node struct {
	address *ecdsa.PrivateKey
	peers   PeerStore
	mu      sync.Mutex
}

func NewNode() *Node {
	n := &Node{}
	prv, _ := crypto.GenerateKey()
	n.address = prv
	n.mu = sync.Mutex{}

	return n
}
func (n *Node) listener() {
	rand.Seed(time.Now().UnixNano())
	port := 5770 + rand.Intn(7)
	port_s := strconv.Itoa(port)
	log.Println("listening on", port)
	l, err := net.Listen("tcp", ":"+port_s)
	if err != nil {
		log.Fatalln(err)
	}
	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			log.Fatalln(err)
		}
		p := NewPeer(c)
		go n.handle(p)
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
				} else {
					fmt.Printf("(%s:%s): %s", peer.Ip, peer.Port, data)
				}
			}
		}
	}
	peer.Sock.Close()
}

func (n *Node) Discover() {
	for i := 0; i < 7; i++ {
		port := strconv.Itoa(5770 + i)
		con, err := net.Dial("tcp", ":"+port)
		if err != nil {
			fmt.Println(err)
			continue
		}
		p := NewPeer(con)
		go n.handle(p)
		pub := crypto.FromECDSAPub(&n.address.PublicKey)
		p.Send(msg.Msg(hex.EncodeToString(pub), []byte("get_peers")))
		break
	}
}

func (n *Node) Run() {
	n.Discover()
	n.listener()
}
