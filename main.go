package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/satindergrewal/shurli/sagoutil"
	// "shurli/sagoutil"

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

	var conf sagoutil.SubAtomicConfig = sagoutil.SubAtomicConfInfo()

	var chains = sagoutil.StrToAppType(conf.Chains)

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
		Results   string `json:"results"`
		SortBy    string `json:"sortby"`
		OrderList []sagoutil.OrderData
	}

	// fmt.Println("r.FormValue", r.FormValue("coin_base"))
	// fmt.Println("r.FormValue", r.FormValue("coin_rel"))
	// fmt.Println("r.FormValue", r.FormValue("result_limit"))
	// fmt.Println("r.FormValue", r.FormValue("sortby"))

	var orderlist []sagoutil.OrderData
	orderlist = sagoutil.OrderBookList(r.FormValue("coin_base"), r.FormValue("coin_rel"), r.FormValue("result_limit"), r.FormValue("sortby"))

	data := OrderPost{
		Base:      r.FormValue("coin_base"),
		Rel:       r.FormValue("coin_rel"),
		Results:   r.FormValue("result_limit"),
		SortBy:    r.FormValue("sortby"),
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
	// fmt.Println(orderData)

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

	var conf sagoutil.SubAtomicConfig = sagoutil.SubAtomicConfInfo()

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
		BaseExplorer: conf.Explorers[strings.ReplaceAll(orderData.Base, "z", "")],
		RelExplorer:  conf.Explorers[strings.ReplaceAll(orderData.Rel, "z", "")],
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

	var conf sagoutil.SubAtomicConfig = sagoutil.SubAtomicConfInfo()
	// fmt.Println("SubatomicExe:", conf.SubatomicExe)
	// fmt.Println("SubatomicDir:", conf.SubatomicDir)

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	c.WriteMessage(1, []byte(`{"state":"Starting..."}`))

	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)

		err = c.WriteMessage(mt, message)

		type opIdMsg struct {
			Opid string `json:"opid"`
			Coin string `json:"coin"`
		}
		var opidmsg opIdMsg
		err = json.Unmarshal([]byte(message), &opidmsg)
		// fmt.Println(opidmsg.Opid)
		// fmt.Println(opidmsg.Coin)

		if len(opidmsg.Opid) > 0 {
			txidMsg, _ := sagoutil.TxIDFromOpID(opidmsg.Coin, opidmsg.Opid)
			fmt.Println(txidMsg)

			err = c.WriteMessage(1, []byte(txidMsg))
		}

		var parsed []string
		err = json.Unmarshal([]byte(message), &parsed)
		// fmt.Println("parsed", parsed)

		if len(parsed) > 0 {
			// fmt.Println("parsed Rel:", parsed[0])
			// fmt.Println("parsed ID:", parsed[1])
			// fmt.Println("parsed Amount:", parsed[2])

			// Create a new context and add a timeout to it
			ctx, cancel := context.WithTimeout(context.Background(), 180*time.Second)
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
			filename := "./swaplogs/" + sagoutil.IntToString(currentUnixTimestamp) + "_" + parsed[1] + ".log"
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

				logstr, err := sagoutil.SwapLogFilter(string(s.Bytes()), "single")
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
		}

		c.WriteMessage(1, []byte(`{"state":"Finished"}`))
	}
}
