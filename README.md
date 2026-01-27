# Stellar Ledger Ingest - Go Implementation

An educational implementation of a client that reads ledgers in real-time from the Stellar blockchain using the official Go SDK.

## What is this project?

This repository demonstrates how to connect to the Stellar network (testnet) and process ledgers continuously using the **RPC Backend** from the Go SDK. It's ideal for learning about:

- How the ledger system works in Stellar
- Using the official Stellar Go SDK
- Architecture patterns for blockchain applications
- Real-time data ingestion from a blockchain

## Key Concepts

### What is Stellar?

[Stellar](https://stellar.org) is a decentralized blockchain network designed for fast and low-cost payments. Unlike other blockchains, Stellar uses a consensus protocol called **Stellar Consensus Protocol (SCP)** that allows transactions to finalize in 3-5 seconds.

### What is a Ledger?

In Stellar, a **ledger** is the equivalent of a "block" in other blockchains. Each ledger contains:
- A unique sequence number
- A set of transactions
- The resulting state changes (balances, contracts, etc.)

Ledgers are created approximately every 5 seconds on the Stellar network.

### What is Soroban?

[Soroban](https://soroban.stellar.org) is Stellar's smart contract platform. The **Stellar RPC** (formerly called Soroban RPC) is the server that allows interaction with the network and retrieval of ledger data.

### What is the Ingest SDK?

The [Ingest SDK](https://github.com/stellar/go-stellar-sdk) is part of the Go SDK that allows you to:
- Read historical and real-time ledgers
- Parse transactions and their effects
- Extract state changes from the blockchain

## Project Architecture

```
ImplRpcBackendLedgerGo/
├── cmd/
│   └── main.go              # Application entry point
├── internal/
│   └── stellar/
│       └── client.go        # Stellar RPC client wrapper
├── bin/                     # Compiled binary
├── Makefile                 # Build commands
├── go.mod                   # Go module definition
└── go.sum                   # Dependency checksums
```

### Component Explanation

#### `cmd/main.go`

The entry point that:
1. Initializes the RPC client connecting to testnet
2. Gets the most recent ledger number
3. Configures the backend to read ledgers continuously
4. Processes each ledger extracting transactions and changes

```go
// Main flow
stellarClient := stellar.NewClient(sorobanTestnetURL)  // 1. Create client
latestSeq, _ := stellarClient.GetLatestLedgerSequence(ctx)  // 2. Get latest ledger
backend := ledgerbackend.NewRPCLedgerBackend(...)  // 3. Configure backend
processLedgers(ctx, backend, latestSeq)  // 4. Process in infinite loop
```

#### `internal/stellar/client.go`

A wrapper that encapsulates the official RPC client, following the **Adapter pattern**:

```go
type Client struct {
    rpc *client.Client
    url string
}

func (c *Client) GetLatestLedgerSequence(ctx context.Context) (uint32, error) {
    response, err := c.rpc.GetLatestLedger(ctx)
    return response.Sequence, nil
}
```

This pattern allows:
- Simplifying the API for specific use cases
- Facilitating testing with mocks
- Decoupling code from the SDK implementation

## How to Run

### Prerequisites

- Go 1.25 or higher
- Internet connection (to connect to testnet)

### Installation

```bash
# Clone the repository
git clone https://github.com/zkCaleb-dev/ImplRpcBackendLedgerGo.git
cd ImplRpcBackendLedgerGo

# Download dependencies
go mod download
```

### Running

Using the Makefile:

```bash
# Run directly (development)
make run

# Build binary
make build

# Run the compiled binary
./bin/main
```

Or directly with Go:

```bash
go run cmd/main.go
```

### Expected Output

```
Starting from ledger: 1234567
Processing sequence: 1234567
Changes: [...]
Processing sequence: 1234568
Changes: [...]
...
```

## Main Dependencies

| Package | Purpose |
|---------|---------|
| `github.com/stellar/go-stellar-sdk` | Official Stellar SDK for Go |
| `github.com/stellar/go-stellar-sdk/ingest` | Ledger reading and parsing |
| `github.com/stellar/go-stellar-sdk/ingest/ledgerbackend` | RPC backend for fetching ledgers |
| `github.com/stellar/go-stellar-sdk/clients/rpcclient` | HTTP client for the RPC |
| `github.com/sirupsen/logrus` | Structured logging |

## Configuration

Currently the configuration is hardcoded in `cmd/main.go`:

```go
const (
    sorobanTestnetURL = "https://soroban-testnet.stellar.org"
    networkPassphrase = "Test SDF Network ; September 2015"
)
```

| Variable | Description | Default Value |
|----------|-------------|---------------|
| `sorobanTestnetURL` | RPC server URL | Stellar Testnet |
| `networkPassphrase` | Network identifier | Testnet |

To connect to **mainnet**, change to:
```go
const (
    sorobanTestnetURL = "https://soroban.stellar.org"
    networkPassphrase = "Public Global Stellar Network ; September 2015"
)
```

## Code Concepts Explained

### UnboundedRange vs BoundedRange

```go
// UnboundedRange: Reads from an initial ledger into the future indefinitely
ledgerRange := ledgerbackend.UnboundedRange(latestSeq)

// BoundedRange: Reads a specific range of ledgers
ledgerRange := ledgerbackend.BoundedRange(startSeq, endSeq)
```

### Transaction Reader

```go
txReader, _ := ingest.NewLedgerTransactionReaderFromLedgerCloseMeta(
    networkPassphrase,
    ledger,
)

// Read transactions one by one
tx, _ := txReader.Read()

// Get the state changes produced by the transaction
changes, _ := tx.GetChanges()
```

### State Changes

The `changes` represent modifications to the blockchain state:
- Account creation/modification
- Balance changes
- Soroban contract updates
- Modifications to trustlines, offers, etc.

## Additional Resources

### Official Documentation

- [Stellar Developer Docs](https://developers.stellar.org)
- [Go SDK (go-stellar-sdk)](https://github.com/stellar/go-stellar-sdk)
- [Stellar RPC API](https://developers.stellar.org/docs/data/apis/rpc)
- [Event Ingestion Guide](https://developers.stellar.org/docs/build/guides/events/ingest)

### Related Tutorials

- [Introducing the Golang Stellar SDK](https://stellar.org/blog/developers/introducing-the-golang-stellar-sdk)
- [pkg.go.dev - SDK Documentation](https://pkg.go.dev/github.com/stellar/go)

## Suggested Next Steps

If you want to expand this project, consider:

1. **Add environment-based configuration**: Use environment variables or config files
2. **Implement filters**: Filter transactions by type, account, or contract
3. **Persistence**: Save processed data to a database
4. **Metrics**: Add observability with Prometheus/OpenTelemetry
5. **Graceful shutdown**: Handle system signals for clean termination

## Contributing

Contributions are welcome. Please:

1. Fork the repository
2. Create a branch for your feature (`git checkout -b feature/new-feature`)
3. Commit your changes (`git commit -m 'Add new feature'`)
4. Push to the branch (`git push origin feature/new-feature`)
5. Open a Pull Request

## License

This project is open source and available for educational use.

---

Created for educational purposes to learn about blockchain application development with Stellar.
