package wallet

import (
	"fmt"
	"os/exec"

	"github.com/google/uuid"
)

func CreateWallet() (string, error) {
	walletName := "airdrop-wallet" + uuid.New().String()
	cmd := exec.Command("cardano-cli", "address", "key-gen", "--verification-key-file", walletName+".vkey", "--signing-key-file", walletName+".skey")
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("wallet creation failed: %w", err)
	}

	// Build address
	cmd = exec.Command("cardano-cli", "address", "build",
		"--payment-verification-key-file", walletName+".vkey",
		"--out-file", walletName+".addr",
		"--mainnet")
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("address creation failed: %w", err)
	}

	return walletName + ".addr", nil
}
