package sagoutil

import (
	"fmt"
	"golang-practice/kmdutil"
	"log"
	"os"
	"os/exec"
	"runtime"

	"github.com/Meshbits/shurli-server/sagoutil"
)

// StartWallet will launch Komodo-Ocean-QT with the specified Wallet
func StartWallet(chain, cmdParams string) error {
	fmt.Println(chain)

	// Check if provided blockchain is already running on system.
	// If chain's pid (ie. "komodo.pid") is present in that chain's data directory, it means
	// - that chain's daemon process is already running
	// - or the previous process did not delete the ie. "komodo.pid" file before exiting due to some reason, i.e. daemon crash etc.
	// 		- In this case, just delete the "komodo.pid" file and next time "shurli" should be able to start that blockchain.
	appName := chain
	dir := kmdutil.AppDataDir(appName, false)
	fmt.Println(dir)
	// If "chain" blockchain is running already, print notification
	if _, err := os.Stat(dir + "/komodod.pid"); err == nil {
		fmt.Println("[Shurli] " + chain + " blockchain already running or " + chain + " pid file exist.")
		sagoutil.Log.Println("[Shurli] " + chain + " blockchain already running or " + chain + " pid file exist.")
		os.Exit(1)
	} else {
		// If provided blockchain isn't found running already, start it.
		cmd := exec.Command("./komodo-qt-mac", cmdParams)
		if runtime.GOOS == "windows" {
			cmd = exec.Command("./komodo-qt-mac.exe", cmdParams)
		}
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Start()
		if err != nil {
			log.Fatalf("cmd.Run() failed with %s\n", err)
		}
		log.Printf("Started %s, with chain daemon params in background\n\t %s \nwith process ID: %d\n\n", chain, cmdParams, cmd.Process.Pid)
	}

	return nil
}
