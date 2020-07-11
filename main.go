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
	"os/signal"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"text/template"
	"time"

	"github.com/satindergrewal/kmdgo/kmdutil"

	"github.com/Meshbits/shurli/sagoutil"
	"github.com/satindergrewal/kmdgo"

	// "shurli/sagoutil"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

// ShurliInfo returns general information about application
// such as version and phases such as alpha, beta, stable etc.
type ShurliInfo struct {
	AppVersion string `json:"appversion"`
	AppPhase   string `json:"appphase"`
}

// DexP2pChain which shurli queries for DEXP2P API
var DexP2pChain string = "SHURLI0"

// Change to SHurli's root directory path
var rootDir string = sagoutil.ShurliRootDir()

// ShurliApp stores the information about applications
var ShurliApp = ShurliInfo{
	AppVersion: "0.0.1",
	AppPhase:   "alpha",
}

var tpl *template.Template

func check(e error) {
	if e != nil {
		panic(e)
		// log.Println(e)
	}
}

// PIDFile file stores the process ID file for shurli process
var PIDFile = "./shurli.pid"

func savePID(pid int) {

	file, err := os.Create(PIDFile)

	if err != nil {
		log.Printf("Unable to create pid file : %v\n", err)
		os.Exit(1)
	}

	defer file.Close()

	_, err = file.WriteString(strconv.Itoa(pid))

	if err != nil {
		log.Printf("Unable to create pid file : %v\n", err)
		os.Exit(1)
	}

	file.Sync() // flush to disk

}

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
}

