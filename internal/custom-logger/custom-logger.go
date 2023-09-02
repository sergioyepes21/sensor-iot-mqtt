package logger

import (
	"flag"
	"log"
	"os"
	"path/filepath"
)

var (
	Log *log.Logger
)

func init() {
	// set location of log file
	var logpath, err = filepath.Abs("./assets/anomalies.log")
	if err != nil {
		panic(err)
	}

	flag.Parse()
	var file, err1 = os.Create(logpath)

	if err1 != nil {
		panic(err1)
	}
	Log = log.New(file, "", log.LstdFlags|log.Lmsgprefix)
	Log.Println("LogFile : " + logpath)
}
