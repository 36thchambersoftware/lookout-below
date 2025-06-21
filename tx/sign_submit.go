// tx/sign_submit.go
package tx

import (
	"fmt"
	"os/exec"
)

// SignTransaction signs the raw transaction using the given signing key file.
func SignTransaction(signingKeyFile string) error {
	cmd := exec.Command("cardano-cli", "transaction", "sign",
		"--tx-body-file", "airdrop-tx.raw",
		"--signing-key-file", signingKeyFile,
		"--mainnet",
		"--out-file", "airdrop-tx.signed",
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("signing failed: %v\n%s", err, output)
	}
	return nil
}

// SubmitTransaction submits the signed transaction to the blockchain.
func SubmitTransaction() error {
	cmd := exec.Command("cardano-cli", "transaction", "submit",
		"--tx-file", "airdrop-tx.signed",
		"--mainnet",
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("submission failed: %v\n%s", err, output)
	}
	return nil
}
