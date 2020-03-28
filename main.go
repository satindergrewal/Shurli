package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	// "github.com/satindergrewal/kmdgo"
	"kmdgo"

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
	Synced  bool
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

	var wallets []WalletInfo

	for _, v := range chains {
		// fmt.Println(i)
		// fmt.Println(v)
		appName := kmdgo.NewAppType(v)

		var info kmdgo.GetInfo

		info, err := appName.GetInfo()
		// fmt.Println(info.Error.Message)
		if err != nil {
			// fmt.Printf("Code: %v\n", info.Error.Code)
			// fmt.Printf("Message: %v\n\n", info.Error.Message)
			if info.Error.Message == "Loading block index..." {
				fmt.Println(v, "- Err happened:", info.Error.Message)
				wallets = append(wallets, WalletInfo{string(v), "Loading...", 0.0, 0, false})
			} else {
				fmt.Println(v, "- Err happened:", err)
				wallets = append(wallets, WalletInfo{string(v), "Offline", 0.0, 0, false})
			}
		} else {
			if info.Error.Message == "connection refused" {
				fmt.Println(v, "- Err happened:", info.Error.Message)
				wallets = append(wallets, WalletInfo{string(v), "Offline", 0.0, 0, false})
			} else {
				// Check status of the blockchain sync
				var tempSyncStatus bool
				if info.Result.Longestchain != info.Result.Blocks {
					tempSyncStatus = false
				} else {
					tempSyncStatus = true
				}

				wallets = append(wallets, WalletInfo{
					Ticker:  info.Result.Name,
					Status:  "Online",
					Balance: info.Result.Balance,
					Blocks:  info.Result.Longestchain,
					Synced:  tempSyncStatus,
				})
			}
		}
	}

	// fmt.Println(wallets)

	err := tpl.ExecuteTemplate(w, "index.gohtml", wallets)
	if err != nil {
		// log.Fatalf("some error")
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
