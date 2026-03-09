package main

import "fmt"

type MsgType uint8

const (
	MsgDiscover MsgType = 0x01
	MsgAnnounce MsgType = 0x02
)

type Message struct {
	Type   MsgType
	NodeID string
}

func (m *Message) Encode() []byte {
	idBytes := []byte(m.NodeID)
	buf := make([]byte, 0, 2+len(idBytes))

	buf = append(buf, uint8(m.Type))
	buf = append(buf, uint8(len(idBytes)))
	buf = append(buf, idBytes...)

	return buf
}

func Decode(data []byte) (*Message, error) {
	if len(data) < 2 {
		return nil, fmt.Errorf("message too short")
	}

	msg := &Message{}
	msg.Type = MsgType(data[0])
	idLen := int(data[1])

	if len(data) < 2+idLen {
		return nil, fmt.Errorf("truncated node id")
	}

	msg.NodeID = string(data[2 : 2+idLen])

	return msg, nil
}
