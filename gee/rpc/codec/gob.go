package codec

import (
	"bufio"
	"encoding/gob"
	"io"
	"log"
)

type GobCodec struct {
	conn io.ReadWriteCloser
	buf  *bufio.Writer
	dec  *gob.Decoder
	enc  *gob.Encoder
}

func (c *GobCodec) ReadHeader(h *Header) error {
	log.Printf("ReadHeader===>:%v\n", h)
	return c.dec.Decode(h)
}

func (c *GobCodec) ReadBody(body interface{}) error {
	log.Printf("ReadBody===>:%v\n", body)
	return c.dec.Decode(body)
}
func (c *GobCodec) Write(h *Header, body interface{}) (err error) {
	defer func() {
		c.buf.Flush()
		if err != nil {
			c.Close()
		}
	}()
	log.Printf("Write====>header:%v body:%v\n", h, body)

	if err = c.enc.Encode(h); err != nil {
		log.Println("rpc codec :gob error encoding header :", err)
		return
	}
	if err = c.enc.Encode(body); err != nil {
		log.Println("rpc codec :gob error encoding body :", err)
		return
	}
	return nil
}
func (c *GobCodec) Close() error {
	return nil
}

var _ Codec = (*GobCodec)(nil)

func NewGobCodec(conn io.ReadWriteCloser) Codec {
	buf := bufio.NewWriter(conn)
	return &GobCodec{
		conn: conn,
		buf:  buf,
		dec:  gob.NewDecoder(conn),
		enc:  gob.NewEncoder(buf),
	}
}
