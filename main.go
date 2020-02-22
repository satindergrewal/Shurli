package main

import (
	"fmt"
	"kmdgo"
	"log"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
)

var tpl *template.Template

var chains = []kmdgo.AppType{"komodo", "PIRATE"}

// WalletInfo stores data to display on Wallet info screen
type WalletInfo struct {
	Ticker  string
	Status  string
	Balance float64
	Blocks  int
	Synched bool
}

// Wallets store array of WalletInfo data type
type Wallets struct {
	Wallet []WalletInfo
}

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", idx)
	r.HandleFunc("/orderbook", orderbook)

	// favicon.ico file
	r.HandleFunc("/favicon.ico", faviconHandler)

	// public assets files
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./public/"))))
	log.Fatal(http.ListenAndServe(":8080", r))
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "favicon.ico")
}

func idx(w http.ResponseWriter, r *http.Request) {

	var wallets Wallets

	for i, v := range chains {
		// fmt.Println(i)
		// fmt.Println(v)
		appName := kmdgo.NewAppType(v)

		var info kmdgo.GetInfo

		info, err := appName.GetInfo()
		if err != nil {
			fmt.Printf("Code: %v\n", info.Error.Code)
			fmt.Printf("Message: %v\n\n", info.Error.Message)
			log.Fatalln("Err happened", err)
			wallets.Wallet[i].Status = "Offline"
		} else {
			wallets.Wallet[i].Ticker = info.Result.Name
			wallets.Wallet[i].Status = "Online"
			wallets.Wallet[i].Balance = info.Result.Balance
			wallets.Wallet[i].Blocks = info.Result.Longestchain
			if info.Result.Longestchain != info.Result.Blocks {
				wallets.Wallet[i].Synched = false
			} else {
				wallets.Wallet[i].Synched = true
			}
		}
	}

	fmt.Println(wallets)

	err := tpl.ExecuteTemplate(w, "index.gohtml", nil)
	if err != nil {
		http.Error(w, err.Error(), 500)
		log.Fatalln(err)
	}
}

func orderbook(w http.ResponseWriter, r *http.Request) {
	err := tpl.ExecuteTemplate(w, "orderbook.gohtml", nil)
	if err != nil {
		http.Error(w, err.Error(), 500)
		log.Fatalln(err)
	}
}
