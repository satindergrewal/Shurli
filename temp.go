package main

import (
	"fmt"
	"kmdgo"
	"log"
)

// DEXHandle stores data to get authorised/unauthorised hanldes from all broadcasting traders
type DEXHandle struct {
	Pubkey    string
	Handle    string
	DEXPubkey string
}

func main() {

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

	fmt.Println(len(handles))
	fmt.Println(handles)
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
}
