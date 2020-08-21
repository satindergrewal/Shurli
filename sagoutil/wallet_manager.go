package sagoutil

import (
	"errors"
	"golang-practice/kmdutil"
	"log"
	"os"
	"os/exec"
	"runtime"

	"github.com/Meshbits/shurli-server/sagoutil"
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
func BackupConfigJSON() {
}

// GenerateDEXP2PAccount generate the transaparent address and shielded address
func GenerateDEXP2PAccount() {
	//Generate Transparent Address
	//Generate Shielded Address
	//Get Transparent Address, Shielded Address and public key of newly generated address. Create new if doesn't exists, or store/Update config.json file.
}

// ImportTAddrPrivKey import private key of DEXP2P transparent address to specified wallet
func ImportTAddrPrivKey() {

}

// ImportZAddrPrivKey import private key of DEXP2P shielded address to specified wallet
func ImportZAddrPrivKey() {

}

// UpdateDEXP2PAccount allow users to update config.json file with user specified DEXP2P params details
func UpdateDEXP2PAccount() {

}

// DlBootstrap download, extract and replace/update bootstrap blockchain files for a specified wallet
func DlBootstrap() {

}

// BackupWallet allows taking a dump or backup of the wallet.dat file
func BackupWallet() {

}
