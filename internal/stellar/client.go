package stellar

import (
	"context"
	"fmt"
	"net/http"

	client "github.com/stellar/go-stellar-sdk/clients/rpcclient"
)

// Client envuelve el RPC client de Stellar para simplificar su uso
type Client struct {
	rpc *client.Client
	url string
}

// NewClient crea una nueva instancia del cliente Stellar
func NewClient(url string) *Client {
	return &Client{
		rpc: client.NewClient(url, &http.Client{}),
		url: url,
	}
}

// NewClientWithHTTP crea un cliente con un http.Client personalizado
// Útil para configurar timeouts, proxies, etc.
func NewClientWithHTTP(url string, httpClient *http.Client) *Client {
	return &Client{
		rpc: client.NewClient(url, httpClient),
		url: url,
	}
}

// GetLatestLedgerSequence obtiene el número de secuencia del ledger más reciente
func (c *Client) GetLatestLedgerSequence(ctx context.Context) (uint32, error) {
	response, err := c.rpc.GetLatestLedger(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get latest ledger: %w", err)
	}
	return response.Sequence, nil
}
