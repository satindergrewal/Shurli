package main

import (
	"log"

	"github.com/Meshbits/shurli/sagoutil"
)

func main() {

	err := sagoutil.DLSubJSONData()
	if err != nil {
		log.Println(err)
	}

}
