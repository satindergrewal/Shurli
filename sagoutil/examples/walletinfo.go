package main

import (
	"fmt"

	"github.com/satindergrewal/kmdgo"

	"github.com/Meshbits/shurli/sagoutil"
)

func main() {
	var chains = []kmdgo.AppType{"komodo", "PIRATE"}

	var wallets []sagoutil.WInfo
	wallets = sagoutil.WalletInfo(chains)

	fmt.Println(wallets)

}
