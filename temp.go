package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {

	res, err := http.Get("https://raw.githubusercontent.com/jl777/komodo/jl777/src/cc/dapps/subatomic.json")
	if err != nil {
		log.Fatal(err)
	}
	robots, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s", robots)
	ioutil.WriteFile("assets/subatomic.json", robots, 0644)

}
