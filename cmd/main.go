package main

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/stellar/go-stellar-sdk/ingest"
	"github.com/stellar/go-stellar-sdk/ingest/ledgerbackend"
	stellarlog "github.com/stellar/go-stellar-sdk/support/log"
	"github.com/zkCaleb-dev/internal/stellar"
)

const (
	sorobanTestnetURL = "https://soroban-testnet.stellar.org"
	networkPassphrase = "Test SDF Network ; September 2015"
)

func main() {
	ctx := context.Background()

	// Setup loggin
	lg := stellarlog.New()
	lg.SetLevel(logrus.ErrorLevel)

	// Crear cliente de Stellar
	stellarClient := stellar.NewClient(sorobanTestnetURL)

	// Obtener ultimo ledger
	latestSeq, err := stellarClient.GetLatestLedgerSequence(ctx)
	if err != nil {
		fmt.Errorf("getting latest ledger: %w", err)
	}

	fmt.Printf("Starting from ledger: %d\n", latestSeq)

	// Setup backend
	backend := ledgerbackend.NewRPCLedgerBackend(ledgerbackend.RPCLedgerBackendOptions{
		RPCServerURL: sorobanTestnetURL,
	})

	// Preparar rango
	ledgerRange := ledgerbackend.UnboundedRange(latestSeq)
	if err := backend.PrepareRange(ctx, ledgerRange); err != nil {
		fmt.Errorf("preparing ledger range: %w", err)
	}

	// Procesar ledgers
	processLedgers(ctx, backend, latestSeq)
}

func processLedgers(ctx context.Context, backend ledgerbackend.LedgerBackend, startSeq uint32) error {
	currentSeq := startSeq

	for {
		if err := processLedger(ctx, backend, currentSeq); err != nil {
			return fmt.Errorf("processing ledger %d: %w", currentSeq, err)
		}
		currentSeq++
	}
}

func processLedger(ctx context.Context, backend ledgerbackend.LedgerBackend, seq uint32) error {

	ledger, err := backend.GetLedger(ctx, seq)
	if err != nil {
		return fmt.Errorf("getting ledger: %w", err)
	}

	// Obtener ledger
	txReader, err := ingest.NewLedgerTransactionReaderFromLedgerCloseMeta(
		networkPassphrase,
		ledger,
	)
	if err != nil {
		return fmt.Errorf("creating transaction reader: %w", err)
	}
	defer txReader.Close()

	fmt.Printf("Processing sequence: %d\n", txReader.GetSequence())

	tx, err := txReader.Read()
	if err != nil {
		return fmt.Errorf("reading transaction: %w", err)
	}

	// Leer transaccion
	changes, err := tx.GetChanges()
	if err != nil {
		return fmt.Errorf("getting changes: %w", err)
	}

	fmt.Printf("Changes: %v\n", changes)

	return nil
}
