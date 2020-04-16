package main

import (
	"fmt"
	"subatomicgo/sagoutil"
)

func main() {

	var handles []sagoutil.DEXHandle
	handles = sagoutil.DEXHandles()

	fmt.Println(handles)
	fmt.Println(handles[0].Handle)
}
