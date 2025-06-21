package tx

import (
	"fmt"
	"log/slog"
	"os/exec"

	"github.com/36thchambersoftware/lookout-below/db"
)

func BuildTransaction(walletAddr string, holders []db.Holder, adaPerNFT float64) error {
	txOuts := []string{}
	for _, h := range holders {
		total := adaPerNFT * float64(h.Count)
		txOuts = append(txOuts, fmt.Sprintf("--tx-out %s+%d", h.Address, int(total*1e6)))
	}

	slog.Default().Info("Building transaction",
		"wallet_address", walletAddr,
		"holders_count", len(holders),
		"ada_per_nft", adaPerNFT,
		"txOuts", txOuts,
	)

	args := append([]string{
		"transaction", "build",
		"--alonzo-era", "--mainnet",
		"--change-address", walletAddr,
		"--out-file", "airdrop-tx.raw",
	}, txOuts...)

	cmd := exec.Command("cardano-cli", args...)
	cmd.Stdout = nil
	cmd.Stderr = nil

	return cmd.Run()
}

