package main

import (
	"fmt"
	"kmdgo"
	"subatomicgo/sagoutil"
)

func main() {
	var chains = []kmdgo.AppType{"komodo", "PIRATE"}

	var wallets []sagoutil.WInfo
	wallets = sagoutil.WalletInfo(chains)

	fmt.Println(wallets)

	// sagoutil.Test()
}
