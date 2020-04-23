package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"text/template"

	// "github.com/satindergrewal/subatomicgo/sagoutil"
	"subatomicgo/sagoutil"

	"github.com/satindergrewal/kmdgo"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var tpl *template.Template

// SubAtomicConfig holds the app's confugration settings
type SubAtomicConfig struct {
	Chains       []string `json:"chains"`
	SubatomicExe string   `json:"subatomic_exe"`
	SubatomicDir string   `json:"subatomic_dir"`
}

//SubAtomicConfInfo returns application's config params
func SubAtomicConfInfo() SubAtomicConfig {
	var conf SubAtomicConfig
	content, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(content, &conf)
	return conf
}

//StrToAppType converts and returns slice of string as slice of kmdgo.AppType
func StrToAppType(chain []string) []kmdgo.AppType {
	var chainskmd []kmdgo.AppType
	for _, v := range chain {
		chainskmd = append(chainskmd, kmdgo.AppType(v))
	}
	return chainskmd
}

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", idx)
	r.HandleFunc("/orderbook", orderbook).Methods("GET", "POST")
	r.HandleFunc("/orderbook/{id}", orderid).Methods("GET")
	r.HandleFunc("/orderbook/swap/{id}/{amount}/{total}", orderinit).Methods("GET")

	r.HandleFunc("/echo", echo)

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

	var conf SubAtomicConfig = SubAtomicConfInfo()

	var chains = StrToAppType(conf.Chains)

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

	type OrderPost struct {
		Base      string `json:"coin_base"`
		Rel       string `json:"coin_rel"`
		OrderList []sagoutil.OrderData
	}

	// fmt.Println("r.FormValue", r.FormValue("coin_base"))
	// fmt.Println("r.FormValue", r.FormValue("coin_rel"))

	var orderlist []sagoutil.OrderData
	orderlist = sagoutil.OrderBookList(r.FormValue("coin_base"), r.FormValue("coin_rel"), "10")

	data := OrderPost{
		Base:      r.FormValue("coin_base"),
		Rel:       r.FormValue("coin_rel"),
		OrderList: orderlist,
	}

	err := tpl.ExecuteTemplate(w, "orderbook.gohtml", data)
	if err != nil {
		http.Error(w, err.Error(), 500)
		log.Fatalln(err)
	}
}

func orderid(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	// fmt.Println(vars)
	// fmt.Println(id)

	var orderData sagoutil.OrderData
	orderData = sagoutil.OrderID(id)

	err := tpl.ExecuteTemplate(w, "orderid.gohtml", orderData)
	if err != nil {
		http.Error(w, err.Error(), 500)
		log.Fatalln(err)
	}
}

func orderinit(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]
	amount := vars["amount"]
	total := vars["total"]

	// fmt.Println(vars)
	// fmt.Println(id)
	// fmt.Println(amount)
	// fmt.Println(total)

	var orderData sagoutil.OrderData
	orderData = sagoutil.OrderID(id)

	cmdString := `./subatomic ` + orderData.Base + ` "" ` + id + ` ` + total
	fmt.Println(cmdString)

	data := struct {
		ID     string
		Amount string
		Total  string
		sagoutil.OrderData
	}{
		ID:        id,
		Amount:    amount,
		Total:     total,
		OrderData: orderData,
	}

	err := tpl.ExecuteTemplate(w, "orderinit.gohtml", data)
	if err != nil {
		http.Error(w, err.Error(), 500)
		log.Fatalln(err)
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func echo(w http.ResponseWriter, r *http.Request) {

	var conf SubAtomicConfig = SubAtomicConfInfo()
	// fmt.Println("SubatomicExe:", conf.SubatomicExe)
	// fmt.Println("SubatomicDir:", conf.SubatomicDir)

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	c.WriteMessage(1, []byte("Starting...\n"))

	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)

		err = c.WriteMessage(mt, message)

		var parsed []string
		err = json.Unmarshal([]byte(message), &parsed)
		fmt.Println("parsed", parsed)
		fmt.Println("parsed Rel:", parsed[0])
		fmt.Println("parsed ID:", parsed[1])
		fmt.Println("parsed Amount:", parsed[2])

		cmd := exec.Command(conf.SubatomicExe, parsed[0], "", parsed[1], parsed[2])
		cmd.Dir = conf.SubatomicDir
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			log.Println(err)
			fmt.Println("StdOut Nil")
			return
		}
		stderr, err := cmd.StderrPipe()
		if err != nil {
			log.Println(err)
			fmt.Println("Err Nil")
			return
		}

		if err := cmd.Start(); err != nil {
			log.Println(err)
			fmt.Println("Start")
			return
		}

		s := bufio.NewScanner(io.MultiReader(stdout, stderr))
		for s.Scan() {
			log.Printf("CMD Bytes: %s", s.Bytes())
			c.WriteMessage(1, s.Bytes())
		}

		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}

		if err := cmd.Wait(); err != nil {
			log.Println(err)
			c.WriteMessage(1, []byte(err.Error()))
			fmt.Println("Wait")
			return
		}

		c.WriteMessage(1, []byte("Finished\n"))
	}
}
