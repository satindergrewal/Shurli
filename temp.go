package main

import (
	"fmt"
	"subatomicgo/sagoutil"
)

func main() {

	var orderlist []sagoutil.OrderData
	orderlist = sagoutil.OrderBookList("KMD", "PIRATE", "10")

	fmt.Println(orderlist)

}