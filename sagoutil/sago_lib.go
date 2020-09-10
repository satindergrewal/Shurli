package sagoutil

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/satindergrewal/kmdgo"
)

// DexP2pChain which shurli queries for DEXP2P API
var DexP2pChain string = "SHURLI0"

// ShurliRootDir returns Shurli root directory path
func ShurliRootDir() string {
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	Log.Println(currentDir)
	return currentDir
}

// WInfo type stores data to display on Wallet info screen
type WInfo struct {
	Name       string
	Ticker     string
	Icon       string
	Status     string
	Balance    float64
	ZBalance   float64
	Blocks     int
	Synced     bool
	Shielded   bool
	TValidAddr bool
	ZValidAddr bool
}

// WalletInfo method returns processed data to display on Dashboard
func WalletInfo(chains []kmdgo.AppType) []WInfo {
	var wallets []WInfo

	// fmt.Println(chains)

	stats, err := kmdgo.NewAppType(kmdgo.AppType(DexP2pChain)).DEXStats()
	if err != nil {
		Log.Printf("Code: %v\n", stats.Error.Code)
		Log.Printf("Message: %v\n\n", stats.Error.Message)
		// log.Fatalln("Err happened", err)
	}

	// fmt.Println("stats value", stats)
	// fmt.Println("Recvaddr", stats.Result.Recvaddr)
	// fmt.Println("RecvZaddr", stats.Result.RecvZaddr)

	for _, v := range chains {
		// fmt.Println(i)
		// fmt.Println("v", v)

		switch v {
		case "KMD":
			v = "komodo"
		case "Pirate":
			v = "PIRATE"
		}

		vWithoutZ := strings.ReplaceAll(string(v), "z", "")

		coinConfInfo := GetCoinConfInfo(vWithoutZ)
		// fmt.Println(strings.ToLower(coinConfInfo.Ticker))
		tmpicon := strings.ToLower(coinConfInfo.Ticker)
		// fmt.Println(coinConfInfo)
		// if !!coinConfInfo.Shielded {
		// 	fmt.Println(coinConfInfo.Shielded)
		// 	fmt.Println(!coinConfInfo.Shielded)
		// }

		appName := kmdgo.NewAppType(kmdgo.AppType(vWithoutZ))

		var info kmdgo.GetInfo

		info, err := appName.GetInfo()
		// fmt.Println(info.Error.Message)
		if err != nil {
			// fmt.Printf("Code: %v\n", info.Error.Code)
			// fmt.Printf("Message: %v\n\n", info.Error.Message)
			if info.Error.Message == "Loading block index..." {
				Log.Println(v, "- Err happened:", info.Error.Message)
				wallets = append(wallets, WInfo{coinConfInfo.Name, coinConfInfo.Ticker, tmpicon, "Loading...", 0.0, 0, 0, false, false, false, false})
			} else if info.Error.Message == "Rescanning..." {
				Log.Println(v, "- Err happened:", info.Error.Message)
				wallets = append(wallets, WInfo{coinConfInfo.Name, coinConfInfo.Ticker, tmpicon, "Rescanning...", 0.0, 0, 0, false, false, false, false})
			} else if info.Error.Message == "Rewinding blocks if needed..." {
				Log.Println(v, "- Err happened:", info.Error.Message)
				wallets = append(wallets, WInfo{coinConfInfo.Name, coinConfInfo.Ticker, tmpicon, "Rewinding blocks if needed...", 0.0, 0, 0, false, false, false, false})
			} else {
				Log.Println(v, "- Err happened:", err)
				wallets = append(wallets, WInfo{coinConfInfo.Name, coinConfInfo.Ticker, tmpicon, "Offline", 0.0, 0, 0, false, false, false, false})
			}
		} else {
			if info.Error.Message == "connection refused" {
				Log.Println(v, "- Err happened:", info.Error.Message)
				wallets = append(wallets, WInfo{coinConfInfo.Name, coinConfInfo.Ticker, tmpicon, "Offline", 0.0, 0, 0, false, false, false, false})
			} else {

				// Check status of the blockchain sync
				var tempSyncStatus bool
				// GetBlockChainInfo to get value of "verificationprogress" to check if blockchain synced or not
				var gb kmdgo.GetBlockchainInfo
				gb, err := appName.GetBlockchainInfo()
				if err != nil {
					Log.Printf("Code: %v\n", gb.Error.Code)
					Log.Printf("Message: %v\n\n", gb.Error.Message)
					Log.Println("Err happened", err)
				}

				// Validate Transaparent Address
				var vldadr kmdgo.ValidateAddress
				taddress := stats.Result.Recvaddr
				vldadr, err = appName.ValidateAddress(taddress)
				if err != nil {
					Log.Printf("Code: %v\n", vldadr.Error.Code)
					Log.Printf("Message: %v\n\n", vldadr.Error.Message)
					Log.Println("Err happened", err)
				}

				// if math.Round(gb.Result.Verificationprogress) != 1 {
				if gb.Result.Verificationprogress >= 0.9999995 {
					tempSyncStatus = true
				} else {
					tempSyncStatus = false
				}

				// if info.Result.Longestchain != info.Result.Blocks {
				// 	tempSyncStatus = false
				// } else {
				// 	tempSyncStatus = true
				// }

				if !!coinConfInfo.Shielded {
					// fmt.Printf("it is %s, Getting it's Z balance...\n", coinConfInfo.Name)
					var zblc kmdgo.ZGetBalance

					args := make(kmdgo.APIParams, 2)
					args[0] = stats.Result.RecvZaddr
					//args[1] = 1
					// fmt.Println(args)

					// Validate if Shielded adddress is = ismine
					var zvldadr kmdgo.ZValidateAddress
					zaddress := stats.Result.RecvZaddr
					zvldadr, err := appName.ZValidateAddress(zaddress)
					if err != nil {
						Log.Printf("Code: %v\n", zvldadr.Error.Code)
						Log.Printf("Message: %v\n\n", zvldadr.Error.Message)
						Log.Println("Err happened", err)
					}

					// Get balance of Shielded address
					zblc, err = appName.ZGetBalance(args)
					if err != nil {
						Log.Printf("Code: %v\n", zblc.Error.Code)
						Log.Printf("Message: %v\n\n", zblc.Error.Message)
						// log.Fatalln("Err happened", err)
					}

					// fmt.Println("zblc value", zblc)
					// fmt.Println("-------")
					// fmt.Printf("\n%0.8f\n", zblc.Result)

					wallets = append(wallets, WInfo{
						Name:       coinConfInfo.Name,
						Ticker:     coinConfInfo.Ticker,
						Icon:       strings.ToLower(coinConfInfo.Ticker),
						Status:     "Online",
						ZBalance:   zblc.Result,
						Balance:    info.Result.Balance,
						Blocks:     info.Result.Blocks,
						Synced:     tempSyncStatus,
						Shielded:   coinConfInfo.Shielded,
						TValidAddr: vldadr.Result.Ismine,
						ZValidAddr: zvldadr.Result.Ismine,
					})

				} else {
					wallets = append(wallets, WInfo{
						Name:       coinConfInfo.Name,
						Ticker:     coinConfInfo.Ticker,
						Icon:       strings.ToLower(coinConfInfo.Ticker),
						Status:     "Online",
						Balance:    info.Result.Balance,
						Blocks:     info.Result.Blocks,
						Synced:     tempSyncStatus,
						Shielded:   coinConfInfo.Shielded,
						TValidAddr: vldadr.Result.Ismine,
					})
				}

			}
		}
	}

	// fmt.Println(wallets)
	return wallets
}

