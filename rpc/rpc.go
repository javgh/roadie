package rpc

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"sync"
	"time"

	"github.com/satori/go.uuid"
	"gitlab.com/NebulousLabs/Sia/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/encoding"

	"github.com/javgh/roadie/bob"
	"github.com/javgh/roadie/trader"
)

type (
	JSONCodec struct{}
)

var (
	ErrNotImplemented = errors.New("interceptor support is not implemented")

	serviceDesc = grpc.ServiceDesc{
		ServiceName: "Roadie",
		HandlerType: (*Server)(nil),
		Methods: []grpc.MethodDesc{
			{
				MethodName: "RequestNonBindingOffer",
				Handler:    requestNonBindingOfferHandler,
			},
		},
		Streams: []grpc.StreamDesc{},
	}
)

func (c JSONCodec) Name() string {
	return "json"
}

func (c JSONCodec) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (c JSONCodec) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func init() {
	encoding.RegisterCodec(JSONCodec{})
}

type (
	Server interface {
		RequestNonBindingOffer(*RNBORequest) (*RNBOResponse, error)
	}

	BobServer struct {
		mutex         sync.Mutex
		atomicSwaps   map[uuid.UUID]*bob.AtomicSwap
		listener      net.Listener
		grpcServer    *grpc.Server
		newAtomicSwap func(now time.Time) *bob.AtomicSwap
	}
)

type (
	RNBORequest struct {
		Siacoin types.Currency
	}

	RNBOResponse struct {
		Offer *trader.Offer
	}
)

func (s *BobServer) RequestNonBindingOffer(req *RNBORequest) (*RNBOResponse, error) {
	var err error
	resp := new(RNBOResponse)

	s.mutex.Lock()
	defer s.mutex.Unlock()

	atomicSwap := s.newAtomicSwap(time.Now())
	s.atomicSwaps[atomicSwap.ID] = atomicSwap

	resp.Offer, err = atomicSwap.RequestNonBindingOffer(req.Siacoin)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func requestNonBindingOfferHandler(srv interface{}, ctx context.Context, dec func(interface{}) error,
	interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	if interceptor != nil {
		return nil, ErrNotImplemented
	}

	in := new(RNBORequest)
	err := dec(in)
	if err != nil {
		return nil, err
	}

	return srv.(Server).RequestNonBindingOffer(in)
}

func NewBobServer(network string, address string, certFile string, keyFile string,
	newAtomicSwap func(now time.Time) *bob.AtomicSwap) (*BobServer, error) {
	opts := []grpc.ServerOption{}
	if certFile != "" && keyFile != "" {
		creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
		if err != nil {
			return nil, err
		}

		opts = append(opts, grpc.Creds(creds))
	}

	listener, err := net.Listen(network, address)
	if err != nil {
		return nil, err
	}

	bobServer := BobServer{
		atomicSwaps:   make(map[uuid.UUID]*bob.AtomicSwap),
		listener:      listener,
		newAtomicSwap: newAtomicSwap,
	}
	bobServer.grpcServer = grpc.NewServer(opts...)
	bobServer.grpcServer.RegisterService(&serviceDesc, &bobServer)

	return &bobServer, nil
}

func (s *BobServer) Serve() error {
	return s.grpcServer.Serve(s.listener)
}

type Client struct {
	conn *grpc.ClientConn
}

func (c *Client) RequestNonBindingOffer(siacoin types.Currency) (*trader.Offer, error) {
	in := RNBORequest{
		Siacoin: siacoin,
	}
	out := new(RNBOResponse)
	err := grpc.Invoke(context.Background(), "/Roadie/RequestNonBindingOffer", &in, out, c.conn)
	if err != nil {
		return nil, err
	}

	return out.Offer, nil
}

func Dial(target string) (*Client, error) {
	conn, err := grpc.Dial(target,
		grpc.WithInsecure(),
		grpc.WithDefaultCallOptions(grpc.CallContentSubtype(JSONCodec{}.Name())),
	)
	if err != nil {
		return nil, err
	}

	client := Client{conn}
	return &client, nil
}
