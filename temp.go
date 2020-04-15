package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func main() {

	testpub := "02b27de3ee5335518b06f69f4fbabb029cfc737613b100996841d5532b324a5a61"
	test, err := matched(testpub)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(test)

	testpub2 := "02b27de3ee5335518b06f69f4fbabb029cfc737613b100996841d5532b324a5a62"
	test2, err2 := matched(testpub2)
	if err2 != nil {
		log.Fatal(err2)
	}
	fmt.Println(test2)

}

func matched(pubkey string) (bool, error) {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println(dir)
	data, err := ioutil.ReadFile(filepath.Join(dir, "/assets/subatomic.json"))
	if err != nil {
		log.Printf("%s", err)
		return false, err
	}
	// fmt.Printf("%s", data)

	// dec := json.NewDecoder(strings.NewReader(string(data[:])))
	// fmt.Println(dec)

	var parsed map[string][]map[string]string
	err = json.Unmarshal([]byte(data), &parsed)
	// fmt.Println(parsed["authorized"])
	// var auth map[string][string]interface{}

	for _, v := range parsed["authorized"] {
		// fmt.Println(v)
		for _, pub := range v {
			// fmt.Println(name)
			if pub == pubkey {
				// fmt.Println(pub)
				// fmt.Println("pubkey:", pubkey)
				return true, nil
			}
		}
	}
	return false, nil
}
