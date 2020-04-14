package sagoutil

import (
	"fmt"
	"kmdgo"
)

// WalletInfo type stores data to display on Wallet info screen
type WalletInfo struct {
	Ticker  string
	Status  string
	Balance float64
	Blocks  int
	Synced  bool
}

// WalletInfo method returns processed data to display on Dashboard
func WalletInfo(chains kmdgo.AppType) WalletInfo {
	var wallets []WalletInfo

	for _, v := range chains {
		// fmt.Println(i)
		// fmt.Println(v)
		appName := kmdgo.NewAppType(v)

		var info kmdgo.GetInfo

		info, err := appName.GetInfo()
		// fmt.Println(info.Error.Message)
		if err != nil {
			// fmt.Printf("Code: %v\n", info.Error.Code)
			// fmt.Printf("Message: %v\n\n", info.Error.Message)
			if info.Error.Message == "Loading block index..." {
				fmt.Println(v, "- Err happened:", info.Error.Message)
				wallets = append(wallets, WalletInfo{string(v), "Loading...", 0.0, 0, false})
			} else {
				fmt.Println(v, "- Err happened:", err)
				wallets = append(wallets, WalletInfo{string(v), "Offline", 0.0, 0, false})
			}
		} else {
			if info.Error.Message == "connection refused" {
				fmt.Println(v, "- Err happened:", info.Error.Message)
				wallets = append(wallets, WalletInfo{string(v), "Offline", 0.0, 0, false})
			} else {
				// Check status of the blockchain sync
				var tempSyncStatus bool
				if info.Result.Longestchain != info.Result.Blocks {
					tempSyncStatus = false
				} else {
					tempSyncStatus = true
				}

				wallets = append(wallets, WalletInfo{
					Ticker:  info.Result.Name,
					Status:  "Online",
					Balance: info.Result.Balance,
					Blocks:  info.Result.Longestchain,
					Synced:  tempSyncStatus,
				})
			}
		}
	}

	// fmt.Println(wallets)
	return wallets
}
