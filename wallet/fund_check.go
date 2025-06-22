package wallet

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os/exec"
)

type AssetAmount struct {
	Unit     string `json:"unit"`   // "lovelace" or asset ID
	Quantity string `json:"quantity"` // string to avoid int overflow
}

type UTxO struct {
	Address string         `json:"address"`
	Value   map[string]int `json:"value"` // keys like "lovelace"
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

	var utxos map[string]UTxO
	if err := json.Unmarshal(output, &utxos); err != nil {
		return fmt.Errorf("failed to decode utxo JSON: %w", err)
	}

	var totalLovelace int
	for _, utxo := range utxos {
		if lovelace, ok := utxo.Value["lovelace"]; ok {
			totalLovelace += lovelace
		}
	}

	slog.Default().Info("Wallet balance", "address", walletAddr, "total_lovelace", totalLovelace)

	if int64(totalLovelace) < requiredLovelace {
		return fmt.Errorf("insufficient funds: required %d lovelace, found %d", requiredLovelace, totalLovelace)
	}
	return nil
}
