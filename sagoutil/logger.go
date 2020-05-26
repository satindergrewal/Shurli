package sagoutil

import (
	"flag"
	"log"
	"os"
)

var (
	// Log ...
	Log *log.Logger
	// ShrliLogFile stores the stdout/stderr output log
	ShrliLogFile = "./shurli.log"
)

func init() {
	// set location of log file
	var logpath = ShrliLogFile

	flag.Parse()
	var file, err1 = os.Create(logpath)

	if err1 != nil {
		panic(err1)
	}
	Log = log.New(file, "", log.LstdFlags|log.Lshortfile)
	// Log.Println("LogFile : " + logpath)
}
