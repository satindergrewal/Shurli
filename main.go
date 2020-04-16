package main

import (
	"log"
	"net/http"
	"subatomicgo/sagoutil"
	"text/template"

	// "github.com/satindergrewal/kmdgo"
	"kmdgo"

	"github.com/gorilla/mux"
)

var tpl *template.Template

var chains = []kmdgo.AppType{"komodo", "PIRATE"}

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

// idx is a Index/Dashboard page and shows all wallet which are supported by this Subatomic Go Web App
func idx(w http.ResponseWriter, r *http.Request) {

	var wallets []sagoutil.WInfo
	wallets = sagoutil.WalletInfo(chains)

	// fmt.Println("wallets: ", wallets)

	err := tpl.ExecuteTemplate(w, "index.gohtml", wallets)
	if err != nil {
		// log.Fatalf("some error")
		http.Error(w, err.Error(), 500)
		log.Fatalln(err)
	}
}

func orderbook(w http.ResponseWriter, r *http.Request) {
	var orderlist []sagoutil.OrderData
	orderlist = sagoutil.OrderBookList("KMD", "PIRATE", "10")

	err := tpl.ExecuteTemplate(w, "orderbook.gohtml", orderlist)
	if err != nil {
		http.Error(w, err.Error(), 500)
		log.Fatalln(err)
	}
}
