package sagoutil

import (
	"errors"
	"fmt"
	"io/ioutil"
	"kmdgo"
	"log"
	"net/http"
)

// WInfo type stores data to display on Wallet info screen
type WInfo struct {
	Ticker  string
	Status  string
	Balance float64
	Blocks  int
	Synced  bool
}

// WalletInfo method returns processed data to display on Dashboard
func WalletInfo(chains []kmdgo.AppType) []WInfo {
	var wallets []WInfo

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
				wallets = append(wallets, WInfo{string(v), "Loading...", 0.0, 0, false})
			} else {
				fmt.Println(v, "- Err happened:", err)
				wallets = append(wallets, WInfo{string(v), "Offline", 0.0, 0, false})
			}
		} else {
			if info.Error.Message == "connection refused" {
				fmt.Println(v, "- Err happened:", info.Error.Message)
				wallets = append(wallets, WInfo{string(v), "Offline", 0.0, 0, false})
			} else {
				// Check status of the blockchain sync
				var tempSyncStatus bool
				if info.Result.Longestchain != info.Result.Blocks {
					tempSyncStatus = false
				} else {
					tempSyncStatus = true
				}

				wallets = append(wallets, WInfo{
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

// DEXHandle stores data to get authorised/unauthorised hanldes from all broadcasting traders
type DEXHandle struct {
	Pubkey     string
	Handle     string
	DEXPubkey  string
	authorised bool
}

// DEXHandles returns public address's public Key, DEX CC specific public key, and handle data set picked from:
// - DEX_list handles API
// - subatomic.json file taken from the source code of komodo code
// The purpose of this function is to return the data about authorised and unauthorised handles and show that in orderbook or orderlist in the GUI application.
func DEXHandles() []DEXHandle {
	var handles []DEXHandle

	// handles = append(handles, DEXHandle{"03732f8ef851ff234c74d0df575c2c5b159e2bab3faca4ec52b3f217d5cda5361d", "satinder", "01b5d5b1991152fd45e4ba7005a5a752c2018634a9a6cdeb06b633e731e7b5f46b"})
	// handles = append(handles, DEXHandle{"03732f8ef851ff234c74d0df575c2c5b159e2bab3faca4ec52b3f217d5cda5361d", "satinder", "01b5d5b1991152fd45e4ba7005a5a752c2018634a9a6cdeb06b633e731e7b5f46b"})

	var appName kmdgo.AppType
	appName = `DEX`

	var list kmdgo.DEXList

	args := make(kmdgo.APIParams, 3)
	// stopat
	args[0] = "0"
	// minpriority
	args[1] = "0"
	// tagA
	args[2] = "handles"
	// fmt.Println(args)

	list, err := appName.DEXList(args)
	if err != nil {
		fmt.Printf("Code: %v\n", list.Error.Code)
		fmt.Printf("Message: %v\n\n", list.Error.Message)
		log.Fatalln("Err happened", err)
	}

	var tmpPubkey string
	tmpPubkey = ""

	for _, v := range list.Result.Matches {
		// fmt.Printf("\n-------\n")
		// fmt.Println(i)

		if tmpPubkey == v.Decrypted {
			// fmt.Println("Temp Pubkey matched")
		} else {
			// fmt.Println("Temp Pubkey did not match\nUpdated it's value")
			tmpPubkey = v.Decrypted
			handles = append(handles, DEXHandle{
				Pubkey:    v.Decrypted,
				Handle:    v.TagB,
				DEXPubkey: v.Pubkey,
			})
		}

	}

	// fmt.Println(len(handles))
	// fmt.Println(handles)
	// fmt.Println(handles[0])
	// fmt.Println(handles[1])

	dexpubkey := "01b5d5b1991152fd45e4ba7005a5a752c2018634a9a6cdeb06b633e731e7b5f46b"
	var handle string
	// var authorised bool

	for _, value := range handles {
		// fmt.Println(value.DEXPubkey)
		if value.DEXPubkey == dexpubkey {
			handle = value.Handle
		}
	}

	fmt.Println(handle)
	return handles
}

// DLSubJSONData downloads the latest updated subatomic.json file data and saves it to "assets/subatomic.json" file in subatomicgo app
func DLSubJSONData() error {
	res, err := http.Get("https://raw.githubusercontent.com/jl777/komodo/jl777/src/cc/dapps/subatomic.json")
	if err != nil {
		log.Fatal(err)
	}
	subJSONData, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	if res.Status == "200 OK" {
		ioutil.WriteFile("assets/subatomic.json", subJSONData, 0644)
		return nil
	} else {
		// fmt.Println(res.Status)
		return errors.New(res.Status)
	}
}
