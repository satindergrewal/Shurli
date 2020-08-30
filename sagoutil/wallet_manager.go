package sagoutil

import (
	"encoding/json"
	"errors"
	"fmt"
	"golang-practice/kmdutil"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"runtime"
	"time"
	"unicode/utf8"

	"github.com/Meshbits/shurli-server/sagoutil"
	"github.com/satindergrewal/kmdgo"
)

// StartWallet launches Komodo-Ocean-QT with the specified Wallet
func StartWallet(chain string, cmdParams []string) error {

	//TODO: Add the capability to start 3rd party wallets, other than Komodo Assetchains

	// fmt.Println(chain)

	// fmt.Println(sagoutil.ShurliRootDir())

	// Check if provided blockchain is already running on system.
	// If chain's pid (ie. "komodo.pid") is present in that chain's data directory, it means
	// - that chain's daemon process is already running
	// - or the previous process did not delete the ie. "komodo.pid" file before exiting due to some reason, i.e. daemon crash etc.
	// 		- In this case, just delete the "komodo.pid" file and next time "shurli" should be able to start that blockchain.
	appName := chain
	dir := kmdutil.AppDataDir(appName, false)
	// fmt.Println(dir)
	// If "chain" blockchain is running already, print notification
	if _, err := os.Stat(dir + "/komodod.pid"); err == nil {
		return errors.New("wallet already running or it's process ID file exist")
	} else {
		// If provided blockchain isn't found running already, start it.
		cmd := exec.Command(sagoutil.ShurliRootDir()+"/assets/komodo-qt", cmdParams...)
		if runtime.GOOS == "windows" {
			cmd = exec.Command(sagoutil.ShurliRootDir()+"/assets/komodo-qt.exe", cmdParams...)
		}
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Start()
		if err != nil {
			// log.Printf("cmd.Start() failed with %s\n", err)
			return err
		}
		log.Printf("Started %s, with chain daemon params in background\n\t %s \nwith process ID: %d\n\n", chain, cmdParams, cmd.Process.Pid)
	}

	return nil
}

