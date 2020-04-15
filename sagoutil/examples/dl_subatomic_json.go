package main

import (
	"log"

	"github.com/satindergrewal/subatomicgo/sagoutil"
)

func main() {

	err := sagoutil.DLSubJSONData()
	if err != nil {
		log.Println(err)
	}

}
