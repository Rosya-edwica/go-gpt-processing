package logger

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var (
	Log *log.Logger
)

func init() {
	logpath := "info.log"
	flag.Parse()
	file, err := os.OpenFile(logpath, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	defer file.Close()
	if err != nil {
		panic(err)
	}
	Log = log.New(file,  "", log.LstdFlags)
	fmt.Println("Logfile: " + logpath)
}