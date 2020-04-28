package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
	"time"

	"github.com/satindergrewal/subatomicgo/sagoutil"
	// "subatomicgo/sagoutil"

	"github.com/satindergrewal/kmdgo"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var tpl *template.Template

func check(e error) {
	if e != nil {
		panic(e)
		// log.Println(e)
	}
}

func String(n int32) string {
	buf := [11]byte{}
	pos := len(buf)
	i := int64(n)
	signed := i < 0
	if signed {
		i = -i
	}
	for {
		pos--
		buf[pos], i = '0'+byte(i%10), i/10
		if i == 0 {
			if signed {
				pos--
				buf[pos] = '-'
			}
			return string(buf[pos:])
		}
	}
}

// SubAtomicConfig holds the app's confugration settings
type SubAtomicConfig struct {
	Chains       []string          `json:"chains"`
	SubatomicExe string            `json:"subatomic_exe"`
	SubatomicDir string            `json:"subatomic_dir"`
	Explorers    map[string]string `json:"explorers"`
}

//SubAtomicConfInfo returns application's config params
func SubAtomicConfInfo() SubAtomicConfig {
	var conf SubAtomicConfig
	content, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(content, &conf)
	// fmt.Println(conf.Explorers["KMD"])
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

	var conf SubAtomicConfig = SubAtomicConfInfo()

	data := struct {
		ID           string
		Amount       string
		Total        string
		BaseExplorer string
		RelExplorer  string
		sagoutil.OrderData
	}{
		ID:           id,
		Amount:       amount,
		Total:        total,
		OrderData:    orderData,
		BaseExplorer: conf.Explorers[orderData.Base],
		RelExplorer:  conf.Explorers[orderData.Rel],
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

	c.WriteMessage(1, []byte(`{"state":"Starting...}"`))

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)

		// err = c.WriteMessage(mt, message)

		var parsed []string
		err = json.Unmarshal([]byte(message), &parsed)
		// fmt.Println("parsed", parsed)
		// fmt.Println("parsed Rel:", parsed[0])
		// fmt.Println("parsed ID:", parsed[1])
		// fmt.Println("parsed Amount:", parsed[2])

		// Create a new context and add a timeout to it
		ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
		defer cancel() // The cancel should be deferred so resources are cleaned up

		// cmd := exec.Command(conf.SubatomicExe, parsed[0], "", parsed[1], parsed[2])
		// Create the command with our context
		cmd := exec.CommandContext(ctx, conf.SubatomicExe, parsed[0], "", parsed[1], parsed[2])
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

		// We want to check the context error to see if the timeout was executed.
		// The error returned by cmd.Output() will be OS specific based on what
		// happens when a process is killed.
		if ctx.Err() == context.DeadlineExceeded {
			fmt.Println("Command timed out")
			return
		}

		s := bufio.NewScanner(io.MultiReader(stdout, stderr))

		newpath := filepath.Join(".", "swaplogs")
		err = os.MkdirAll(newpath, 0755)
		check(err)

		currentUnixTimestamp := int32(time.Now().Unix())
		filename := "./swaplogs/" + String(currentUnixTimestamp) + "_" + parsed[1] + ".log"
		fmt.Println(filename)
		// fmt.Println(String(currentUnixTimestamp))

		// If the file doesn't exist, create it, or append to the file
		f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		check(err)
		defer f.Close()

		w := bufio.NewWriter(f)

		for s.Scan() {
			log.Printf("CMD Bytes: %s", s.Bytes())
			// c.WriteMessage(1, s.Bytes())

			logstr, err := sagoutil.SwapLogFilter(string(s.Bytes()))
			if err != nil {
				// fmt.Println(err)
			} else {
				// fmt.Println(logstr)
				c.WriteMessage(1, []byte(logstr))
			}

			l := s.Bytes()
			newLine := "\n"
			l = append(l, newLine...)
			_, err = w.Write(l)
			check(err)
			// fmt.Printf("wrote %d bytes\n", n4)
		}

		w.Flush()

		// err = c.WriteMessage(mt, message)
		// if err != nil {
		// 	log.Println("write:", err)
		// 	break
		// }

		if err := cmd.Wait(); err != nil {
			log.Println(err)
			c.WriteMessage(1, []byte(`{"state": "`+err.Error()+`"}`))
			fmt.Println("Wait")
			return
		}

		c.WriteMessage(1, []byte(`{"state":"Finished"}`))
	}
}
