package tx

import (
	"fmt"
	"log/slog"
	"os/exec"

	"github.com/36thchambersoftware/lookout-below/db"
)

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

	slog.Default().Info("Building transaction",
		"wallet_address", walletAddr,
		"holders_count", len(holders),
		"ada_per_nft", adaPerNFT,
		"txOuts", txOuts,
	)

	args := append([]string{
		"conway",
		"transaction", "build",
		"--mainnet",
		"--change-address", walletAddr,
		"--out-file", "airdrop-tx.raw",
	}, txOuts...)

	slog.Default().Info("Executing cardano-cli command", "args", args)

	cmd := exec.Command("cardano-cli", args...)
	cmd.Stdout = nil
	cmd.Stderr = nil

	return cmd.Run()
}

