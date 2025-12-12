package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/stellar/go-stellar-sdk/ingest"
	"github.com/stellar/go-stellar-sdk/support/log"
	"github.com/zkCaleb-dev/internal/rpc"

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

	//* Intento pseudo desacoplado.
	rpcConfig := rpc.Config{
		Url:        rpcServerUrl,
		HttpClient: &http.Client{},
	}
	rpcClient := rpc.NewRpcClient(rpcConfig)
	latestLedgerSeq, err := rpcClient.GetLatestLedgerSequence(ctx)
	if err != nil {
		println("Error en GetLatestLedgerSequence() de main: ", err)
	}

	//* Intento con funcion
	// latestLedgerSeq := getLatestLedger(rpcServerUrl, ctx)

	//* Intento estructurado
	// latestLedgerSeq, err := backend.GetLatestLedgerSequence(ctx)
	// if err != nil {
	// 	fmt.Println("Error en PrepareRange:", err)
	// 	return
	// }

	// ledgerRange := backends.BoundedRange(startingSeq, startingSeq+ledgersToRead)
	ledgerRangeUnbounded := ledgerbackend.UnboundedRange(latestLedgerSeq)
	err = backend.PrepareRange(ctx, ledgerRangeUnbounded)
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

// func getLatestLedger(url string, ctx context.Context) uint32 {

// 	// Creamos el cliente.
// 	rpcClient := client.NewClient(url, &http.Client{})

// 	// Con el cliente usamos el metodo GetHealth(), que retorna //* (protocol.GetHealthResponse, error)
// 	// Que tiene las propiedades Status, LatestLedger, OldestLedger, LedgerRetentionWindow
// 	// health, err := rpcClient.GetHealth(ctx)
// 	// if err != nil {
// 	// 	fmt.Println("Error with GetHealth(): ", err)
// 	// }

// 	latestLedgerSeq, err := rpcClient.GetLatestLedger(ctx)
// 	if err != nil {
// 		fmt.Println("Error with GetLatesLedger(): ", err)
// 	}

// 	// Asignamos el valor que necesitamos a una variable.
// 	// latestLedger := health.LatestLedger

// 	// Se imprime el valor.
// 	// fmt.Println("Latest Ledger by getLatestLedger: ", latestLedger)

// 	// Se retorna.
// 	return latestLedgerSeq.Sequence
// }
