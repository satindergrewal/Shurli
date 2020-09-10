package main

import (
	"fmt"

	"github.com/Meshbits/shurli/sagoutil"
	"github.com/satindergrewal/kmdgo/kmdutil"
)

func main() {
	// var chain string = "PIRATE"
	// var chainParams = []string{"-ac_name=PIRATE", "-ac_supply=0", "-ac_reward=25600000000", "-ac_halving=77777", "-ac_private=1", "-addnode=178.63.77.56"}
	// err := sagoutil.StartWallet(chain, chainParams)
	// if err != nil {
	// 	log.Printf("StartWallet returned error: %s", err)
	// }

	// backupDir := sagoutil.ShurliRootDir()
	// sagoutil.BackupConfigJSON(backupDir)

	// sagoutil.GenerateDEXP2PAccount()
	// sagoutil.ImportTAddrPrivKey(`Komodo`)
	// sagoutil.ImportZAddrPrivKey(`PIRATE`)

	// Testing DexP2P Account update function
	// var accountData sagoutil.SubAtomicConfig
	// accountData.DexHandle = "Satinder"

	// err := sagoutil.UpdateDEXP2PAccount(accountData)
	// if err != nil {
	// 	fmt.Printf("%s", err)
	// }

	var wallet kmdutil.IguanaWallet
	wallet, err := sagoutil.GenerateDEXP2PWallet()
	if err != nil {
		fmt.Printf("%s", err)
	}
	fmt.Println(wallet)
}
