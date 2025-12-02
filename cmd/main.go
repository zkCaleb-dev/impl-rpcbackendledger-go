package main

import (
	"context"
	"fmt"

	backends "github.com/stellar/go-stellar-sdk/ingest/ledgerbackend"
)

func main() {

	rpcServerUrl := "https://soroban-testnet.stellar.org"
	rpcBackendOptions := backends.RPCLedgerBackendOptions{
		RPCServerURL: rpcServerUrl,
	}
	ctx := context.Background()
	backend := backends.NewRPCLedgerBackend(rpcBackendOptions)

	// Prepare a single ledger to be ingested,
	err := backend.PrepareRange(ctx, backends.BoundedRange(1900133, 1900133))
	if err != nil {
		fmt.Errorf("Tiro error en PrepareRange")
	}

	// then retrieve it:
	ledger, err := backend.GetLedger(ctx, 1900133)
	if err != nil {
		fmt.Errorf("Tiro error en PrepareRange")
	}

	// Now `ledger` is a raw `xdr.LedgerCloseMeta` object containing the
	// transactions contained within this ledger.
	fmt.Printf("\nHello, Sequence %d. \n", ledger.LedgerSequence())

}
