package frame

import (
	"bytes"
	"encoding/binary"
	"testing"
)

func TestNewMyFrameCodec(t *testing.T) {
	codec := NewMyFrameCodec()
	if codec == nil {
		t.Errorf("want non-nil,actual nil")
	}
}

func TestEncode(t *testing.T) {
	codec := NewMyFrameCodec()
	buf := make([]byte, 0, 128)
	rw := bytes.NewBuffer(buf)

	if err := codec.Encode(rw, []byte("hello")); err != nil {
		t.Errorf("want nil,actual %s", err.Error())
	}

	var totalLen int32
	if err := binary.Read(rw, binary.BigEndian, &totalLen); err != nil {
		t.Errorf("want nil,actual %s", err.Error())
	}

	if totalLen != 9 {
		t.Errorf("want 9 , actual %d", totalLen)
	}

	left := rw.Bytes()

	if string(left) != "hello" {
		t.Errorf("want hello,actual %s", string(left))

	}
}

func TestDecode(t *testing.T) {
	var codec = NewMyFrameCodec()
	var data []byte = []byte{0x0, 0x0, 0x0, 0x9, 'h', 'e', 'l', 'l', 'o'}
	payload, err := codec.Decode(bytes.NewBuffer(data))
	if err != nil {
		t.Errorf("want nil,actual %s", err.Error())
	}
	if string(payload) != "hello" {
		t.Errorf("want hello,actual %s", string(payload))
	}
}
