package main

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strconv"

	"github.com/36thchambersoftware/lookout-below/config"
	"github.com/36thchambersoftware/lookout-below/db"
	"github.com/36thchambersoftware/lookout-below/tx"
	"github.com/36thchambersoftware/lookout-below/wallet"
	_ "github.com/lib/pq"
)

func main() {
	config.LoadEnv()

	// dbURL := config.GetEnv("DB_URL")
	// dbConn, err := sql.Open("postgres", dbURL)
	// if err != nil {
	// 	log.Fatalf("DB connection failed: %v", err)
	// }

	// policyID := "d47ae93a752e0e38abe793c84a7320fa959ce23a795c22f52ff73523"
	// holders, err := db.GetHolders(dbConn, policyID)
	// if err != nil {
	// 	log.Fatalf("Failed to fetch holders: %v", err)
	// }
	holdersFile := "holders.json"
	bytes, err := os.ReadFile(holdersFile)
	if err != nil {
		log.Fatalf("Failed to read holders file: %v", err)
	}
	var holders []db.Holder
	if err = json.Unmarshal(bytes, &holders); err != nil {
		log.Fatalf("Failed to parse holders JSON: %v", err)
	}
	// Ensure holders are not empty
	if len(holders) == 0 {
		log.Fatal("No holders found in the provided file")
	}

	walletAddrFile, err := wallet.CreateWallet()
	if err != nil {
		log.Fatalf("Wallet creation failed: %v", err)
	}
	addrBytes, _ := os.ReadFile(walletAddrFile)
	walletAddr := string(addrBytes)

	var totalNeeded int64
	adaPerNFT, err := strconv.ParseInt(config.GetEnv("AIR_DROP_AMOUNT"), 10, 64)
	if err != nil {
		log.Fatalf("Invalid AIR_DROP_AMOUNT: %v", err)
	}

	slog.Default().Info("Airdrop configuration",
		"wallet_address", walletAddr,
		"ada_per_nft", adaPerNFT,
		"holders_count", len(holders),
	)
	for _, h := range holders {
		totalNeeded += h.Count * adaPerNFT
	}

	if err := wallet.CheckFunds(walletAddr, int64(totalNeeded*1e6)); err != nil {
		log.Fatalf("Funding error: %v", err)
	}

	if err := tx.BuildTransaction(walletAddr, holders, adaPerNFT); err != nil {
		log.Fatalf("Transaction build failed: %v", err)
	}

	fmt.Println("Transaction draft built as airdrop-tx.raw")
}