// DEXHandle stores data to get authorised/unauthorised hanldes from all broadcasting traders
type DEXHandle struct {
	Pubkey     string
	Handle     string
	DEXPubkey  string
	Authorised bool
}

// DEXHandles returns public address's public Key, DEX CC specific public key, and handle data set picked from:
// - DEX_list handles API
// - subatomic.json file taken from the source code of komodo code
// The purpose of this function is to return the data about authorised and unauthorised handles and show that in orderbook or orderlist in the GUI application.
func DEXHandles() []DEXHandle {
	var handles []DEXHandle

	// handles = append(handles, DEXHandle{"03732f8ef851ff234c74d0df575c2c5b159e2bab3faca4ec52b3f217d5cda5361d", "satinder", "01b5d5b1991152fd45e4ba7005a5a752c2018634a9a6cdeb06b633e731e7b5f46b"})

	appName := kmdgo.AppType(DexP2pChain)

	var list kmdgo.DEXList

	args := make(kmdgo.APIParams, 3)
	// stopat
	args[0] = "0"
	// minpriority
	args[1] = "0"
	// tagA
	args[2] = "handles"
	// fmt.Println(args)

	list, err := appName.DEXList(args)
	if err != nil {
		fmt.Printf("Code: %v\n", list.Error.Code)
		fmt.Printf("Message: %v\n\n", list.Error.Message)
		// log.Fatalln("Err happened", err)
	}

	var tmpPubkey string
	tmpPubkey = ""

	for _, v := range list.Result.Matches {
		// fmt.Printf("\n-------\n")
		// fmt.Println(i)

		if tmpPubkey == v.Decrypted {
			// fmt.Println("Temp Pubkey matched")
		} else {
			// fmt.Println("Temp Pubkey did not match\nUpdated it's value")
			tmpPubkey = v.Decrypted

			// Checking if the found pubkey is authorized in subatomic.json pubkey list
			authorized, err := MatchedAuthorized(v.Decrypted)
			if err != nil {
				Log.Fatal(err)
			}
			// fmt.Println(authorized)

			handles = append(handles, DEXHandle{
				Pubkey:     v.Decrypted,
				Handle:     v.TagB,
				DEXPubkey:  v.Pubkey,
				Authorised: authorized,
			})
		}

	}

	// fmt.Println(len(handles))
	// fmt.Println(handles)
	// fmt.Println(handles[0])
	// fmt.Println(handles[1])

	return handles
}

