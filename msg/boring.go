package msg

import (
	"crypto/sha256"
	"encoding/hex"
)

type Message struct {
	Header string
	Body   string
}

func Msg(address string, body []byte) *Message {
	m := &Message{}
	// enc := base58.RippleEncoding
	// b, _ := enc.Encode(body)
	m.Body = string(body)
	m.Header = address

	return m
}

func (m *Message) Marshal() []byte {
	sha := sha256.New()
	sha.Write([]byte(m.Header + m.Body))
	hash := hex.EncodeToString(sha.Sum(nil))
	return []byte(m.Header + "\n" + hash + "\n" + m.Body + "\n" + "$")
}
