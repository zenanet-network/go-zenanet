package harmoniagrpc

import (
	"time"

	"github.com/zenanet-network/go-zenanet/log"

	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	proto "github.com/maticnetwork/polyproto/heimdall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	stateFetchLimit = 50
)

type HarmoniaGRPCClient struct {
	conn   *grpc.ClientConn
	client proto.HarmoniaClient
}

func NewHarmoniaGRPCClient(address string) *HarmoniaGRPCClient {
	opts := []grpc_retry.CallOption{
		grpc_retry.WithMax(10000),
		grpc_retry.WithBackoff(grpc_retry.BackoffLinear(5 * time.Second)),
		grpc_retry.WithCodes(codes.Internal, codes.Unavailable, codes.Aborted, codes.NotFound),
	}

	conn, err := grpc.NewClient(address,
		grpc.WithStreamInterceptor(grpc_retry.StreamClientInterceptor(opts...)),
		grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(opts...)),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Crit("Failed to connect to Harmonia gRPC", "error", err)
	}

	log.Info("Connected to Harmonia gRPC server", "address", address)

	return &HarmoniaGRPCClient{
		conn:   conn,
		client: proto.NewHarmoniaClient(conn),
	}
}

func (h *HarmoniaGRPCClient) Close() {
	log.Debug("Shutdown detected, Closing Harmonia gRPC client")
	h.conn.Close()
}
