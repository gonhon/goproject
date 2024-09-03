package frame

import (
	"encoding/binary"
	"errors"
	"io"
)

var (
	ErrShortWrite = errors.New("short write")
	ErrShortRead  = errors.New("short read")
)

type FramePaload []byte

type StreamFrameCodec interface {
	Encode(io.Writer, FramePaload) error
	Decode(io.Reader) (FramePaload, error)
}

type myFrameCodec struct{}

func NewMyFrameCodec() StreamFrameCodec {
	return &myFrameCodec{}
}

func (*myFrameCodec) Encode(write io.Writer, framePaload FramePaload) error {
	var f = framePaload
	var totalLen int32 = int32(len(f)) + 4

	err := binary.Write(write, binary.BigEndian, &totalLen)
	if err != nil {
		return nil
	}

	n, err := write.Write([]byte(f))
	if err != nil {
		return nil
	}

	if n != len(f) {
		return ErrShortWrite
	}

	return nil
}

func (*myFrameCodec) Decode(read io.Reader) (FramePaload, error) {
	var totalLen int32
	if err := binary.Read(read, binary.BigEndian, &totalLen); err != nil {
		return nil, err
	}

	buf := make([]byte, totalLen-4)
	n, err := io.ReadFull(read, buf)
	if err != nil {
		return nil, err
	}

	if n != int(totalLen-4) {
		return nil, ErrShortRead
	}
	return FramePaload(buf), nil
}
