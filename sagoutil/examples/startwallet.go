package main

import (
	"log"

	"github.com/Meshbits/shurli/sagoutil"
)

func main() {
	var chain string = "PIRATE"
	var chainParams = []string{"-ac_name=PIRATE", "-ac_supply=0", "-ac_reward=25600000000", "-ac_halving=77777", "-ac_private=1", "-addnode=178.63.77.56"}
	err := sagoutil.StartWallet(chain, chainParams)
	if err != nil {
		log.Printf("StartWallet returned error: %s", err)
	}
}
