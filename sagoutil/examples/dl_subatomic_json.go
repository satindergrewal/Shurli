package main

import (
	"log"

	"github.com/satindergrewal/shurli/sagoutil"
)

func main() {

	err := sagoutil.DLSubJSONData()
	if err != nil {
		log.Println(err)
	}

}
