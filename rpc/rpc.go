package rpc

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"math/big"
	"net"
	"sort"
	"sync"
	"time"

	"github.com/HyperspaceApp/ed25519"
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
	ErrUnknownID      = errors.New("unknown id")

	serviceDesc = grpc.ServiceDesc{
		ServiceName: "Roadie",
		HandlerType: (*Server)(nil),
		Methods: []grpc.MethodDesc{
			{
				MethodName: "RequestNonBindingOffer",
				Handler:    requestNonBindingOfferHandler,
			},
			{
				MethodName: "RequestBindingOffer",
				Handler:    requestBindingOfferHandler,
			},
			{
				MethodName: "AcceptOffer",
				Handler:    acceptOfferHandler,
			},
			{
				MethodName: "EnableFunding",
				Handler:    enableFundingHandler,
			},
			{
				MethodName: "RequestAdaptorDetails",
				Handler:    requestAdaptorDetailsHandler,
			},
			{
				MethodName: "AnnounceDeposit",
				Handler:    announceDepositHandler,
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
		RequestBindingOffer(req *RBORequest) (*RBOResponse, error)
		AcceptOffer(req *AORequest) (*AOResponse, error)
		EnableFunding(req *EFRequest) (*EFResponse, error)
		RequestAdaptorDetails(req *RADRequest) (*RADResponse, error)
		AnnounceDeposit(req *ADRequest) (*ADResponse, error)
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
		ID    uuid.UUID
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

	log.Printf("[%s] RequestNonBindingOffer; %s\n", atomicSwap.ID, req.Siacoin.HumanString())

	resp.Offer, err = atomicSwap.RequestNonBindingOffer(req.Siacoin, time.Now())
	if err != nil {
		return nil, err
	}
	resp.ID = atomicSwap.ID

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

type (
	RBORequest struct {
		ID         uuid.UUID
		AntiSpamID big.Int
	}

	RBOResponse struct {
		Offer *trader.Offer
	}
)

func (s *BobServer) RequestBindingOffer(req *RBORequest) (*RBOResponse, error) {
	var err error
	resp := new(RBOResponse)

	s.mutex.Lock()
	defer s.mutex.Unlock()

	atomicSwap, ok := s.atomicSwaps[req.ID]
	if !ok {
		return nil, ErrUnknownID
	}

	log.Printf("[%s] RequestBindingOffer; %s\n", atomicSwap.ID, req.AntiSpamID.String())

	resp.Offer, err = atomicSwap.RequestBindingOffer(req.AntiSpamID, time.Now())
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func requestBindingOfferHandler(srv interface{}, ctx context.Context, dec func(interface{}) error,
	interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	if interceptor != nil {
		return nil, ErrNotImplemented
	}

	in := new(RBORequest)
	err := dec(in)
	if err != nil {
		return nil, err
	}

	return srv.(Server).RequestBindingOffer(in)
}

type (
	AORequest struct {
		ID          uuid.UUID
		AlicePubKey ed25519.PublicKey
	}

	AOResponse struct {
		RefundDetails *bob.RefundDetails
	}
)

func (s *BobServer) AcceptOffer(req *AORequest) (*AOResponse, error) {
	var err error
	resp := new(AOResponse)

	s.mutex.Lock()
	defer s.mutex.Unlock()

	atomicSwap, ok := s.atomicSwaps[req.ID]
	if !ok {
		return nil, ErrUnknownID
	}

	log.Printf("[%s] AcceptOffer\n", atomicSwap.ID)

	resp.RefundDetails, err = atomicSwap.AcceptOffer(req.AlicePubKey, time.Now())
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func acceptOfferHandler(srv interface{}, ctx context.Context, dec func(interface{}) error,
	interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	if interceptor != nil {
		return nil, ErrNotImplemented
	}

	in := new(AORequest)
	err := dec(in)
	if err != nil {
		return nil, err
	}

	return srv.(Server).AcceptOffer(in)
}

type (
	EFRequest struct {
		ID                    uuid.UUID
		AliceRefundNoncePoint ed25519.CurvePoint
		RefundSigAlice        []byte
	}

	EFResponse struct {
		TxID *types.TransactionID
	}
)

func (s *BobServer) EnableFunding(req *EFRequest) (*EFResponse, error) {
	var err error
	resp := new(EFResponse)

	s.mutex.Lock()
	defer s.mutex.Unlock()

	atomicSwap, ok := s.atomicSwaps[req.ID]
	if !ok {
		return nil, ErrUnknownID
	}

	log.Printf("[%s] EnableFunding\n", atomicSwap.ID)

	resp.TxID, err = atomicSwap.EnableFunding(req.AliceRefundNoncePoint, req.RefundSigAlice)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func enableFundingHandler(srv interface{}, ctx context.Context, dec func(interface{}) error,
	interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	if interceptor != nil {
		return nil, ErrNotImplemented
	}

	in := new(EFRequest)
	err := dec(in)
	if err != nil {
		return nil, err
	}

	return srv.(Server).EnableFunding(in)
}

type (
	RADRequest struct {
		ID                   uuid.UUID
		AliceClaimUnlockHash types.UnlockHash
		AliceClaimNoncePoint ed25519.CurvePoint
	}

	RADResponse struct {
		AdaptorDetails *bob.AdaptorDetails
	}
)

func (s *BobServer) RequestAdaptorDetails(req *RADRequest) (*RADResponse, error) {
	var err error
	resp := new(RADResponse)

	s.mutex.Lock()
	defer s.mutex.Unlock()

	atomicSwap, ok := s.atomicSwaps[req.ID]
	if !ok {
		return nil, ErrUnknownID
	}

	log.Printf("[%s] RequestAdaptorDetails\n", atomicSwap.ID)

	resp.AdaptorDetails, err = atomicSwap.RequestAdaptorDetails(req.AliceClaimUnlockHash, req.AliceClaimNoncePoint)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func requestAdaptorDetailsHandler(srv interface{}, ctx context.Context, dec func(interface{}) error,
	interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	if interceptor != nil {
		return nil, ErrNotImplemented
	}

	in := new(RADRequest)
	err := dec(in)
	if err != nil {
		return nil, err
	}

	return srv.(Server).RequestAdaptorDetails(in)
}

type (
	ADRequest struct {
		ID uuid.UUID
	}

	ADResponse struct{}
)

func (s *BobServer) AnnounceDeposit(req *ADRequest) (*ADResponse, error) {
	var err error
	resp := new(ADResponse)

	s.mutex.Lock()
	defer s.mutex.Unlock()

	atomicSwap, ok := s.atomicSwaps[req.ID]
	if !ok {
		return nil, ErrUnknownID
	}

	log.Printf("[%s] AnnounceDeposit\n", atomicSwap.ID)

	err = atomicSwap.AnnounceDeposit()
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func announceDepositHandler(srv interface{}, ctx context.Context, dec func(interface{}) error,
	interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	if interceptor != nil {
		return nil, ErrNotImplemented
	}

	in := new(ADRequest)
	err := dec(in)
	if err != nil {
		return nil, err
	}

	return srv.(Server).AnnounceDeposit(in)
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

func (s *BobServer) Report() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	var keys []string
	for k := range s.atomicSwaps {
		keys = append(keys, k.String())
	}

	sort.Strings(keys)
	for _, uuidStr := range keys {
		k := uuid.Must(uuid.FromString(uuidStr))
		log.Printf("State of %s: %s\n", k, s.atomicSwaps[k].StateText())
	}
}

type Client struct {
	conn *grpc.ClientConn
}

func (c *Client) RequestNonBindingOffer(siacoin types.Currency) (*uuid.UUID, *trader.Offer, error) {
	in := RNBORequest{
		Siacoin: siacoin,
	}
	out := new(RNBOResponse)
	err := grpc.Invoke(context.Background(), "/Roadie/RequestNonBindingOffer", &in, out, c.conn)
	if err != nil {
		return nil, nil, err
	}

	return &out.ID, out.Offer, nil
}

func (c *Client) RequestBindingOffer(id uuid.UUID, antiSpamID big.Int) (*trader.Offer, error) {
	in := RBORequest{
		ID:         id,
		AntiSpamID: antiSpamID,
	}
	out := new(RBOResponse)
	err := grpc.Invoke(context.Background(), "/Roadie/RequestBindingOffer", &in, out, c.conn)
	if err != nil {
		return nil, err
	}

	return out.Offer, nil
}

func (c *Client) AcceptOffer(id uuid.UUID, alicePubKey ed25519.PublicKey) (*bob.RefundDetails, error) {
	in := AORequest{
		ID:          id,
		AlicePubKey: alicePubKey,
	}
	out := new(AOResponse)
	err := grpc.Invoke(context.Background(), "/Roadie/AcceptOffer", &in, out, c.conn)
	if err != nil {
		return nil, err
	}

	return out.RefundDetails, nil
}

func (c *Client) EnableFunding(id uuid.UUID,
	aliceRefundNoncePoint ed25519.CurvePoint, refundSigAlice []byte) (*types.TransactionID, error) {
	in := EFRequest{
		ID:                    id,
		AliceRefundNoncePoint: aliceRefundNoncePoint,
		RefundSigAlice:        refundSigAlice,
	}
	out := new(EFResponse)
	err := grpc.Invoke(context.Background(), "/Roadie/EnableFunding", &in, out, c.conn)
	if err != nil {
		return nil, err
	}

	return out.TxID, nil
}

func (c *Client) RequestAdaptorDetails(id uuid.UUID,
	aliceClaimUnlockHash types.UnlockHash, aliceClaimNoncePoint ed25519.CurvePoint) (*bob.AdaptorDetails, error) {
	in := RADRequest{
		ID:                   id,
		AliceClaimUnlockHash: aliceClaimUnlockHash,
		AliceClaimNoncePoint: aliceClaimNoncePoint,
	}
	out := new(RADResponse)
	err := grpc.Invoke(context.Background(), "/Roadie/RequestAdaptorDetails", &in, out, c.conn)
	if err != nil {
		return nil, err
	}

	return out.AdaptorDetails, nil
}

func (c *Client) AnnounceDeposit(id uuid.UUID) error {
	in := ADRequest{
		ID: id,
	}
	out := new(ADResponse)
	err := grpc.Invoke(context.Background(), "/Roadie/AnnounceDeposit", &in, out, c.conn)
	if err != nil {
		return err
	}

	return nil
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
