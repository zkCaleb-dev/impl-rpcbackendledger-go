package main

import (
	"context"
	"fmt"
	"io"

	"github.com/stellar/go-stellar-sdk/ingest"
	"github.com/stellar/go-stellar-sdk/support/log"

	"github.com/sirupsen/logrus"
	backends "github.com/stellar/go-stellar-sdk/ingest/ledgerbackend"
)

func main() {

	rpcServerUrl := "https://soroban-testnet.stellar.org"
	rpcBackendOptions := backends.RPCLedgerBackendOptions{
		RPCServerURL: rpcServerUrl,
	}
	ctx := context.Background()

	lg := log.New()
	lg.SetLevel(logrus.ErrorLevel)

	backend := backends.NewRPCLedgerBackend(rpcBackendOptions)

	// Prepare range to be ingested
	var startingSeq uint32 = 1900940
	var ledgersToRead uint32 = 60

	fmt.Printf("Preparing range (%d ledgers)...\n", ledgersToRead)
	ledgerRange := backends.BoundedRange(startingSeq, startingSeq+ledgersToRead)
	err := backend.PrepareRange(ctx, ledgerRange)
	if err != nil {
		fmt.Errorf("Tiro error en PrepareRange")
	}

	// These are the statistics that we're tracking.
	var successfulTransactions, failedTransactions int
	var operationsInSuccessful, operationsInFailed int

	for seq := startingSeq; seq <= startingSeq+ledgersToRead; seq++ {
		fmt.Printf("Processed ledger %d...\r", seq)

		txReader, err := ingest.NewLedgerTransactionReader(
			ctx, backend, "Test SDF Network ; September 2015", seq,
		)
		if err != nil {
			fmt.Errorf("Failed txReader")
		}
		defer txReader.Close()

		// Read each transaction within the ledger, extract its operations, and
		// accumulate the statistics we're interested in.
		for {
			tx, err := txReader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				fmt.Errorf("Failed tx")
			}

			envelope := tx.Envelope
			operationCount := len(envelope.Operations())
			if tx.Result.Successful() {
				successfulTransactions++
				operationsInSuccessful += operationCount
			} else {
				failedTransactions++
				operationsInFailed += operationCount
			}
		}
	}

	fmt.Println("\nDone. Results:")
	fmt.Printf("  - total transactions: %d\n", successfulTransactions+failedTransactions)
	fmt.Printf("  - succeeded / failed: %d / %d\n", successfulTransactions, failedTransactions)
	fmt.Printf("  - total operations:   %d\n", operationsInSuccessful+operationsInFailed)
	fmt.Printf("  - succeeded / failed: %d / %d\n", operationsInSuccessful, operationsInFailed)
} // end of main
