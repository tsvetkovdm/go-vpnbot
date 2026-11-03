package core_xray

import (
	"fmt"
	"time"

	"github.com/xtls/xray-core/app/proxyman/command"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"
)

type User struct {
	Email string
	ID    string
}

type Client struct {
	conn    *grpc.ClientConn
	handler command.HandlerServiceClient
	timeout time.Duration
}

func NewClient(config Config) (*Client, error) {
	conn, err := grpc.NewClient(
		config.APIAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithConnectParams(grpc.ConnectParams{
			Backoff: backoff.Config{
				BaseDelay:  time.Second,
				Multiplier: 1.6,
				Jitter:     0.2,
				MaxDelay:   5 * time.Second,
			},
			MinConnectTimeout: 5 * time.Second,
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("new grpc client %s: %w", config.APIAddr, err)
	}

	return &Client{
		conn:    conn,
		handler: command.NewHandlerServiceClient(conn),
		timeout: config.Timeout,
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}
