package prismaClient

import (
    "context"

    "github.com/steebchen/prisma-client-go/runtime/transaction"
    prisma "github.com/ZXyoukai/InfraWatch-Backend/prisma/client"
)

type Client struct{
    Client *prisma.PrismaClient
}

func NewClient(ctx context.Context) (*Client, error) {
    p := prisma.NewClient()
    if err := p.Connect(); err != nil {
        return nil, err
    }
    // optional: test connection
    _ = transaction.NewTxManager(nil)
    return &Client{Client: p}, nil
}

func (c *Client) Disconnect() error {
    return c.Client.Disconnect()
}