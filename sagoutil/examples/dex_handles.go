package main

import (
	"fmt"

	"github.com/satindergrewal/subatomicgo/sagoutil"
)

func main() {

	var handles []sagoutil.DEXHandle
	handles = sagoutil.DEXHandles()

	fmt.Println(handles)

}
