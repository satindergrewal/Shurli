package sagoutil

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/satindergrewal/kmdgo"
)

// ConfigCoins type holds the values of coin information
type ConfigCoins struct {
	Coin     string `json:"coin"`
	Ticker   string `json:"ticker"`
	Name     string `json:"name"`
	Shielded bool   `json:"shielded"`
}

// SubAtomicConfig holds the app's confugration settings
type SubAtomicConfig struct {
	Chains       []ConfigCoins     `json:"chains,omitempty"`
	SubatomicExe string            `json:"subatomic_exe"`
	SubatomicDir string            `json:"subatomic_dir"`
	DexNSPV      string            `json:"dex_nSPV"`
	DexAddnode   string            `json:"dex_addnode"`
	DexPubkey    string            `json:"dex_pubkey"`
	DexHandle    string            `json:"dex_handle"`
	DexRecvZAddr string            `json:"dex_recvzaddr"`
	DexRecvTAddr string            `json:"dex_recvtaddr"`
	Explorers    map[string]string `json:"explorers,omitempty"`
}

// ConfigChains holds the list of chains and it's other details like explorers links etc.
type ConfigChains struct {
	Chains    []ConfigCoins     `json:"chains"`
	Explorers map[string]string `json:"explorers"`
}

//SubAtomicConfInfo returns application's config params
func SubAtomicConfInfo() SubAtomicConfig {
	var conf SubAtomicConfig
	confJSONContent, err := ioutil.ReadFile("config.json")
	if err != nil {
		fmt.Println("config.json file not found")
		log.Fatal(err)
	}
	err = json.Unmarshal(confJSONContent, &conf)
	// fmt.Println("conf1:", conf)

	var chains ConfigChains
	chainsJSONContent, err := ioutil.ReadFile("chains.json")
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println(chainsJSONContent)
	err = json.Unmarshal(chainsJSONContent, &chains)
	conf.Chains = chains.Chains
	conf.Explorers = chains.Explorers
	// fmt.Println("chains: ", chains)
	// fmt.Println("conf2:", conf)
	// fmt.Println(chains.Explorers["VRSC"])
	return conf
}

//StrToAppType converts and returns slice of string as slice of kmdgo.AppType
func StrToAppType(chain []ConfigCoins) []kmdgo.AppType {
	var chainskmd []kmdgo.AppType
	for _, v := range chain {
		// fmt.Println(v.Coin)
		chainskmd = append(chainskmd, kmdgo.AppType(v.Coin))
	}
	return chainskmd
}

// GetCoinConfInfo returns single coin info from config.json cofiguration list
func GetCoinConfInfo(coin string) ConfigCoins {
	var conf SubAtomicConfig = SubAtomicConfInfo()
	var confChains []ConfigCoins = conf.Chains
	// fmt.Println(confChains)

	if coin == "KMD" || coin == "Komodo" {
		coin = "komodo"
	}

	// fmt.Println("getCoinConfInfo:", coin)

	var coininfo ConfigCoins
	for _, v := range confChains {
		// fmt.Println(v)
		if coin == v.Coin {
			coininfo = v
			return coininfo
		}
	}

	return coininfo
}

// ShurliStartMsg printing message at start
func ShurliStartMsg() {
	// 	heart := `
	//    ░█████   ░█████
	//  ░████████░█████████
	//  ░██████████████████
	//   ░████████████████
	//     ░████████████
	//        ░███████
	//          ░██
	// 	 `
	// 	fmt.Println(heart)

	iLShurli := `
	█████      ░█████   ░█████         █████████  █████                           ████   ███ 
	░░███    ░████████░█████████      ███░░░░░███░░███                           ░░███  ░░░  
	 ░███    ░██████████████████     ░███    ░░░  ░███████   █████ ████ ████████  ░███  ████ 
	 ░███     ░████████████████      ░░█████████  ░███░░███ ░░███ ░███ ░░███░░███ ░███ ░░███ 
	 ░███       ░████████████         ░░░░░░░░███ ░███ ░███  ░███ ░███  ░███ ░░░  ░███  ░███ 
	 ░███         ░███████            ███    ░███ ░███ ░███  ░███ ░███  ░███      ░███  ░███ 
	 █████           ░██             ░░█████████  ████ █████ ░░████████ █████     █████ █████
	░░░░░                              ░░░░░░░░░  ░░░░ ░░░░░   ░░░░░░░░ ░░░░░     ░░░░░ ░░░░░ 
	`

	fmt.Println(iLShurli)
}
