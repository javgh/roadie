package rpc

import (
	"context"
	"encoding/json"
	"errors"
	"net"

	"gitlab.com/NebulousLabs/Sia/types"
	"google.golang.org/grpc"
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
		HandlerType: (*RoadieServer)(nil),
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
	RoadieServer interface {
		RequestNonBindingOffer(*RNBORequest) (*RNBOResponse, error)
	}

	BobRoadieServer struct {
		Playground bob.AtomicSwap
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

func (s *BobRoadieServer) RequestNonBindingOffer(req *RNBORequest) (*RNBOResponse, error) {
	var err error
	resp := new(RNBOResponse)

	resp.Offer, err = s.Playground.RequestNonBindingOffer(req.Siacoin)
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

	return srv.(RoadieServer).RequestNonBindingOffer(in)
}

func Playground(s bob.AtomicSwap) error {
	roadieServer := BobRoadieServer{s}
	grpcServer := grpc.NewServer() // TODO: grpc.Creds(creds)
	grpcServer.RegisterService(&serviceDesc, &roadieServer)

	lis, err := net.Listen("tcp", "localhost:9000")
	if err != nil {
		return err
	}

	err = grpcServer.Serve(lis)
	if err != nil {
		return err
	}

	return nil
}

type RoadieClient struct {
	conn *grpc.ClientConn
}

func (c *RoadieClient) RequestNonBindingOffer(siacoin types.Currency) (*trader.Offer, error) {
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

func Dial(target string) (*RoadieClient, error) {
	conn, err := grpc.Dial(target,
		grpc.WithInsecure(),
		grpc.WithDefaultCallOptions(grpc.CallContentSubtype(JSONCodec{}.Name())),
	)
	if err != nil {
		return nil, err
	}

	client := RoadieClient{conn}
	return &client, nil
}