// BackupConfigJSON take backup of existing config.json file and store it with filename + timestamp
func BackupConfigJSON(confPath string) {
	// Get current time in unixtime format
	currentUnixTimestamp := sagoutil.IntToString(int32(time.Now().Unix()))
	// fmt.Println(currentUnixTimestamp)

	// create directory if it does't alredy exists
	if _, err := os.Stat(confPath + "/backups"); os.IsNotExist(err) {
		os.Mkdir(confPath+"/backups", 0755)
	}

	// read contents of existing config.json file
	from, err := os.Open(confPath + "/config.json")
	if err != nil {
		log.Fatal(err)
	}
	defer from.Close()

	// set copy file path to copy contents of config.json to
	to, err := os.OpenFile(confPath+"/backups/config_"+currentUnixTimestamp+".json", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer to.Close()

	// copy contents from config.json to backups/config_<unix time stamp>.json file
	_, err = io.Copy(to, from)
	if err != nil {
		log.Fatal(err)
	}
}

// GenerateDEXP2PAccount generate the transaparent address and shielded address
func GenerateDEXP2PAccount() error {
	// set DEXP2P2 chain's name to get RPC details
	appName := kmdgo.NewAppType(kmdgo.AppType(DexP2pChain))

	//Generate Transparent Address
	var DexP2PTransparentAddr kmdgo.GetNewAddress

	DexP2PTransparentAddr, err := appName.GetNewAddress()
	if err != nil {
		fmt.Printf("Code: %v\n", DexP2PTransparentAddr.Error.Code)
		fmt.Printf("Message: %v\n\n", DexP2PTransparentAddr.Error.Message)
		// log.Fatalln("Err happened", err)
		return errors.New(err.Error())
	}

	// fmt.Println(DexP2PTransparentAddr.Result)

	//Generate Shielded Address
	var DexP2PShieldedAddr kmdgo.ZGetNewAddress

	zAddrType := `sapling`

	DexP2PShieldedAddr, err = appName.ZGetNewAddress(zAddrType)
	if err != nil {
		fmt.Printf("Code: %v\n", DexP2PShieldedAddr.Error.Code)
		fmt.Printf("Message: %v\n\n", DexP2PShieldedAddr.Error.Message)
		// log.Fatalln("Err happened", err)
		return errors.New(err.Error())
	}

	// fmt.Println(DexP2PShieldedAddr.Result)

	// Get Transparent Address, Shielded Address and public key of newly generated address. Create new if doesn't exists, or store/Update config.json file.
	/// Get Transparent Address's public key
	var DexP2PPubkey kmdgo.ValidateAddress

	_DexP2PTAddr := DexP2PTransparentAddr.Result

	DexP2PPubkey, err = appName.ValidateAddress(_DexP2PTAddr)
	if err != nil {
		fmt.Printf("Code: %v\n", DexP2PPubkey.Error.Code)
		fmt.Printf("Message: %v\n\n", DexP2PPubkey.Error.Message)
		// log.Fatalln("Err happened", err)
		return errors.New(err.Error())
	}

	// fmt.Println("Pubkey: ", DexP2PPubkey.Result.Pubkey)

	// Generate a temporary random Handle based on pubkey
	_, i := utf8.DecodeRuneInString(DexP2PPubkey.Result.Pubkey)
	_tempHandle := "Anon" + DexP2PPubkey.Result.Pubkey[i:7]

	// Get contents of config.json sample file
	var conf SubAtomicConfig
	confJSONContent, err := ioutil.ReadFile("config.json.sample")
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(confJSONContent, &conf)

	// Generate new contents for config.json file and store newly generated address info to it
	var newConf sagoutil.SubAtomicConfig
	newConf.SubatomicExe = conf.SubatomicExe
	newConf.SubatomicDir = conf.SubatomicDir
	newConf.DexNSPV = conf.DexNSPV
	newConf.DexAddnode = conf.DexAddnode
	newConf.DexPubkey = DexP2PPubkey.Result.Pubkey      // Store public key of newly generated transparent address
	newConf.DexHandle = _tempHandle                     // A temporary handle generated based on public key's first 6 characters
	newConf.DexRecvTAddr = DexP2PTransparentAddr.Result // Store newly generated transparent address
	newConf.DexRecvZAddr = DexP2PShieldedAddr.Result    // Store newly generated shielded address

	// get indented JSON output of nelwy generated config.json
	var confJSON []byte
	confJSON, err = json.MarshalIndent(newConf, "", "	")
	if err != nil {
		return err
	}
	fmt.Println(string(confJSON))

	// Check if config.json already exists.
	// Take backup if exists before write new config.json file.
	// If doesn't, then create a new one
	_, err = os.Stat("config.json")
	if os.IsNotExist(err) {
		fmt.Println("config.json file does not exists. Creating a new one")
	} else {
		fmt.Println("config.json file already exists. Taking backup of it to backups/ directory")
		backupDir := sagoutil.ShurliRootDir()
		BackupConfigJSON(backupDir)
	}

	// Write newly genrated config.json to file
	err = ioutil.WriteFile("config.json", confJSON, 0644)
	if err != nil {
		return err
	}

	return nil
}

// ImportTAddrPrivKey import private key of DEXP2P transparent address to specified wallet
func ImportTAddrPrivKey(toChain string) error {
	// Get contents of config.json file
	var conf SubAtomicConfig
	confJSONContent, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(confJSONContent, &conf)

	// Get DEXP2P transparent address' private key using dumpprivkey
	var DexP2PDumpTPrivKey kmdgo.DumpPrivKey
	DexP2PDumpTPrivKey, err = kmdgo.NewAppType(kmdgo.AppType(DexP2pChain)).DumpPrivKey(conf.DexRecvTAddr)
	if err != nil {
		fmt.Printf("Code: %v\n", DexP2PDumpTPrivKey.Error.Code)
		fmt.Printf("Message: %v\n\n", DexP2PDumpTPrivKey.Error.Message)
		// log.Fatalln("Err happened", err)
		return err
	}
	// fmt.Println(DexP2PDumpTPrivKey.Result)

	// Import privkey to the target chain
	var DexP2PImpTPrivkey kmdgo.ImportPrivKey
	args := make(kmdgo.APIParams, 3)
	args[0] = DexP2PDumpTPrivKey.Result
	args[1] = `shurli`
	args[2] = false
	// fmt.Println(args)

	DexP2PImpTPrivkey, err = kmdgo.NewAppType(kmdgo.AppType(toChain)).ImportPrivKey(args)
	if err != nil {
		fmt.Printf("Code: %v\n", DexP2PImpTPrivkey.Error.Code)
		fmt.Printf("Message: %v\n\n", DexP2PImpTPrivkey.Error.Message)
		// log.Fatalln("Err happened", err)
		return err
	}
	// fmt.Println(DexP2PImpTPrivkey.Result)

	return nil
}

// ImportZAddrPrivKey import private key of DEXP2P shielded address to specified wallet
func ImportZAddrPrivKey(toChain string) error {
	// Get contents of config.json file
	var conf SubAtomicConfig
	confJSONContent, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(confJSONContent, &conf)

	// Get DEXP2P shielded address' private key using z_exportkey
	var DexP2PExportZPrivKey kmdgo.ZExportKey
	DexP2PExportZPrivKey, err = kmdgo.NewAppType(kmdgo.AppType(DexP2pChain)).ZExportKey(conf.DexRecvZAddr)
	if err != nil {
		fmt.Printf("Code: %v\n", DexP2PExportZPrivKey.Error.Code)
		fmt.Printf("Message: %v\n\n", DexP2PExportZPrivKey.Error.Message)
		// log.Fatalln("Err happened", err)
		return err
	}
	// fmt.Println(DexP2PExportZPrivKey.Result)

	// Import shielded address' privkey to the target chain
	var DexP2PImpZPrivkey kmdgo.ZImportKey
	args := make(kmdgo.APIParams, 3)
	args[0] = DexP2PExportZPrivKey.Result
	args[1] = `no`
	args[2] = 0
	// fmt.Println(args)

	DexP2PImpZPrivkey, err = kmdgo.NewAppType(kmdgo.AppType(toChain)).ZImportKey(args)
	if err != nil {
		fmt.Printf("Code: %v\n", DexP2PImpZPrivkey.Error.Code)
		fmt.Printf("Message: %v\n\n", DexP2PImpZPrivkey.Error.Message)
		// log.Fatalln("Err happened", err)
		return err
	}
	// fmt.Println(DexP2PImpZPrivkey.Result)

	return nil
}

// UpdateDEXP2PAccount allow users to update config.json file with user specified DEXP2P params details
func UpdateDEXP2PAccount(data SubAtomicConfig) error {
	//Debug print of checking what data we are getting
	// fmt.Println(data)
	// fmt.Println(len(data.DexHandle))

	// Get contents of config.json file
	var conf SubAtomicConfig
	confJSONContent, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(confJSONContent, &conf)

	if len(data.DexNSPV) != 0 {
		conf.DexNSPV = data.DexNSPV
	}
	if len(data.DexAddnode) != 0 {
		conf.DexAddnode = data.DexAddnode
	}
	if len(data.DexPubkey) != 0 {
		conf.DexPubkey = data.DexPubkey
	}
	if len(data.DexHandle) != 0 {
		conf.DexHandle = data.DexHandle
	}
	if len(data.DexRecvTAddr) != 0 {
		conf.DexRecvTAddr = data.DexRecvTAddr
	}
	if len(data.DexRecvZAddr) != 0 {
		conf.DexRecvZAddr = data.DexRecvZAddr
	}

	// fmt.Println(conf.DexNSPV)
	// fmt.Println(conf.DexAddnode)
	// fmt.Println(conf.DexPubkey)
	// fmt.Println(conf.DexHandle)
	// fmt.Println(conf.DexRecvTAddr)
	// fmt.Println(conf.DexRecvZAddr)

	// get indented JSON output of nelwy generated config.json
	var confJSON []byte
	confJSON, err = json.MarshalIndent(conf, "", "	")
	if err != nil {
		return err
	}
	// fmt.Println(string(confJSON))

	// Write newly genrated config.json to file
	err = ioutil.WriteFile("config.json", confJSON, 0644)
	if err != nil {
		return err
	}

	return nil
}

// DlBootstrap download, extract and replace/update bootstrap blockchain files for a specified wallet
func DlBootstrap() {

}

// BackupWallet allows taking a dump or backup of the wallet.dat file
func BackupWallet() {

}
