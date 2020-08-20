package sagoutil

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

// StartWallet will launch Komodo-Ocean-QT with the specified Wallet
func StartWallet(chain, pubkey string) error {
	fmt.Println(chain)

	cmdParams := "-pubkey" + pubkey

	cmd := exec.Command("./komodo-qt-mac", cmdParams)
	// if runtime.GOOS == "windows" {
	// 	cmd = exec.Command("tasklist")
	// }
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	log.Printf("Started %s, with public key %s in background with process ID: %d", chain, pubkey, cmd.Process.Pid)
	return nil
}
