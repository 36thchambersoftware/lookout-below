package tx

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"os/exec"

	"github.com/36thchambersoftware/lookout-below/db"
)

type UTxOMap map[string]struct {
		TxHash string `json:"tx_hash"`
		TxIx   int    `json:"tx_index"`
	}

func BuildTransaction(walletAddr string, holders []db.Holder, adaPerNFT int64) error {
	txOuts := []string{}
	for _, h := range holders {
		total := adaPerNFT * h.Quantity
		slog.Default().Info("Calculating transaction output",
			"holder_address", h.Address,
			"nft_count", h.Quantity,
			"ada_per_nft", adaPerNFT,
			"total_ada", total,
		)
		txOuts = append(txOuts, fmt.Sprintf("--tx-out %s+%d", h.Address, total))
	}

	// 1. Query UTXOs in JSON format
	utxoCmd := exec.Command("cardano-cli", "query", "utxo",
		"--address", walletAddr,
		"--mainnet",
		"--out-file", "/dev/stdout",
		"--output-json",
	)

	utxoOutput, err := utxoCmd.Output()
	if err != nil {
		return fmt.Errorf("failed to query UTXOs: %w", err)
	}

	// 2. Parse JSON into map
	var utxos UTxOMap
	if err := json.Unmarshal(utxoOutput, &utxos); err != nil {
		return fmt.Errorf("failed to parse UTXO JSON: %w", err)
	}

	// 3. Construct --tx-in arguments
	txIns := []string{}
	for key := range utxos {
		// key is like "txhash#txix"
		txIns = append(txIns, fmt.Sprintf("--tx-in=%s", key))
	}
	if len(txIns) == 0 {
		return fmt.Errorf("no UTXOs found at address %s", walletAddr)
	}

	slog.Default().Info("Building transaction",
		"wallet_address", walletAddr,
		"holders_count", len(holders),
		"ada_per_nft", adaPerNFT,
		"txOuts", txOuts,
	)

	// 4. Build full CLI command
	args := append([]string{
		"conway", "transaction", "build",
		"--mainnet",
		"--change-address", walletAddr,
		"--out-file", "airdrop-tx.raw",
	}, append(txIns, txOuts...)...)

	slog.Default().Info("Executing cardano-cli command", "args", args)

	cmd := exec.Command("cardano-cli", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

