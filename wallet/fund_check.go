package wallet

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os/exec"
	"strconv"
)

type AssetAmount struct {
	Unit     string `json:"unit"`   // "lovelace" or asset ID
	Quantity string `json:"quantity"` // string to avoid int overflow
}

type UTxO struct {
	TxHash   string         `json:"tx_hash"`
	TxIndex  int            `json:"tx_index"`
	Amount   []AssetAmount  `json:"amount"`
}

func CheckFunds(walletAddr string, requiredLovelace int64) error {
	// Run `cardano-cli query utxo` with JSON output
	cmd := exec.Command("cardano-cli", "query", "utxo",
		"--address", walletAddr,
		"--mainnet",
		"--out-file", "/dev/stdout",
		"--output-json",
	)

	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to run cardano-cli query utxo: %w", err)
	}

	// Parse the JSON into a map of tx_in => UTxO
	var utxos map[string]UTxO
	if err := json.Unmarshal(output, &utxos); err != nil {
		return fmt.Errorf("failed to parse UTXO JSON: %w", err)
	}

	var totalLovelace int64
	for _, utxo := range utxos {
		for _, amt := range utxo.Amount {
			if amt.Unit == "lovelace" {
				q, err := strconv.ParseInt(amt.Quantity, 10, 64)
				if err != nil {
					return fmt.Errorf("invalid lovelace quantity: %w", err)
				}
				totalLovelace += q
			}
		}
	}

	slog.Default().Info("Wallet balance", "address", walletAddr, "total_lovelace", totalLovelace)

	if totalLovelace < requiredLovelace {
		return fmt.Errorf("insufficient funds: required %d lovelace, found %d", requiredLovelace, totalLovelace)
	}
	return nil
}
