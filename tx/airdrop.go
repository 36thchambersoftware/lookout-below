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

// Send total balance to another wallet.
// cardano-cli conway transaction build-raw \
//   --tx-in c3195ffa0fcc3fc70a681f64bbdfbc52e0eff3f70cbd767714675d088d7b1a14#63 \
//   --tx-out addr1q8ur464mlqsqslh0dn9dqg88zn0q0sqag2hkxc0vhtrn5c7wkhumlr876ehcm8ltdwt7s49mwxfw47c4hcf5p6qdlavqaawfcs+0 \
//   --fee 0 \
//   --out-file tx.raw

//   cardano-cli conway transaction calculate-min-fee \
//   --tx-body-file tx.raw \
//   --tx-in-count 1 \
//   --tx-out-count 1 \
//   --witness-count 1 \
//   --mainnet \
//   --protocol-params-file /home/cardano/cardano/pparams.json

// cardano-cli conway transaction build-raw \
//   --tx-in c3195ffa0fcc3fc70a681f64bbdfbc52e0eff3f70cbd767714675d088d7b1a14#63 \
//   --tx-out addr1q8ur464mlqsqslh0dn9dqg88zn0q0sqag2hkxc0vhtrn5c7wkhumlr876ehcm8ltdwt7s49mwxfw47c4hcf5p6qdlavqaawfcs+1475662 \
//   --fee 165281 \
//   --out-file tx.raw

// cardano-cli conway transaction sign \
//   --tx-body-file tx.raw \
//   --signing-key-file airdrop-walletdd9e0b14-d69f-4dab-9ac4-db5cb18301e5.skey \
//   --mainnet \
//   --out-file tx.signed

// cardano-cli conway transaction submit --tx-file tx.signed --mainnet

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
		txOuts = append(txOuts, "--tx-out", fmt.Sprintf("%s+%d", h.Address, total))
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
	for utxo := range utxos {
		// key is like "txhash#txix"
		txIns = append(txIns, "--tx-in", utxo)
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

