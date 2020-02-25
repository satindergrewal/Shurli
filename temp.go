package main

import (
	"fmt"
	"kmdgo"
)

var chains = []kmdgo.AppType{"komodo", "PIRATE"}

// WalletInfo stores data to display on Wallet info screen
type WalletInfo struct {
	Ticker  string
	Status  string
	Balance float64
	Blocks  int
	Synched bool
}

func main() {

	var wallets []WalletInfo

	wallets = append(wallets, WalletInfo{"KMD", "Offline", 0.3, 2, true})
	wallets = append(wallets, WalletInfo{"PIRATE", "Online", 4.5, 56, false})

	// for i, v := range chains {
	// 	fmt.Println(i)
	// 	fmt.Println(wallets.Wallet)
	// 	// fmt.Println(v)
	// 	appName := kmdgo.NewAppType(v)

	// 	var info kmdgo.GetInfo

	// 	info, err := appName.GetInfo()
	// 	fmt.Println(info)
	// 	if err != nil {
	// 		// fmt.Printf("Code: %v\n", info.Error.Code)
	// 		// fmt.Printf("Message: %v\n\n", info.Error.Message)
	// 		fmt.Println("Err happened", err)
	// 		fmt.Println("wallets.Wallet[", i, "].Status")
	// 		// wallets.Wallet[i].Status = "Offline"
	// 	} else {
	// 		fmt.Println(i, info.Result.Name)
	// 		// wallets.Wallet[i].Ticker = info.Result.Name
	// 		fmt.Println(info.Result.Name)
	// 		// wallets.Wallet[i].Status = "Online"
	// 		// wallets.Wallet[i].Balance = info.Result.Balance
	// 		// wallets.Wallet[i].Blocks = info.Result.Longestchain
	// 		// if info.Result.Longestchain != info.Result.Blocks {
	// 		// 	wallets.Wallet[i].Synched = false
	// 		// } else {
	// 		// 	wallets.Wallet[i].Synched = true
	// 		// }
	// 	}
	// }

	appName := kmdgo.NewAppType(`komodo`)

	var info kmdgo.GetInfo

	info, err := appName.GetInfo()
	fmt.Println(info)

	if err != nil {
		fmt.Println("Err happened", err)
	}

	fmt.Println(len(wallets))
	fmt.Println(wallets[0])
}
