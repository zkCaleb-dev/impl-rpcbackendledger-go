package main

import (
	"context"
	"fmt"
	"io"

	"github.com/stellar/go-stellar-sdk/ingest"
	"github.com/stellar/go-stellar-sdk/support/log"
	"github.com/stellar/go-stellar-sdk/xdr"

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
	var startingSeq uint32 = 1907183
	var ledgersToRead uint32 = 60

	fmt.Printf("Preparing range (%d ledgers)...\n", ledgersToRead)

	// ledgerRange := backends.BoundedRange(startingSeq, startingSeq+ledgersToRead)
	ledgerRangeUnbounded := backends.UnboundedRange(startingSeq)
	err := backend.PrepareRange(ctx, ledgerRangeUnbounded)
	if err != nil {
		fmt.Errorf("Tiro error en PrepareRange")
	}

	// These are the statistics that we're tracking.
	var successfulTransactions, failedTransactions int
	var operationsInSuccessful, operationsInFailed int

	// for seq := startingSeq; seq <= startingSeq+ledgersToRead; seq++ {
	// 	fmt.Printf("Processed ledger %d...\r", seq)

	// 	txReader, err := ingest.NewLedgerTransactionReader(
	// 		ctx, backend, "Test SDF Network ; September 2015", seq,
	// 	)
	// 	if err != nil {
	// 		fmt.Errorf("Failed txReader")
	// 	}
	// 	defer txReader.Close()

	// 	// Read each transaction within the ledger, extract its operations, and
	// 	// accumulate the statistics we're interested in.
	// 	for {
	// 		tx, err := txReader.Read()
	// 		if err == io.EOF {
	// 			break
	// 		}
	// 		if err != nil {
	// 			fmt.Errorf("Failed tx")
	// 		}

	// 		envelope := tx.Envelope
	// 		operationCount := len(envelope.Operations())
	// 		if tx.Result.Successful() {
	// 			successfulTransactions++
	// 			operationsInSuccessful += operationCount
	// 		} else {
	// 			failedTransactions++
	// 			operationsInFailed += operationCount
	// 		}
	// 	}
	// }

	for {
		// fmt.Printf("Processed ledger %d...\r", startingSeq)

		txReader, err := ingest.NewLedgerTransactionReader(ctx, backend, "Test SDF Network ; September 2015", startingSeq)
		if err != nil {
			fmt.Errorf("Failed txReader")
		}
		defer txReader.Close()

		for {
			tx, err := txReader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				fmt.Errorf("Failed tx")
			}

			// Imprimir información de la transacción
			fmt.Printf("Reading tx %d\n", startingSeq)
			fmt.Printf("Transaction Hash: %s\n", tx.Result.TransactionHash.HexString())
			fmt.Printf("Successful: %v\n", tx.Result.Successful())
			fmt.Printf("Operations: %d\n", len(tx.Envelope.Operations()))
			fmt.Println("---")

		}

		startingSeq++

	}

	fmt.Println("\nDone. Results:")
	fmt.Printf("  - total transactions: %d\n", successfulTransactions+failedTransactions)
	fmt.Printf("  - succeeded / failed: %d / %d\n", successfulTransactions, failedTransactions)
	fmt.Printf("  - total operations:   %d\n", operationsInSuccessful+operationsInFailed)
	fmt.Printf("  - succeeded / failed: %d / %d\n", operationsInSuccessful, operationsInFailed)
} // end of main

func hasMatchingMemo(tx ingest.LedgerTransaction, targetName string) bool {
	memo := tx.Envelope.Memo()

	switch memo.Type {
	case xdr.MemoTypeMemoText:
		if string(memo.MustText()) == targetName {
			return true
		}
	case xdr.MemoTypeMemoId:
		// Puedes agregar lógica para IDs si es necesario
	case xdr.MemoTypeMemoHash, xdr.MemoTypeMemoReturn:
		// Comparar hashes si es necesario
	}

	return true
}

// Función para imprimir detalles de la transacción
func printTransactionDetails(tx ingest.LedgerTransaction) {
	fmt.Printf("  Operaciones: %d\n", len(tx.Envelope.Operations()))
	fmt.Printf("  Exitosa: %v\n", tx.Result.Successful())

	// Imprimir cada operación
	for i, op := range tx.Envelope.Operations() {
		fmt.Printf("    Operación %d: %s\n", i, op.Body.Type)
	}
	fmt.Println()
}