// DLSubJSONData downloads the latest updated subatomic.json file data and saves it to "assets/subatomic.json" file in subatomicgo app
func DLSubJSONData() error {
	res, err := http.Get("https://raw.githubusercontent.com/jl777/komodo/jl777/src/cc/dapps/subatomic.json")
	if err != nil {
		log.Fatal(err)
	}
	subJSONData, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	if res.Status == "200 OK" {
		ioutil.WriteFile("assets/subatomic.json", subJSONData, 0644)
		return nil
	}
	return errors.New(res.Status)
}

// MatchedAuthorized checks the pubkey against the subatomic.json file's authorized pubkey list.
// If the searched pubkey is present in that list, it returns true, otherwise false.
func MatchedAuthorized(pubkey string) (bool, error) {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println(dir)
	data, err := ioutil.ReadFile(filepath.Join(dir, "/assets/subatomic.json"))
	if err != nil {
		log.Printf("%s", err)
		return false, err
	}
	// fmt.Printf("%s", data)

	// dec := json.NewDecoder(strings.NewReader(string(data[:])))
	// fmt.Println(dec)

	var parsed map[string][]map[string]string
	err = json.Unmarshal([]byte(data), &parsed)
	// fmt.Println(parsed["authorized"])
	// var auth map[string][string]interface{}

	for _, v := range parsed["authorized"] {
		// fmt.Println(v)
		for _, pub := range v {
			// fmt.Println(name)
			if pub == pubkey {
				// fmt.Println(pub)
				// fmt.Println("pubkey:", pubkey)
				return true, nil
			}
		}
	}
	return false, nil
}

// OrderData type is used to get formated data to display on Orderbook page
type OrderData struct {
	Price        string  `json:"price"`
	MaxVolume    string  `json:"max_volume"`
	DexPubkey    string  `json:"dex_pubkey"`
	Base         string  `json:"base"`
	ZBase        bool    `json:"zbase"`
	Rel          string  `json:"rel"`
	ZRel         bool    `json:"zrel"`
	OrderID      int64   `json:"orderid"`
	TimestampStr string  `json:"timestampstr"`
	Timestamp    string  `json:"timestamp"`
	Handle       string  `json:"handle"`
	Pubkey       string  `json:"pubkey"`
	Authorized   bool    `json:"authorized"`
	BaseBal      float64 `json:"basebal"`
	ZBaseBal     float64 `json:"zbasebal"`
	RelBal       float64 `json:"relbal"`
	ZRelBal      float64 `json:"zrelbal"`
	BaseIcon     string  `json:"baseicon"`
	RelIcon      string  `json:"relicon"`
}

