package utils

import (
	"log"
	"os"
	"time"
)

var ServLog *log.Logger

type PatientDetails struct {
	FullName     string
	HKID         string
	MedicalNotes string
	Email        string
}

type HyperledgerDetails struct {
	FullName         string
	HKID             string
	MedicalNotes     string
	Email            string
	FileName         string
	Size             int64
	ModificationTime time.Time
	EncryptImage     []byte
}

func InitLogger() {
	log_file, err := os.OpenFile("ServerLog.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	//defer log_file.Close()
	logger := log.New(log_file, "", log.Llongfile)
	ServLog = logger
}
