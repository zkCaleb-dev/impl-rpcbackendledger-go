package main

import (
	"context"
	"fmt"
	"net/http"

	client "github.com/stellar/go-stellar-sdk/clients/rpcclient"
	"github.com/stellar/go-stellar-sdk/ingest"
	"github.com/stellar/go-stellar-sdk/support/log"

	"github.com/sirupsen/logrus"
	"github.com/stellar/go-stellar-sdk/ingest/ledgerbackend"
)

func main() {

	rpcServerUrl := "https://soroban-testnet.stellar.org"
	rpcBackendOptions := ledgerbackend.RPCLedgerBackendOptions{
		RPCServerURL: rpcServerUrl,
	}
	ctx := context.Background()

	lg := log.New()
	lg.SetLevel(logrus.ErrorLevel)

	backend := ledgerbackend.NewRPCLedgerBackend(rpcBackendOptions)

	latestLedgerSeq := getLatestLedger(rpcServerUrl, ctx)
	// latestLedgerSeq, err := backend.GetLatestLedgerSequence(ctx)
	// if err != nil {
	// 	fmt.Println("Error en PrepareRange:", err)
	// 	return
	// }

	// ledgerRange := backends.BoundedRange(startingSeq, startingSeq+ledgersToRead)
	ledgerRangeUnbounded := ledgerbackend.UnboundedRange(latestLedgerSeq)
	err := backend.PrepareRange(ctx, ledgerRangeUnbounded)
	if err != nil {
		fmt.Println("Error en ledgerRangeUnbounded:", err)
		return
	}

	for {
		ledger, err := backend.GetLedger(ctx, latestLedgerSeq)
		if err != nil {
			fmt.Println("Error en ledger:", err)
			return
		}

		transaction, err := ingest.NewLedgerTransactionReaderFromLedgerCloseMeta("Test SDF Network ; September 2015", ledger)
		if err != nil {
			fmt.Println("Error en transaction:", err)
			return
		}
		defer transaction.Close()

		fmt.Println("Actual Sequence: ", transaction.GetSequence())

		read, err := transaction.Read()

		changes, err := read.GetChanges()
		fmt.Println("Changes: ", changes)

		latestLedgerSeq++
	}

}

func getLatestLedger(url string, ctx context.Context) uint32 {

	rpcClient := client.NewClient(url, &http.Client{})
	health, err := rpcClient.GetHealth(ctx)
	if err != nil {
		fmt.Println("Error with GetHealth(): ", err)
	}
	latestLedger := health.LatestLedger

	fmt.Println("Latest Ledger by getLatestLedger: ", latestLedger)

	return latestLedger
}
