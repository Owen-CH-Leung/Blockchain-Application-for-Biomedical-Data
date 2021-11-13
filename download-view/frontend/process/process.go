package process

import (
	"crypto/cipher"
	"download-view/backend/encryption"
	"download-view/backend/image"
	b_logger "download-view/backend/logger"
	"download-view/backend/readfabric"
	f_utils "download-view/frontend/utils"
	"encoding/base64"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func DownloadKeyFile(r *http.Request) (string, string, string) {
	log := f_utils.ServLog
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		log.Println("Parsing MultiPart Form encounters error")
		log.Println(err)
		os.Exit(1)
	}

	hkid := r.FormValue("HKID")
	email := r.FormValue("Email")
	// do something with details
	log.Println("HKID :", hkid)
	log.Println("Email :", email)

	file, handler, err := r.FormFile("Keyfile")
	if err != nil {
		log.Println("Error Retrieving the Key File")
		log.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	log.Printf("Uploaded File: %+v\n", handler.Filename)
	// log.Printf("File Size: %+v\n", handler.Size)
	// log.Printf("MIME Header: %+v\n", handler.Header)

	keypath := filepath.Join("..", "downloaded_key")
	err = os.MkdirAll(keypath, os.ModePerm)
	if err != nil {
		log.Println("Creating downloaded_key folder faces error")
		log.Println(err)
	}
	download_file, err := os.OpenFile(filepath.Join(keypath, handler.Filename), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
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
	privatekeypath := filepath.Join(keypath, handler.Filename)
	log.Println("Private Key Path: ")
	log.Println(privatekeypath)

	return email, hkid, privatekeypath
}

func GetAESKeyFiles(hkid string) string {
	keypath := filepath.Join("..", "key")
	dir_files, err := ioutil.ReadDir(filepath.Clean(keypath))
	if err != nil {
		log.Fatal(err)
	}

	var target_path string
	for _, file := range dir_files {
		file_hkid := strings.Split(file.Name(), "_")[0]
		if file_hkid == hkid && strings.Contains(file.Name(), ".key") {
			log.Println("Found Key")
			log.Println(file.Name())
			target_path = filepath.Clean(filepath.Join(keypath, file.Name()))
		}
	}
	log.Println("Target Path for AES Key")
	log.Println(target_path)

	return target_path
}

func GetPrivateKeyFiles(hkid string) string {
	keypath := filepath.Join("..", "key")
	dir_files, err := ioutil.ReadDir(filepath.Clean(keypath))
	if err != nil {
		log.Fatal(err)
	}

	var target_path string
	for _, file := range dir_files {
		file_hkid := strings.Split(file.Name(), "_")[0]
		if file_hkid == hkid && strings.Contains(file.Name(), "_Private.pem") {
			log.Println("Found Key")
			log.Println(file.Name())
			target_path = filepath.Clean(filepath.Join(keypath, file.Name()))
		}
	}
	log.Println("Target Path for Private Key")
	log.Println(target_path)

	return target_path
}

func DecryptAndSaveAsJson(hkid string, info string, aespath string, privatekeypath string) []string {
	currenttime := b_logger.GetCurrentTime()
	logfile := hkid + "_" + currenttime + ".log"
	b_logger.InitLogger(logfile)

	record, err := readfabric.ReadLatest(hkid, info)
	if err != nil {
		log.Println("Reading record from Fabric encounters Error")
		log.Println(err)
		os.Exit(1)
	}
	log.Println("record is :", record)
	log.Printf("It is of type %T", record)
	aeskey := encryption.ReadAESKey(aespath)
	privatekey := encryption.ReadPrivateKey(privatekeypath)
	// cb, nonce := encryption.GetCipherAndNonce(aeskey)
	decrypted_AES := encryption.DecryptKeyByPrivateKey(privatekey, aeskey)
	cb, nonce := encryption.GetCipherAndNonce(decrypted_AES)

	DecryptStringsInStruct(reflect.ValueOf(record).Elem(), cb, nonce)
	var attachment []string
	if info == "chestpainreferral" {
		filename := hkid + "_chestXray.jpg"
		new_record := record.(*f_utils.ChestPainReferral)
		img_path := image.ByteArrayToJpg([]byte(new_record.ChestXRayImage), filename)
		attachment = append(attachment, img_path)
	}
	//TODO : Add logic for chestpain, marshal without the image field
	jsonByte, err := json.MarshalIndent(record, "", " ")
	if err != nil {
		log.Println("JSON encoding encounters Error")
		log.Println(err)
		os.Exit(1)
	}

	jsonpath := filepath.Join("..", "sent_files")
	err = os.MkdirAll(jsonpath, os.ModePerm)
	if err != nil {
		log.Println("Creating json directory encounters error")
		log.Println(err)
	}
	target_path := hkid + "_" + info + "_SENT.json"
	_ = ioutil.WriteFile(filepath.Join(jsonpath, target_path), jsonByte, 0644)
	filepath := filepath.Join(jsonpath, target_path)
	attachment = append(attachment, filepath)
	return attachment
}

func DecryptStringsInStruct(v reflect.Value, cb cipher.Block, nonce []byte) {
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)
		log.Println("Handling Field - ", t.Field(i).Name)
		switch sf.Type.Kind() {
		case reflect.Struct:
			DecryptStringsInStruct(v.Field(i), cb, nonce)
		case reflect.String:
			field_str := v.Field(i).Interface().(string)
			//log.Println("The field_str is : ", field_str)
			encrypted_byte_arr, err := base64.StdEncoding.DecodeString(field_str)
			if err != nil {
				log.Println("Convert Base64 to Byte encounters Error")
				log.Println(err)
				os.Exit(1)
			}
			//log.Println("B64-Decoded string = ", string(encrypted_byte_arr))
			decrypted_img := encryption.DecryptDataByAESKey(cb, nonce, encrypted_byte_arr)
			//log.Println("Decrypted string = ", string(decrypted_img))
			v.Field(i).SetString(string(decrypted_img))
			// decrypted_img := encryption.DecryptDataByAESKey(cb, nonce, []byte(field_str))
			// v.Field(i).SetString(string(decrypted_img))
		case reflect.Slice:
			field_str := v.Field(i).Interface().([]byte)
			encryptedmsg := encryption.DecryptDataByAESKey(cb, nonce, field_str)
			v.Field(i).SetBytes(encryptedmsg)
			//log.Println("The Decrypted String is : ", string(encryptedmsg))
		}
	}
}

