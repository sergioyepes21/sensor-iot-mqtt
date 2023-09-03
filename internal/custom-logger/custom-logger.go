package logger

import (
	"flag"
	"log"
	"os"
)

var (
	Log *log.Logger
)

func init() {

	var logpath = "./assets/anomalies.log"

	flag.Parse()
	var file, err1 = os.Create(logpath)

	if err1 != nil {
		panic(err1)
	}
	Log = log.New(file, "", log.LstdFlags|log.Lmsgprefix)
	Log.Println("LogFile : " + logpath)
}
