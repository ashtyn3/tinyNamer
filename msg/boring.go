package msg

import (
	"crypto/sha256"
	"encoding/hex"

	"google.golang.org/protobuf/proto"
)

type Message struct {
	internal *ProtoMessage
}

func Msg(address string, command string, body []byte) *Message {
	m := &Message{
		internal: &ProtoMessage{
			Address: address,
			Hash:    "",
			Command: command,
			Data:    hex.EncodeToString(body),
		},
	}

	// enc := base58.RippleEncoding
	// b, _ := enc.Encode(body)
	return m
}

func (m *Message) Marshal() []byte {
	sha := sha256.New()
	sha.Write([]byte(m.internal.Address + m.internal.Command))
	hash := hex.EncodeToString(sha.Sum(nil))
	m.internal.Hash = hash

	b, _ := proto.Marshal(m.internal)

	b = append(b, '$')

	return b
}