func DecryptAllAndSaveAsJson(hkid string, info string, aespath string, privatekeypath string) []string {
	currenttime := b_logger.GetCurrentTime()
	logfile := hkid + "_" + currenttime + ".log"
	b_logger.InitLogger(logfile)
	aeskey := encryption.ReadAESKey(aespath)
	privatekey := encryption.ReadPrivateKey(privatekeypath)
	// cb, nonce := encryption.GetCipherAndNonce(aeskey)
	decrypted_AES := encryption.DecryptKeyByPrivateKey(privatekey, aeskey)
	cb, nonce := encryption.GetCipherAndNonce(decrypted_AES)
	switch {
	case info == "patient":
		record, err := readfabric.ReadPatient_All(hkid)
		if err != nil {
			log.Println("Reading record from Fabric encounters Error")
			log.Println(err)
			os.Exit(1)
		}
		for idx, rec := range *record {
			//log.Println("The record looks like : ", &rec)
			//log.Printf("It is of type %T", &rec)
			DecryptStringsInStruct(reflect.ValueOf(&rec).Elem(), cb, nonce)
			//log.Println("Inside For Loop : rec = ", rec)
			(*record)[idx] = rec
		}
		//log.Println("After Decrypting : Record = ", record)
		attachment := GetAttachmentList(hkid, info, record)
		return attachment

	case info == "vaccine":
		record, err := readfabric.ReadVaccine_All(hkid)
		if err != nil {
			log.Println("Reading record from Fabric encounters Error")
			log.Println(err)
			os.Exit(1)
		}
		for idx, rec := range *record {
			DecryptStringsInStruct(reflect.ValueOf(&rec).Elem(), cb, nonce)
			(*record)[idx] = rec
		}
		attachment := GetAttachmentList(hkid, info, record)
		return attachment

	case info == "allergy":
		record, err := readfabric.ReadAllergy_All(hkid)
		if err != nil {
			log.Println("Reading record from Fabric encounters Error")
			log.Println(err)
			os.Exit(1)
		}
		for idx, rec := range *record {
			DecryptStringsInStruct(reflect.ValueOf(&rec).Elem(), cb, nonce)
			(*record)[idx] = rec
		}
		attachment := GetAttachmentList(hkid, info, record)
		return attachment

	case info == "bodymeasurement":
		record, err := readfabric.ReadBodyMeasurement_All(hkid)
		if err != nil {
			log.Println("Reading record from Fabric encounters Error")
			log.Println(err)
			os.Exit(1)
		}
		for idx, rec := range *record {
			DecryptStringsInStruct(reflect.ValueOf(&rec).Elem(), cb, nonce)
			(*record)[idx] = rec
		}
		attachment := GetAttachmentList(hkid, info, record)
		return attachment

	case info == "mealplan":
		record, err := readfabric.ReadMealPlan_All(hkid)
		if err != nil {
			log.Println("Reading record from Fabric encounters Error")
			log.Println(err)
			os.Exit(1)
		}
		for idx, rec := range *record {
			DecryptStringsInStruct(reflect.ValueOf(&rec).Elem(), cb, nonce)
			(*record)[idx] = rec
		}
		attachment := GetAttachmentList(hkid, info, record)
		return attachment

	case info == "assessment":
		record, err := readfabric.ReadAssessment_All(hkid)
		if err != nil {
			log.Println("Reading record from Fabric encounters Error")
			log.Println(err)
			os.Exit(1)
		}
		for idx, rec := range *record {
			DecryptStringsInStruct(reflect.ValueOf(&rec).Elem(), cb, nonce)
			(*record)[idx] = rec
		}
		attachment := GetAttachmentList(hkid, info, record)

		return attachment

	case info == "riskfactor":
		record, err := readfabric.ReadRiskFactor_All(hkid)
		if err != nil {
			log.Println("Reading record from Fabric encounters Error")
			log.Println(err)
			os.Exit(1)
		}
		for idx, rec := range *record {
			DecryptStringsInStruct(reflect.ValueOf(&rec).Elem(), cb, nonce)
			(*record)[idx] = rec
		}
		attachment := GetAttachmentList(hkid, info, record)
		return attachment

	case info == "procedure":
		record, err := readfabric.ReadProcedure_All(hkid)
		if err != nil {
			log.Println("Reading record from Fabric encounters Error")
			log.Println(err)
			os.Exit(1)
		}
		for idx, rec := range *record {
			DecryptStringsInStruct(reflect.ValueOf(&rec).Elem(), cb, nonce)
			(*record)[idx] = rec
		}
		attachment := GetAttachmentList(hkid, info, record)
		return attachment

	case info == "chestpainreferral":
		record, err := readfabric.ReadChestPain_All(hkid)
		if err != nil {
			log.Println("Reading record from Fabric encounters Error")
			log.Println(err)
			os.Exit(1)
		}
		for idx, rec := range *record {
			DecryptStringsInStruct(reflect.ValueOf(&rec).Elem(), cb, nonce)
			(*record)[idx] = rec
		}
		attachment := GetAttachmentList(hkid, info, record)
		for idx, rec := range *record {
			filename := hkid + "_chestXray_" + strconv.Itoa(idx) + ".jpg"
			//new_record := record.(*f_utils.ChestPainReferral)
			img_path := image.ByteArrayToJpg([]byte(rec.ChestXRayImage), filename)
			attachment = append(attachment, img_path)
		}
		return attachment

	default:
		record, err := readfabric.ReadPatient_All(hkid)
		if err != nil {
			log.Println("Reading record from Fabric encounters Error")
			log.Println(err)
			os.Exit(1)
		}
		for _, rec := range *record {
			DecryptStringsInStruct(reflect.ValueOf(rec).Elem(), cb, nonce)
		}
		attachment := GetAttachmentList(hkid, info, record)
		return attachment
	}
}

