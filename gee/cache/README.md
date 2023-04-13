grpc.go

package geecache

import (
	"context"
	"fmt"
	"geecache/geecache/consistent"
	pb "geecache/geecache/geecachepb"
	"geecache/geecache/peers"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"sync"
)

type grpcGetter struct {
	addr string
}

func (g *grpcGetter) Get(in *pb.Request, out *pb.Response) error {
	c, err := grpc.Dial(g.addr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	client := pb.NewGroupCacheClient(c)
	response, err := client.Get(context.Background(), in)
	out.Value = response.Value
	return err
}

var _ peers.PeerGetter = (*grpcGetter)(nil)

type GrpcPool struct {
	pb.UnimplementedGroupCacheServer

	self        string
	mu          sync.Mutex
	peers       *consistent.Map
	grpcGetters map[string]*grpcGetter
}

func NewGrpcPool(self string) *GrpcPool {
	return &GrpcPool{
		self:        self,
		peers:       consistent.New(defaultReplicas, nil),
		grpcGetters: map[string]*grpcGetter{},
	}
}

func (p *GrpcPool) Set(peers ...string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.peers.Add(peers...)
	for _, peer := range peers {
		p.grpcGetters[peer] = &grpcGetter{
			addr: peer,
		}
	}
}

func (p *GrpcPool) PickPeer(key string) (peers.PeerGetter, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if peer := p.peers.Get(key); peer != "" && peer != p.self {
		return p.grpcGetters[peer], true
	}
	return nil, false
}

var _ peers.PeerPicker = (*GrpcPool)(nil)

func (p *GrpcPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s", p.self, fmt.Sprintf(format, v...))
}

func (p *GrpcPool) Get(ctx context.Context, in *pb.Request) (*pb.Response, error) {
	p.Log("%s %s", in.Group, in.Key)
	response := &pb.Response{}

	group := GetGroup(in.Group)
	if group == nil {
		p.Log("no such group %v", in.Group)
		return response, fmt.Errorf("no such group %v", in.Group)
	}
	value, err := group.Get(in.Key)
	if err != nil {
		p.Log("get key %v error %v", in.Key, err)
		return response, err
	}

	response.Value = value.ByteSlice()
	return response, nil
}

func (p *GrpcPool) Run() {
	lis, err := net.Listen("tcp", p.self)
	if err != nil {
		panic(err)
	}

	server := grpc.NewServer()
	pb.RegisterGroupCacheServer(server, p)

	reflection.Register(server)
	err = server.Serve(lis)
	if err != nil {
		panic(err)
	}
}
main.go

func startCacheServerGrpc(addr string, addrs []string, gee *geecache.Group) {
	peers := geecache.NewGrpcPool(addr)
	peers.Set(addrs...)
	gee.RegisterPeers(peers)
	log.Println("geecache is running at", addr)
	peers.Run()
}

func startGRPCServer() {
	var port int
	var api bool
	flag.IntVar(&port, "port", 8001, "Geecache server port")
	flag.BoolVar(&api, "api", false, "Start a api server?")
	flag.Parse()

	apiAddr := "http://localhost:9999"
	addrMap := map[int]string{
		8001: ":8001",
		8002: ":8002",
		8003: ":8003",
	}

	var addrs []string
	for _, v := range addrMap {
		addrs = append(addrs, v)
	}

	gee := createGroup()
	if api {
		go startAPIServer(apiAddr, gee)
	}
	startCacheServerGrpc(addrMap[port], addrs, gee)
}