func main() {

	if len(os.Args) != 2 {
		fmt.Printf("Usage : %s [start|stop] \n ", os.Args[0]) // return the program name back to %s
		os.Exit(0)                                            // graceful exit
	}

	// If running with command "./shurli main"
	// this condition will trigger the code to run without exiting the stdout.
	// User has to press CTRL or CMD + C to interrup the process
	if strings.ToLower(os.Args[1]) == "main" {

		// Make arrangement to remove PID file upon receiving the SIGTERM from kill command
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, os.Interrupt, os.Kill, syscall.SIGTERM)

		go func() {
			signalType := <-ch
			signal.Stop(ch)
			fmt.Println("Exit command received. Exiting...")

			// this is a good place to flush everything to disk
			// before terminating.
			fmt.Println("Received signal type : ", signalType)

			// remove PID file
			os.Remove(PIDFile)

			os.Exit(0)

		}()

		// Insert blank lines before starting next log
		sagoutil.Log.Printf("\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n")
		// Display Shurli Application's version and phase
		sagoutil.Log.Printf(">>> Shurli version: %s %s\n", ShurliApp.AppVersion, ShurliApp.AppPhase)
		// shurli mux code start here
		sagoutil.ShurliStartMsg()

		// Setup/Define http (Gorilla) Mux
		r := mux.NewRouter()
		r.HandleFunc("/", idx)
		r.HandleFunc("/orderbook", orderbook).Methods("GET", "POST")
		r.HandleFunc("/orderbook/{id}", orderid).Methods("GET")
		r.HandleFunc("/orderbook/swap/{id}/{amount}/{total}", orderinit).Methods("GET")
		r.HandleFunc("/history", swaphistory)

		// Gorilla WebSockets echo example used to do give subatomic trade data updates to orderinit
		r.HandleFunc("/echo", echo)

		// favicon.ico file
		r.HandleFunc("/favicon.ico", faviconHandler)

		// public assets files
		r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./public/"))))
		sagoutil.Log.Fatal(http.ListenAndServe(":8080", r))
	}

	// using command "./shurli start" will show the daemon process info and exit stdout to terminal
	// leaving "shurli" process running the background, and storing it's process ID info in "shurli.pid" file
	if strings.ToLower(os.Args[1]) == "start" {

		// // Display Shurli Application's version and phase
		// sagoutil.Log.Printf(">>> Shurli version: %s %s\n", ShurliApp.AppVersion, ShurliApp.AppPhase)

		// check if daemon already running.
		if _, err := os.Stat(PIDFile); err == nil {
			fmt.Println("Already running or ./shirli.pid file exist.")
			os.Exit(1)
		}

		// fmt.Println(os.Args)

		cmd := exec.Command(os.Args[0], "main")
		cmd.Start()
		fmt.Println("Shurli started as daemon. Process ID is : ", cmd.Process.Pid)
		sagoutil.Log.Println("Shurli started as daemon. Process ID is : ", cmd.Process.Pid)
		savePID(cmd.Process.Pid)
		// os.Exit(0)

		// Check if "DEX" blockchain is already running on system.
		// If "komodo.pid" is present in "DEX" data directory, it means
		// - "DEX" blockchain is already running
		// - or the previous process did not delete the "komodo.pid" file before exiting due to some reason, i.e. daemon crash etc.
		// 		- In this case, just delete the "komodo.pid" file and next time "shurli" should be able to start "DEX" blockchain.
		appName := DexP2pChain
		dir := kmdutil.AppDataDir(appName, false)
		// fmt.Println(dir)
		// If "DEX" blockchain is running already, print notification
		if _, err := os.Stat(dir + "/komodod.pid"); err == nil {
			fmt.Println("[Shurli] " + DexP2pChain + " blockchain already running or " + DexP2pChain + " pid file exist.")
			sagoutil.Log.Println("[Shurli] " + DexP2pChain + " blockchain already running or " + DexP2pChain + " pid file exist.")
			os.Exit(1)
		} else {
			// If "DEX" blockchain isn't found running already, start it in daemon mode.
			var conf sagoutil.SubAtomicConfig = sagoutil.SubAtomicConfInfo()
			// fmt.Prinln("DexNSPV: " conf.DexNSPV)
			// fmt.Prinln("DexAddnode: " conf.DexAddnode)
			// fmt.Prinln("DexPubkey: " conf.DexPubkey)
			// fmt.Prinln("DexHandle: " conf.DexHandle)
			// fmt.Prinln("DexRecvzaddr: " conf.DexRecvzaddr)
			// fmt.Prinln("DexRecvtaddr: " conf.DexRecvtaddr)

			// dexnspv := "-nSPV=" + conf.DexNSPV
			// if conf.DexNSPV == "0" {
			// 	dexnspv = ""
			// }
			acname := "-ac_name=" + DexP2pChain
			dexaddnode := "-addnode=" + conf.DexAddnode
			dexpubkey := "-pubkey=" + conf.DexPubkey
			dexhandle := "-handle=" + conf.DexHandle
			dexrecvzaddr := "-recvZaddr=" + conf.DexRecvZAddr
			dexrecvtaddr := "-recvTaddr=" + conf.DexRecvTAddr
			dexcmd := exec.Command("assets/komodod", acname, "-daemon", "-server", "-ac_supply=10", "-dexp2p=2", dexaddnode, dexpubkey, dexhandle, dexrecvzaddr, dexrecvtaddr)
			if runtime.GOOS == "windows" {
				dexcmd = exec.Command("assets/komodod.exe", acname, "-daemon", "-server", "-ac_supply=10", "-dexp2p=2", dexaddnode, dexpubkey, dexhandle, dexrecvzaddr, dexrecvtaddr)
			}
			// fmt.Println(conf.SubatomicDir)
			// dexcmd.Dir = conf.SubatomicDir
			// out, err := dexcmd.Output()
			// if err != nil {
			// 	log.Fatalf("dexcmd.Start() failed with %s\n", err)
			// } else {
			// 	fmt.Printf("%s", out)
			// }
			err := dexcmd.Start()
			if err != nil {
				log.Fatalf("dexcmd.Start() failed with %s\n", err)
			}
			fmt.Println("[Shurli] Started "+DexP2pChain+" komodod. Process ID is : ", dexcmd.Process.Pid)
			fmt.Println("[Shurli] " + DexP2pChain + " chain params: ")
			// fmt.Println("\t" + DexP2pChain + " nSPV: ", conf.DexNSPV)
			// sagoutil.Log.Println("\t" + DexP2pChain + " nSPV: ", conf.DexNSPV)
			fmt.Println("\t"+DexP2pChain+" addnode: ", conf.DexAddnode)
			sagoutil.Log.Println("\t"+DexP2pChain+" addnode: ", conf.DexAddnode)
			fmt.Println("\t"+DexP2pChain+" pubkey: ", conf.DexPubkey)
			sagoutil.Log.Println("\t"+DexP2pChain+" pubkey: ", conf.DexPubkey)
			fmt.Println("\t"+DexP2pChain+" handle: ", conf.DexHandle)
			sagoutil.Log.Println("\t"+DexP2pChain+" handle: ", conf.DexHandle)
			fmt.Println("\t"+DexP2pChain+" recvZaddr: ", conf.DexRecvZAddr)
			sagoutil.Log.Println("\t"+DexP2pChain+" recvZaddr: ", conf.DexRecvZAddr)
			fmt.Println("\t"+DexP2pChain+" recvTaddr: ", conf.DexRecvTAddr)
			sagoutil.Log.Println("\t"+DexP2pChain+" recvTaddr: ", conf.DexRecvTAddr)
			sagoutil.Log.Println("[Shurli] Started DEX komodod. Process ID is : ", dexcmd.Process.Pid)
			os.Exit(0)
		}
	}

	// upon receiving the stop command
	// read the Process ID stored in PIDfile
	// kill the process using the Process ID
	// and exit. If Process ID does not exist, prompt error and quit

	if strings.ToLower(os.Args[1]) == "stop" {

		appName := kmdgo.NewAppType(kmdgo.AppType(DexP2pChain))
		var info kmdgo.Stop
		info, err := appName.Stop()
		if err != nil {
			fmt.Printf("Code: %v\n", info.Error.Code)
			fmt.Printf("Message: %v\n\n", info.Error.Message)
			log.Println("Err happened", err)
		}
		// fmt.Println(info)
		fmt.Println("[Shurli] ", info.Result)
		sagoutil.Log.Println("[Shurli] ", info.Result)

		if _, err := os.Stat(PIDFile); err == nil {
			data, err := ioutil.ReadFile(PIDFile)
			if err != nil {
				fmt.Println("Shurli is not running.")
				os.Exit(1)
			}
			ProcessID, err := strconv.Atoi(string(data))

			if err != nil {
				fmt.Println("[Shurli] Unable to read and parse process id found in ", PIDFile)
				os.Exit(1)
			}

			process, err := os.FindProcess(ProcessID)

			if err != nil {
				fmt.Printf("[Shurli] Unable to find process ID [%v] with error %v \n", ProcessID, err)
				os.Exit(1)
			}
			// remove PID file
			os.Remove(PIDFile)

			fmt.Printf("Stopping Shurli daemon... Killing process ID [%v] now.\n", ProcessID)
			sagoutil.Log.Printf("Stopping Shurli daemon... Killing process ID [%v] now.\n", ProcessID)
			// kill process and exit immediately
			err = process.Kill()

			if err != nil {
				fmt.Printf("[Shurli] Unable to kill process ID [%v] with error %v \n", ProcessID, err)
				sagoutil.Log.Printf("[Shurli] Unable to kill process ID [%v] with error %v \n", ProcessID, err)
				os.Exit(1)
			} else {
				fmt.Printf("[Shurli stopped] Killed process ID [%v]\n", ProcessID)
				sagoutil.Log.Printf("[Shurli stopped] Killed process ID [%v]\n", ProcessID)
				os.Exit(0)
			}

		} else {

			fmt.Println("Shurli is not running.")
			os.Exit(1)
		}
	} else {
		fmt.Printf("[Shurli] Unknown command : %v\n", os.Args[1])
		fmt.Printf("Shurli Usage : %s [start|stop]\n", os.Args[0]) // return the program name back to %s
		os.Exit(1)
	}
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "favicon.png")
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
		sagoutil.Log.Fatalln(err)
	}
}

