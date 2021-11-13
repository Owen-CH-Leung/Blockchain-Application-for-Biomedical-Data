package logger

import (
	"log"
	"os"
	"time"
	"path/filepath"
	"fmt"
)

var BlockchainLog *log.Logger

func InitLogger(logname string) {
	logpath := filepath.Join(".", "log")
	err := os.MkdirAll(logpath, os.ModePerm)
	if err != nil {
		fmt.Println("Unexpected Error: ")
		fmt.Println(err)
		os.Exit(1)
	}
	log_file, err := os.OpenFile(filepath.Join(logpath, logname), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Unexpected Error: ")
		fmt.Println(err)
		os.Exit(1)
	}
	//defer log_file.Close()
	logger := log.New(log_file, "", log.Llongfile)
	BlockchainLog = logger
}


func GetCurrentTime() string {
	current_t := time.Now()
	formatedTime := current_t.Format("20060102_150405")
	return formatedTime
}
