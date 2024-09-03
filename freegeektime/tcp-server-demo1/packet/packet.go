package packet

import (
	"bytes"
	"fmt"
)

const (
	CommandConn = iota + 0x01
	CommandSubmit
)

const (
	CommandConnAck = iota + 0x80
	CommandSubmitAck
)

type Packet interface {
	Decode([]byte) error
	Encode() ([]byte, error)
}

type Submit struct {
	ID      string
	Payload []byte
}

func (s *Submit) Decode(body []byte) error {
	s.ID = string(body[:8])
	s.Payload = body[8:]
	return nil
}

func (s *Submit) Encode() ([]byte, error) {
	return bytes.Join([][]byte{[]byte(s.ID[:8]), s.Payload}, nil), nil
}

type SubmitAck struct {
	ID     string
	Result uint8
}

func (s *SubmitAck) Decode(body []byte) error {
	s.ID = string(body[0:8])
	s.Result = uint8(body[8])
	return nil
}

func (s *SubmitAck) Encode() ([]byte, error) {
	return bytes.Join([][]byte{[]byte(s.ID[:8]), []byte{s.Result}}, nil), nil
}

func Decode(packet []byte) (Packet, error) {
	commandID := packet[0]

	body := packet[1:]

	switch commandID {
	case CommandConn:
		return nil, nil
	case CommandConnAck:
		return nil, nil
	case CommandSubmit:
		s := &Submit{}
		if err := s.Decode(body); err != nil {
			return nil, err
		}
		return s, nil
	case CommandSubmitAck:
		s := &SubmitAck{}
		if err := s.Decode(body); err != nil {
			return nil, err
		}
		return s, nil
	default:
		return nil, fmt.Errorf("unkonwn commandID [%d]", commandID)
	}
}

func Encode(p Packet) ([]byte, error) {
	var commandID uint8
	var pktBody []byte
	var err error

	switch t := p.(type) {
	case *Submit:
		commandID = CommandSubmit
		pktBody, err = p.Encode()
		if err != nil {
			return nil, err
		}
	case *SubmitAck:
		commandID = CommandSubmitAck
		pktBody, err = p.Encode()
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unknown type [%s]", t)
	}
	return bytes.Join([][]byte{[]byte{commandID}, pktBody}, nil), nil
}