func orderbook(w http.ResponseWriter, r *http.Request) {

	// Change to SHurli's root directory path
	os.Chdir(rootDir)

	type OrderPost struct {
		Base      string `json:"coin_base"`
		Rel       string `json:"coin_rel"`
		Results   string `json:"results"`
		SortBy    string `json:"sortby"`
		BaseBal   float64
		RelBal    float64
		BaseIcon  string
		RelIcon   string
		OrderList []sagoutil.OrderData
	}

	// fmt.Println("r.FormValue", r.FormValue("coin_base"))
	// fmt.Println("r.FormValue", r.FormValue("coin_rel"))
	// fmt.Println("r.FormValue", r.FormValue("result_limit"))
	// fmt.Println("r.FormValue", r.FormValue("sortby"))

	var orderlist []sagoutil.OrderData
	orderlist = sagoutil.OrderBookList(r.FormValue("coin_base"), r.FormValue("coin_rel"), r.FormValue("result_limit"), r.FormValue("sortby"))

	var baseRelWallet = []kmdgo.AppType{kmdgo.AppType(r.FormValue("coin_base")), kmdgo.AppType(r.FormValue("coin_rel"))}

	var wallets []sagoutil.WInfo
	wallets = sagoutil.WalletInfo(baseRelWallet)
	// fmt.Println(wallets[0].Balance)
	// fmt.Println(wallets[0].ZBalance)
	// fmt.Println(wallets[1].Balance)
	// fmt.Println(wallets[1].ZBalance)

	var relBalance, baseBalance float64
	if strings.HasPrefix(r.FormValue("coin_base"), "z") {
		baseBalance = wallets[0].ZBalance
	} else if strings.HasPrefix(r.FormValue("coin_base"), "PIRATE") {
		baseBalance = wallets[0].ZBalance
	} else {
		baseBalance = wallets[0].Balance
	}

	if strings.HasPrefix(r.FormValue("coin_rel"), "z") {
		relBalance = wallets[1].ZBalance
	} else if strings.HasPrefix(r.FormValue("coin_rel"), "PIRATE") {
		relBalance = wallets[1].ZBalance
	} else {
		relBalance = wallets[1].Balance
	}

	data := OrderPost{
		Base:      r.FormValue("coin_base"),
		Rel:       r.FormValue("coin_rel"),
		Results:   r.FormValue("result_limit"),
		SortBy:    r.FormValue("sortby"),
		BaseBal:   baseBalance,
		RelBal:    relBalance,
		BaseIcon:  wallets[0].Icon,
		RelIcon:   wallets[1].Icon,
		OrderList: orderlist,
	}

	err := tpl.ExecuteTemplate(w, "orderbook.gohtml", data)
	if err != nil {
		http.Error(w, err.Error(), 500)
		sagoutil.Log.Fatalln(err)
	}
}