func GetAttachmentList(hkid string, info string, record interface{}) []string {
	var attachment []string
	jsonByte, err := json.MarshalIndent(record, "", " ")
	if err != nil {
		log.Println("JSON encoding encounters Error")
		log.Println(err)
		os.Exit(1)
	}

	jsonpath := filepath.Join("..", "sent_files")
	err = os.MkdirAll(jsonpath, os.ModePerm)
	if err != nil {
		log.Println("Creating json directory encounters error")
		log.Println(err)
	}
	target_path := hkid + "_" + info + "_ALL_SENT.json"
	_ = ioutil.WriteFile(filepath.Join(jsonpath, target_path), jsonByte, 0644)
	filepath := filepath.Join(jsonpath, target_path)
	attachment = append(attachment, filepath)
	return attachment
}

func GetTotalNumPatient() []string {
	keypath := filepath.Join("..", "key")
	dir_files, err := ioutil.ReadDir(filepath.Clean(keypath))
	if err != nil {
		log.Fatal(err)
	}
	var hkid_arr []string
	for _, file := range dir_files {
		file_hkid := strings.Split(file.Name(), "_")[0]
		if !CheckExist(hkid_arr, file_hkid) {
			hkid_arr = append(hkid_arr, file_hkid)
		}
	}
	return hkid_arr
}

