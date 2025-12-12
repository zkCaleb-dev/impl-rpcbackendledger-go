package rpc

import (
	"context"
	"net/http"

	client "github.com/stellar/go-stellar-sdk/clients/rpcclient"
)

type Config struct {
	Url        string
	HttpClient *http.Client
}

type RpcClient struct {
	client *client.Client
}

func NewRpcClient(config Config) *RpcClient {

	newClient := client.NewClient(config.Url, config.HttpClient)
	rpcClient := RpcClient{
		client: newClient,
	}

	return &rpcClient
}

func (rc RpcClient) GetLatestLedgerSequence(ctx context.Context) (uint32, error) {

	latestLedgerSequence, err := rc.client.GetLatestLedger(ctx)
	if err != nil {
		return 0, err
	}

	return latestLedgerSequence.Sequence, nil

}
