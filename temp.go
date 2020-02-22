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

// Wallets store array of WalletInfo data type
type Wallets struct {
	Wallet []WalletInfo
}

func main() {

	var wallets Wallets

	for i, v := range chains {
		// fmt.Println(i)
		// fmt.Println(v)
		appName := kmdgo.NewAppType(v)

		var info kmdgo.GetInfo

		info, err := appName.GetInfo()
		fmt.Println(info)
		if err != nil {
			// fmt.Printf("Code: %v\n", info.Error.Code)
			// fmt.Printf("Message: %v\n\n", info.Error.Message)
			fmt.Println("Err happened", err)
			wallets.Wallet[i].Status = "Offline"
		} else {
			fmt.Println(i, info.Result.Name)
			// wallets.Wallet[i].Ticker = info.Result.Name
			// wallets.Wallet[i].Status = "Online"
			// wallets.Wallet[i].Balance = info.Result.Balance
			// wallets.Wallet[i].Blocks = info.Result.Longestchain
			// if info.Result.Longestchain != info.Result.Blocks {
			// 	wallets.Wallet[i].Synched = false
			// } else {
			// 	wallets.Wallet[i].Synched = true
			// }
		}
	}

	fmt.Println(wallets)
}
