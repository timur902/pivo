package orderclient

import (
	"beer/pkg/grpcmiddleware"
	"beer/proto/order"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Dial(target string) (*grpc.ClientConn, orderpb.OrderServiceClient, error) {
	conn, err := grpc.NewClient(
		target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(grpcmiddleware.UnaryClientLogger()),
	)
	if err != nil {
		return nil, nil, err
	}
	return conn, orderpb.NewOrderServiceClient(conn), nil
}
