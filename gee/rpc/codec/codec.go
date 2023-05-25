package codec

import "io"

type Header struct {
	//服务名和方法名
	ServiceMethod string
	//请求的序号
	Seq uint64
	//错误信息
	Error string
}

//抽象出对消息体进行编解码的接口
type Codec interface {
	io.Closer
	//读取头
	ReadHeader(*Header) error
	//获取body数据
	ReadBody(interface{}) error
	//写入数据
	Write(*Header, interface{}) error
}

type NewCodeFunc func(io.ReadWriteCloser) Codec

type Type string

const (
	GobType  Type = "application/gob"
	JsonType Type = "application/json"
)

var NewCodecFuncMap map[Type]NewCodeFunc

func init() {
	NewCodecFuncMap = make(map[Type]NewCodeFunc)
	NewCodecFuncMap[GobType] = NewGobCodec
}