// IsLower returns true if supplied string is lower case, false if upper case
func IsLower(s string) bool {
	for _, r := range s {
		if !unicode.IsLower(r) && unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

// OrderBookList returns processed data for Orderbook page
func OrderBookList(base, rel, maxentries, sortby string) []OrderData {
	var orderList []OrderData

	appName := kmdgo.AppType(DexP2pChain)

	var obook kmdgo.DEXOrderbook

	args := make(kmdgo.APIParams, 4)
	// maxentries eg. "10"
	args[0] = maxentries
	// minpriority
	args[1] = "0"
	// tagA/Base eg. "KMD"
	args[2] = base
	// tagB/Rel eg. "DEX"
	args[3] = rel
	// fmt.Println(args)

	// Debug outputs
	fmt.Println("compiled command is:")
	fmt.Println("dex-cli DEX_orderbook", maxentries, args[1], base, rel, " | jq .asks")
	fmt.Println(`
	from Asks :-
	buying = ` + base + `
	selling = ` + rel + `

	from Bids :-
	buying = ` + rel + `
	selling = ` + base + `
	`)

	obook, err := appName.DEXOrderbook(args)
	if err != nil {
		fmt.Printf("Code: %v\n", obook.Error.Code)
		fmt.Printf("Message: %v\n\n", obook.Error.Message)
		// log.Fatalln("Err happened", err)
	}

	for _, v := range obook.Result.Asks {
		handle, pubkey, auth := GetHandle(v.Pubkey)
		// fmt.Println(handle)
		// fmt.Println(pubkey)
		// fmt.Println(auth)

		var unixTime int64 = int64(v.Timestamp)
		t := time.Unix(unixTime, 0)
		// fmt.Println("t", t)
		strDate := t.Format(time.UnixDate)
		// fmt.Println(strDate)

		orderList = append(orderList, OrderData{
			Price: v.Price,
			// MaxVolume:  v.Relamount,
			MaxVolume:    v.Baseamount,
			DexPubkey:    v.Pubkey,
			Base:         obook.Result.Rel,
			Rel:          obook.Result.Base,
			OrderID:      v.ID,
			TimestampStr: strDate,
			Timestamp:    IntToString(int32(v.Timestamp)),
			Handle:       handle,
			Pubkey:       pubkey,
			Authorized:   auth,
		})
	}

	sort.Slice(orderList, func(i, j int) bool {
		return orderList[i].Authorized
	})

	// Sort by the soon to expire orders
	if sortby == "soon" {
		sort.Slice(orderList, func(i, j int) bool {
			return orderList[i].Timestamp < orderList[j].Timestamp
		})
	}
	// Sort by the late to expire orders
	if sortby == "late" {
		sort.Slice(orderList, func(i, j int) bool {
			return orderList[i].Timestamp > orderList[j].Timestamp
		})
	}

	// TODO
	// Sort by low/high price
	// Sort by Maxvolume
	// Sort by authorised pubkeys first
	// Sort by autorhised pubkeys first + other sorted formats

	// fmt.Println(orderList)
	return orderList
}

// GetHandle returns Handle, Public Key and Authorized status of that pubkey
func GetHandle(pubkey string) (string, string, bool) {
	var handles []DEXHandle
	handles = DEXHandles()

	for _, value := range handles {
		// fmt.Println(index)
		// fmt.Println(value)
		if pubkey == value.DEXPubkey {
			// fmt.Println(value.Handle, value.Pubkey, value.Authorised)
			return value.Handle, value.Pubkey, value.Authorised
		}
	}

	return "", "", false
}

// OrderID returns data to display on orderbook/{orderid} page
func OrderID(id string) OrderData {

	var orderData OrderData

	appName := kmdgo.AppType(DexP2pChain)

	var orderid kmdgo.DEXGet

	args := make(kmdgo.APIParams, 1)
	args[0] = id
	// args[1] = "0138d849d6bc81ff1c5389aae9a60ba3ee9cfd7858d93a3864679c25937e70951f"
	// args[2] = "BTC"
	// args[3] = "KMD"
	// fmt.Println(args)

	orderid, err := appName.DEXGet(args)
	if err != nil {
		fmt.Printf("Code: %v\n", orderid.Error.Code)
		fmt.Printf("Message: %v\n\n", orderid.Error.Message)
		// log.Fatalln("Err happened", err)
	}

	// fmt.Println(orderid.Result)

	handle, pubkey, auth := GetHandle(orderid.Result.Pubkey)
	// fmt.Println(handle)
	// fmt.Println(pubkey)
	// fmt.Println(auth)

	var unixTime int64 = int64(orderid.Result.Timestamp)
	t := time.Unix(unixTime, 0)
	strDate := t.Format(time.UnixDate)
	// fmt.Println(strDate)̉̉̉

	amountA, err := strconv.ParseFloat(orderid.Result.AmountA, 64)
	if err != nil {
		fmt.Println(err)
	}
	amountB, err := strconv.ParseFloat(orderid.Result.AmountB, 64)
	if err != nil {
		fmt.Println(err)
	}

	// fmt.Println("amountA:", amountA)
	// fmt.Println("amountB: ", amountB)
	price := amountB / amountA
	// fmt.Println("price:", price)

	var baseRelWallet = []kmdgo.AppType{kmdgo.AppType(orderid.Result.TagB), kmdgo.AppType(orderid.Result.TagA)}

	// fmt.Println(baseRelWallet)

	var wallets []WInfo
	wallets = WalletInfo(baseRelWallet)
	// fmt.Println(wallets[0].Balance)
	// fmt.Println(wallets[1].Balance)

	var relBalance, baseBalance float64
	if strings.HasPrefix(orderid.Result.TagB, "z") {
		baseBalance = wallets[0].ZBalance
	} else if strings.HasPrefix(orderid.Result.TagB, "PIRATE") {
		baseBalance = wallets[0].ZBalance
	} else {
		baseBalance = wallets[0].Balance
	}

	if strings.HasPrefix(orderid.Result.TagA, "z") {
		relBalance = wallets[1].ZBalance
	} else if strings.HasPrefix(orderid.Result.TagA, "PIRATE") {
		relBalance = wallets[1].ZBalance
	} else {
		relBalance = wallets[1].Balance
	}

	orderData = OrderData{
		Price:        fmt.Sprintf("%f", price),
		MaxVolume:    orderid.Result.AmountA,
		DexPubkey:    orderid.Result.Pubkey,
		Base:         orderid.Result.TagB,
		ZBase:        IsLower(orderid.Result.TagB[0:1]),
		Rel:          orderid.Result.TagA,
		ZRel:         IsLower(orderid.Result.TagA[0:1]),
		OrderID:      int64(orderid.Result.ID),
		TimestampStr: strDate,
		Timestamp:    IntToString(int32(orderid.Result.Timestamp)),
		Handle:       handle,
		Pubkey:       pubkey,
		Authorized:   auth,
		BaseBal:      baseBalance,
		RelBal:       relBalance,
		BaseIcon:     wallets[0].Icon,
		RelIcon:      wallets[1].Icon,
	}

	return orderData
}

// TxIDFromOpID returns TxID for provided opid and coin
func TxIDFromOpID(coin, opid string) (string, error) {
	var appName kmdgo.AppType
	appName = kmdgo.AppType(strings.ReplaceAll(string(coin), "z", ""))

	var oprst kmdgo.ZGetOperationStatus

	args := make(kmdgo.APIParams, 1)
	args[0] = []string{opid}
	// fmt.Println(args)

	oprst, err := appName.ZGetOperationStatus(args)
	if err != nil {
		fmt.Printf("Code: %v\n", oprst.Error.Code)
		fmt.Printf("Message: %v\n\n", oprst.Error.Message)
		Log.Println("Err happened", err)
	}

	for _, v := range oprst.Result {
		state6 := SwapStatus{
			State:    "opid_txid",
			Status:   "6",
			BaseTxID: v.Result.Txid,
		}
		state6JSON, _ := json.Marshal(state6)
		// fmt.Println("state6 JSON:", string(state6JSON))

		return string(state6JSON), nil
	}

	return "", nil
}

// IntToString Converts Int value to string
func IntToString(n int32) string {
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
