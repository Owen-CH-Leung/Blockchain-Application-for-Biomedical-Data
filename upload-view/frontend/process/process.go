package process

import (
	"upload-view/backend/email"
	"upload-view/backend/encryption"
	"upload-view/backend/image"
	b_logger "upload-view/backend/logger"
	f_utils "upload-view/frontend/utils"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func DownloadFile(r *http.Request) (f_utils.PatientDetails, string) {
	log := f_utils.ServLog
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		log.Println("Parsing MultiPart Form encounters error")
		log.Println(err)
		os.Exit(1)
	}

	details := f_utils.PatientDetails{
		FullName:     r.FormValue("FullName"),
		HKID:         r.FormValue("HKID"),
		MedicalNotes: r.FormValue("MedicalNotes"),
		Email:        r.FormValue("Email"),
	}
	// do something with details
	log.Println("Full Name :", details.FullName)
	log.Println("HKID :", details.HKID)
	log.Println("MedicalNotes :", details.MedicalNotes)
	log.Println("Email :", details.Email)

	file, handler, err := r.FormFile("ChestXRay")
	if err != nil {
		log.Println("Error Retrieving the File")
		log.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	log.Printf("Uploaded File: %+v\n", handler.Filename)
	log.Printf("File Size: %+v\n", handler.Size)
	log.Printf("MIME Header: %+v\n", handler.Header)

	//name := strings.TrimSuffix(handler.Filename, filepath.Ext(handler.Filename))
	//download_file, err := os.OpenFile(name+"_new"+".jpg", os.O_WRONLY|os.O_CREATE, 0666)
	jpgpath := filepath.Join(".", "jpg")
	err = os.MkdirAll(jpgpath, os.ModePerm)
    if err != nil {
        log.Println("Creating jpg folder faces error")
        log.Println(err)
    }
	download_file, err := os.OpenFile(filepath.Join(jpgpath, handler.Filename), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Println("Error Creating a new File for copy")
		log.Println(err)
		os.Exit(1)
	}
	_, err = io.Copy(download_file, file)
	if err != nil {
		log.Println("Copying file encounters error")
		log.Println(err)
		os.Exit(1)
	}
	download_file.Close()

	//return details, handler.Filename
	return details, handler.Filename
}

func ReadEncryptEmail(src_img string, details f_utils.PatientDetails) (string, image.Metadata_jpg2json) {
	currenttime := b_logger.GetCurrentTime()
	logfile := details.HKID + "_" + currenttime + ".log"

	b_logger.InitLogger(logfile)

	destname := details.FullName

	jpgpath := filepath.Join(".", "jpg")
	mapping := image.JpgToMetaData(filepath.Join(jpgpath, src_img))

	private_key, public_key := encryption.InitatePubicPrivateKey()
	aes_key := encryption.GenerateAESKey()
	cb, nonce := encryption.GetCipherAndNonce(aes_key)

	encryptedmsg := encryption.EncryptDataByAESKey(cb, nonce, mapping.Image)
	encrypted_key := encryption.EncryptKeyByPublicKey(public_key, aes_key)

	mapping.Image = encryptedmsg

	Private_Key_Name, Public_Key_Name, encryption_key_name := encryption.CreateKeyFiles(public_key, private_key, encrypted_key, destname)

	keypath := filepath.Join(".", "key")
	prikeypath := filepath.Join(keypath, Private_Key_Name)
	pubkeypath := filepath.Join(keypath, Public_Key_Name)
	encryptkeypath := filepath.Join(keypath, encryption_key_name)

	recipient := details.Email
	attachment := []string{prikeypath, pubkeypath, encryptkeypath}
	email.SendEmailViaAWS(recipient, attachment...)

	basefile := strings.TrimSuffix(logfile, filepath.Ext(logfile))
	return basefile, mapping
}

func CreateBlockchainRecordJSON(basefile string, mapping image.Metadata_jpg2json, details f_utils.PatientDetails) {
	logfile := basefile + ".log"
	b_logger.InitLogger(logfile)
	log := b_logger.BlockchainLog
	log.Println("Creating Record for Hyperledger Fabric...")

	jsonpath := filepath.Join(".", "json")
	err := os.MkdirAll(jsonpath, os.ModePerm)
    if err != nil {
        log.Println("Creating json directory encounters error")
        log.Println(err)
    }

	destfile := basefile + ".json"
	record := f_utils.HyperledgerDetails{
		FullName:         details.FullName,
		HKID:             details.HKID,
		MedicalNotes:     details.MedicalNotes,
		Email:            details.Email,
		FileName:         mapping.FileName,
		Size:             mapping.Size,
		ModificationTime: mapping.ModificationTime,
		EncryptImage:     mapping.Image,
	}

	jsonByte, err := json.MarshalIndent(record, "", " ")
	if err != nil {
		log.Println("When creating blockchain record, JSON encoding encounters Error")
		log.Println(err)
		os.Exit(1)
	}
	err = ioutil.WriteFile(filepath.Join(jsonpath, destfile), jsonByte, 0644)
	if err != nil {
		log.Println("When creating blockchain record, writing JSON encounters Error")
		log.Println(err)
		os.Exit(1)
	}
}