func CheckExist(arr_str []string, str string) bool {
	for _, record := range arr_str {
		if record == str {
			return true
		}
	}
	return false
}

func CheckRecordMatchPeriod(start string, end string, file_name string) bool {
	file_time := strings.Split(file_name, "_")[1]
	start_date, _ := time.Parse("2006-01-02", start)
	end_date, _ := time.Parse("2006-01-02", end)
	check_date, _ := time.Parse("20060102", file_time)
	start_date_minus1 := start_date.AddDate(0, 0, -1)
	end_date_plus1 := end_date.AddDate(0, 0, 1)
	return InTimeSpan(start_date_minus1, end_date_plus1, check_date)
}

func InTimeSpan(start, end, check time.Time) bool {
	return check.After(start) && check.Before(end)
}

func Get_Patient_Vaccine_ChestPain_Info(r *http.Request) ([]string, []string, []string) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		log.Println("Parsing MultiPart Form encounters error")
		log.Println(err)
		os.Exit(1)
	}
	ID := 0
	start_date := r.FormValue("date_from")
	end_date := r.FormValue("date_to")
	hkid_arr := GetTotalNumPatient()
	vaccine_path := filepath.Join("..", "vaccine")
	chestpain_path := filepath.Join("..", "chestpain")
	dir_files, _ := ioutil.ReadDir(filepath.Clean(vaccine_path))
	chest_dir_files, _ := ioutil.ReadDir(filepath.Clean(chestpain_path))
	patient_arr := []string{}
	vaccine_arr := []string{}
	chestpain_arr := []string{}
	for _, hkid := range hkid_arr {
		mask_name := "Patient" + strconv.Itoa(ID)
		patient_arr = append(patient_arr, mask_name)
		for _, file := range dir_files {
			file_hkid := strings.Split(file.Name(), "_")[0]
			if file_hkid == hkid {
				if CheckRecordMatchPeriod(start_date, end_date, file.Name()) {
					vaccine_arr = append(vaccine_arr, mask_name)
					for _, another_file := range chest_dir_files {
						chest_hkid := strings.Split(another_file.Name(), "_")[0]
						if hkid == chest_hkid {
							chestpain_arr = append(chestpain_arr, mask_name)
						}
					}
				}
			} else {
				continue
			}
		}
		ID += 1
	}
	return patient_arr, vaccine_arr, chestpain_arr
}
