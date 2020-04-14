package main

import (
	"fmt"
	"kmdgo"
	"log"
)

// DEXHandles stores data to get authorised/unauthorised hanldes from all broadcasting traders
type DEXHandles struct {
	Pubkey    string
	Handle    string
	DEXPubkey string
}

func main() {

	var handles []DEXHandles

	// handles = append(handles, DEXHandles{"03732f8ef851ff234c74d0df575c2c5b159e2bab3faca4ec52b3f217d5cda5361d", "satinder", "01b5d5b1991152fd45e4ba7005a5a752c2018634a9a6cdeb06b633e731e7b5f46b"})
	// handles = append(handles, DEXHandles{"03732f8ef851ff234c74d0df575c2c5b159e2bab3faca4ec52b3f217d5cda5361d", "satinder", "01b5d5b1991152fd45e4ba7005a5a752c2018634a9a6cdeb06b633e731e7b5f46b"})

	var appName kmdgo.AppType
	appName = `DEX`

	var list kmdgo.DEXList

	args := make(kmdgo.APIParams, 10)
	// stopat
	args[0] = "0"
	// minpriority
	args[1] = "0"
	// tagA
	args[2] = "handles"
	// tagB
	args[3] = ""
	// pubkey33
	args[4] = ""
	// minA
	args[5] = ""
	// maxA
	args[6] = ""
	// minB
	args[7] = ""
	// maxB
	args[8] = ""
	// stophash
	args[9] = ""
	fmt.Println(args)

	list, err := appName.DEXList(args)
	if err != nil {
		fmt.Printf("Code: %v\n", list.Error.Code)
		fmt.Printf("Message: %v\n\n", list.Error.Message)
		log.Fatalln("Err happened", err)
	}

	var tmpPubkey string
	tmpPubkey = ""

	for i, v := range list.Result.Matches {
		fmt.Printf("\n-------\n")
		fmt.Println(i)
		// fmt.Println("Timestamp", v.Timestamp)
		// fmt.Println("ID", v.ID)
		// fmt.Println("Hash", v.Hash)
		// fmt.Println("TagA", v.TagA)
		// fmt.Println("TagB", v.TagB)
		// fmt.Println("Pubkey", v.Pubkey)
		// fmt.Println("Payload", v.Payload)
		// fmt.Println("Hex", v.Hex)
		// fmt.Println("Decrypted", v.Decrypted)
		// fmt.Println("Decryptedhex", v.Decryptedhex)
		// fmt.Println("Senderpub", v.Senderpub)
		// fmt.Println("Error", v.Error)
		// fmt.Println("AmountA", v.AmountA)
		// fmt.Println("AmountB", v.AmountB)
		// fmt.Println("Priority", v.Priority)
		// fmt.Println("Recvtime", v.Recvtime)
		// fmt.Println("Cancelled", v.Cancelled)

		if tmpPubkey == v.Decrypted {
			fmt.Println("Temp Pubkey matched")
		} else {
			fmt.Println("Temp Pubkey did not match")
		}

		handles = append(handles, DEXHandles{
			Pubkey:    v.Decrypted,
			Handle:    v.TagB,
			DEXPubkey: v.Pubkey,
		})
	}

	// for i, v := range chains {
	// 	// fmt.Println(i)
	// 	// fmt.Println(v)
	// 	appName := kmdgo.NewAppType(v)

	// 	var info kmdgo.GetInfo

	// 	info, err := appName.GetInfo()
	// 	fmt.Println(info)
	// 	if err != nil {
	// 		// fmt.Printf("Code: %v\n", info.Error.Code)
	// 		// fmt.Printf("Message: %v\n\n", info.Error.Message)
	// 		fmt.Println(v, "- Err happened:", err)
	// 		fmt.Println("wallets.Wallet[", i, "].Status")
	// 		wallets = append(wallets, WalletInfo{string(v), "Offline", 0.0, 0, false})
	// 	} else {

	// 		// Check status of the blockchain sync
	// 		var tempSyncStatus bool
	// 		if info.Result.Longestchain != info.Result.Blocks {
	// 			tempSyncStatus = false
	// 		} else {
	// 			tempSyncStatus = true
	// 		}

	// 		wallets = append(wallets, WalletInfo{
	// 			Ticker:  info.Result.Name,
	// 			Status:  "Online",
	// 			Balance: info.Result.Balance,
	// 			Blocks:  info.Result.Longestchain,
	// 			Synched: tempSyncStatus,
	// 		})
	// 	}
	// }

	// appName := kmdgo.NewAppType(`komodo`)

	// var info kmdgo.GetInfo

	// info, err := appName.GetInfo()
	// fmt.Println(info)

	// if err != nil {
	// 	fmt.Println("Err happened", err)
	// }

	// fmt.Println(len(handles))
	// fmt.Println(handles)
	fmt.Println(handles[0])
	fmt.Println(handles[1])
}
