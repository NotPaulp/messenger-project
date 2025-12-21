package grpc

import (
	"context"
	"fmt"
	"messenger-project/pkg/config"
	pb "messenger-project/proto/user"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn   *grpc.ClientConn
	client pb.UserServiceClient
}

var GlobalClient *Client

func NewClient() *Client {
	cfg := config.Load()
	var err error
	conn, err := grpc.NewClient(cfg.UserGrpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(fmt.Sprintf("New gRPC client %s: %v", cfg.UserGrpcAddr, err))
	}

	return &Client{
		conn:   conn,
		client: pb.NewUserServiceClient(conn),
	}
}

func (c *Client) UserExists(ctx context.Context, username string) (bool, error) {
	resp, err := c.client.GetUser(ctx, &pb.GetUserRequest{Username: username})
	if err != nil {
		return false, err
	}
	return resp.Exists, nil
}

func (c *Client) Close() {
	c.conn.Close()
}