func orderid(w http.ResponseWriter, r *http.Request) {

	// Change to SHurli's root directory path
	os.Chdir(rootDir)

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
		sagoutil.Log.Fatalln(err)
	}
}

func orderinit(w http.ResponseWriter, r *http.Request) {

	// Change to SHurli's root directory path
	os.Chdir(rootDir)

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

	orderDataJSON, _ := json.Marshal(orderData)
	sagoutil.Log.Println("orderData JSON:", string(orderDataJSON))

	cmdString := `[subatomic] ./subatomic ` + orderData.Base + ` "" ` + id + ` ` + total
	sagoutil.Log.Println(cmdString)
	log.Println(cmdString)

	var conf sagoutil.SubAtomicConfig = sagoutil.SubAtomicConfInfo()

	data := struct {
		ID           string
		Amount       string
		Total        string
		BaseExplorer string
		RelExplorer  string
		sagoutil.OrderData
		OrderDataJson string
	}{
		ID:            id,
		Amount:        amount,
		Total:         total,
		OrderData:     orderData,
		OrderDataJson: string(orderDataJSON),
		BaseExplorer:  conf.Explorers[strings.ReplaceAll(orderData.Base, "z", "")],
		RelExplorer:   conf.Explorers[strings.ReplaceAll(orderData.Rel, "z", "")],
	}

	err := tpl.ExecuteTemplate(w, "orderinit.gohtml", data)
	if err != nil {
		http.Error(w, err.Error(), 500)
		sagoutil.Log.Fatalln(err)
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func echo(w http.ResponseWriter, r *http.Request) {

	var conf sagoutil.SubAtomicConfig = sagoutil.SubAtomicConfInfo()
	sagoutil.Log.Println("SubatomicExe:", conf.SubatomicExe)
	sagoutil.Log.Println("SubatomicDir:", conf.SubatomicDir)

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	c.WriteMessage(1, []byte(`{"state":"Starting..."}`))

	exPath := filepath.Join(rootDir, "assets")
	os.Chdir(exPath)
	sagoutil.Log.Println(exPath)

	var filename string
	newLine := "\n"

	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			sagoutil.Log.Println("read:", err)
			break
		}
		sagoutil.Log.Printf("recv: %s", message)

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
			sagoutil.Log.Println(txidMsg)

			err = c.WriteMessage(1, []byte(txidMsg))
		}

		var parsed []string
		err = json.Unmarshal([]byte(message), &parsed)
		sagoutil.Log.Println("parsed", parsed)

		if len(parsed) > 0 && parsed[0] == "subatomic_cmd" {
			// fmt.Println("parsed Rel:", parsed[0])
			// fmt.Println("parsed ID:", parsed[1])
			// fmt.Println("parsed Amount:", parsed[2])

			// Create a new context and add a timeout to it
			ctx, cancel := context.WithTimeout(context.Background(), 180*time.Second)
			defer cancel() // The cancel should be deferred so resources are cleaned up

			// cmd := exec.Command(conf.SubatomicExe, parsed[0], "", parsed[1], parsed[2])
			// Create the command with our context
			cmd := exec.CommandContext(ctx, "./subatomic", parsed[1], "", parsed[2], parsed[3])
			if runtime.GOOS == "windows" {
				cmd = exec.CommandContext(ctx, "./subatomic.exe", parsed[1], "", parsed[2], parsed[3])
			}
			// cmd.Dir = conf.SubatomicDir

			stdout, err := cmd.StdoutPipe()
			if err != nil {
				sagoutil.Log.Println(err)
				sagoutil.Log.Println("StdOut Nil")
				return
			}
			stderr, err := cmd.StderrPipe()
			if err != nil {
				sagoutil.Log.Println(err)
				sagoutil.Log.Println("Err Nil")
				return
			}

			if err := cmd.Start(); err != nil {
				sagoutil.Log.Println(err)
				sagoutil.Log.Println("Start")
				return
			}

			// We want to check the context error to see if the timeout was executed.
			// The error returned by cmd.Output() will be OS specific based on what
			// happens when a process is killed.
			if ctx.Err() == context.DeadlineExceeded {
				sagoutil.Log.Println("Command timed out")
				return
			}

			s := bufio.NewScanner(io.MultiReader(stdout, stderr))

			newpath := filepath.Join(rootDir, "swaplogs")
			err = os.MkdirAll(newpath, 0755)
			check(err)

			currentUnixTimestamp := int32(time.Now().Unix())
			filename = rootDir + "/swaplogs/" + sagoutil.IntToString(currentUnixTimestamp) + "_" + parsed[2] + ".log"
			fmt.Println(filename)
			// fmt.Println(String(currentUnixTimestamp))

			// If the file doesn't exist, create it, or append to the file
			f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			check(err)
			defer f.Close()

			w := bufio.NewWriter(f)

			for s.Scan() {
				sagoutil.Log.Printf("[subatomic] CMD Bytes: %s", s.Bytes())
				log.Printf("[subatomic] CMD Bytes: %s", s.Bytes())
				// c.WriteMessage(1, s.Bytes())

				logstr, err := sagoutil.SwapLogFilter(string(s.Bytes()), "single")
				if err != nil {
					// fmt.Println(err)
				} else {
					// fmt.Println(logstr)
					c.WriteMessage(1, []byte(logstr))
				}

				l := s.Bytes()
				l = append(l, newLine...)
				_, err = w.Write(l)
				check(err)
				// fmt.Printf("wrote %d bytes\n", n4)
			}

			m := message
			m = append(m, newLine...)
			_, err = w.Write(m)
			check(err)

			w.Flush()

			// err = c.WriteMessage(mt, message)
			// if err != nil {
			// 	log.Println("write:", err)
			// 	break
			// }

			if err := cmd.Wait(); err != nil {
				sagoutil.Log.Println("[subatomic]", err)
				log.Println("[subatomic]", err)
				c.WriteMessage(1, []byte(`{"state": "`+err.Error()+`"}`))
				sagoutil.Log.Println("[subatomic] Wait")
				log.Println("[subatomic] Wait")
				return
			}
		}

		// fmt.Println("filename", filename)

		os.Chdir(rootDir)

		c.WriteMessage(1, []byte(`{"state":"Finished"}`))
	}
}

func swaphistory(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Content-Type", "application/json")
	var history sagoutil.SwapsHistory
	allhistory, err := history.SwapsHistory()
	// fmt.Println(allhistory)

	// if err != nil {
	// 	json.NewEncoder(w).Encode(err.Error())
	// } else {
	// 	json.NewEncoder(w).Encode(allhistory)
	// }

	err = tpl.ExecuteTemplate(w, "swaphistory.gohtml", allhistory)
	if err != nil {
		// log.Fatalf("some error")
		http.Error(w, err.Error(), 500)
		sagoutil.Log.Fatalln(err)
	}
}
