/*
 * @Author: gaoh
 * @Date: 2023-05-22 23:25:09
 * @LastEditTime: 2023-05-22 23:51:32
 */
package rpc

import (
	"errors"
	"io"
	"log"
	"sync"

	"github.com/limerence-code/goproject/gee/rpc/codec"
)

type Call struct {
	Seq uint64
	//service.method
	ServiceMethod string
	Args          interface{}
	Reply         interface{}
	Error         error
	Done          chan *Call
}

// 支持异步调用 通知回调
func (call *Call) done() {
	call.Done <- call
}

type Client struct {
	//编解码器，和服务端类似，用来序列化将要发送出去的请求，以及反序列化接收到的响应
	cc  codec.Codec
	opt *Option
	//保证请求的有序发送，即防止出现多个请求报文混淆
	sending sync.Mutex
	//每个请求的消息头，header 只有在请求发送时才需要，而请求发送是互斥的，因此每个客户端只需要一个，声明在 Client 结构体中可以复用
	header codec.Header
	mutex  sync.Mutex
	//用于给发送的请求编号，每个请求拥有唯一编号
	seq uint64
	//储未处理完的请求，键是编号，值是 Call 实例
	pending map[uint64]*Call
	//closing 和 shutdown 任意一个值置为 true，则表示 Client 处于不可用的状态，但有些许的差别，closing 是用户主动关闭的，
	//即调用 Close 方法，而 shutdown 置为 true 一般是有错误发生
	closing  bool
	shutdown bool
}

var ErrShutdown = errors.New("connection is shutting down")

var _ io.Closer = (*Client)(nil)

func (client *Client) Close() error {
	client.mutex.Lock()
	defer client.mutex.Unlock()
	if client.closing {
		return ErrShutdown
	}
	client.closing = true
	return client.cc.Close()
}

func (client *Client) IsAvailable() bool {
	client.mutex.Lock()
	defer client.mutex.Unlock()
	return !client.closing && !client.shutdown
}

// 将参数 call 添加到 client.pending 中，并更新 client.seq
func (client *Client) registerCall(call *Call) (uint64, error) {
	client.mutex.Lock()
	defer client.mutex.Unlock()
	if client.closing || client.shutdown {
		return 0, ErrShutdown
	}
	call.Seq = client.seq
	client.pending[call.Seq] = call
	client.seq++
	return client.seq, nil
}

// 根据 seq，从 client.pending 中移除对应的 call，并返回
func (client *Client) removeCall(seq uint64) *Call {
	client.mutex.Lock()
	defer client.mutex.Unlock()

	if call, ok := client.pending[seq]; ok {
		delete(client.pending, seq)
		return call
	}
	log.Printf("seq %d  not found call ...", seq)
	return nil
}

// 服务端或客户端发生错误时调用，将 shutdown 设置为 true，且将错误信息通知所有 pending 状态的 call
func (client *Client) terminalCalls(err error) {
	client.sending.Lock()
	defer client.sending.Unlock()

	client.mutex.Lock()
	defer client.mutex.Unlock()

	client.shutdown = true
	for _, call := range client.pending {
		call.Error = err
		call.done()
	}
}
