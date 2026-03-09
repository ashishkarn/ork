package main

import (
	"encoding/binary"
	"fmt"
)

type MsgType uint8

const (
	MsgDiscover MsgType = 0x01
	MsgAnnounce MsgType = 0x02
)

type Message struct {
	Type   MsgType
	NodeID string
	Port   uint16
}

func (m *Message) Encode() []byte {
	idBytes := []byte(m.NodeID)
	buf := make([]byte, 0, 2+len(idBytes)+2)

	buf = append(buf, uint8(m.Type))
	buf = append(buf, uint8(len(idBytes)))
	buf = append(buf, idBytes...)

	if m.Type == MsgAnnounce {
		port := make([]byte, 2)
		binary.BigEndian.PutUint16(port, m.Port)
		buf = append(buf, port...)
	}

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

	if msg.Type == MsgAnnounce {
		if len(data) < 2+idLen+2 {
			return nil, fmt.Errorf("truncated tcp port")
		}
		msg.Port = binary.BigEndian.Uint16(data[2+idLen : 2+idLen+2])
	}

	return msg, nil
}
