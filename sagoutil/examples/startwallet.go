package main

import (
	"github.com/Meshbits/shurli/sagoutil"
)

func main() {
	var chain string = "komodo"
	var chainParams string = "-pubkey=03b7d226d9f9b14e815c138506087eee844271454ae4d0d3168ffed6873de01b89"
	// var params string = ""
	sagoutil.StartWallet(chain, chainParams)
}
