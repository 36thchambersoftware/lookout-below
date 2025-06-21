package wallet

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

func CheckFunds(addr string, minLovelace int64) error {
	cmd := exec.Command("cardano-cli", "query", "utxo", "--address", addr, "--mainnet")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to query UTXO: %w", err)
	}

	total := int64(0)
	lines := strings.Split(out.String(), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) >= 3 {
			amount, _ := strconv.ParseInt(fields[2], 10, 64)
			total += amount
		}
	}

	if total < minLovelace {
		return fmt.Errorf("insufficient funds: have %d, need %d", total, minLovelace)
	}
	return nil
}
