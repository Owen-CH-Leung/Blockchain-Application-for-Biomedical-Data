package download

import (
	"crypto/cipher"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"upload-view/backend/email"
	"upload-view/backend/encryption"
	"upload-view/backend/image"
	b_logger "upload-view/backend/logger"
	f_utils "upload-view/frontend/utils"
)

func EncryptStringsInStruct(HKID string, v reflect.Value, cb cipher.Block, nonce []byte) {
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)
		switch sf.Type.Kind() {
		case reflect.Struct:
			EncryptStringsInStruct(HKID, v.Field(i), cb, nonce)
		case reflect.String:
			field_str := v.Field(i).Interface().(string)
			if field_str == HKID {
				continue
			}
			encryptedmsg := encryption.EncryptDataByAESKey(cb, nonce, []byte(field_str))
			v.Field(i).SetString(string(encryptedmsg))
		case reflect.Slice:
			field_str := v.Field(i).Interface().([]byte)
			encryptedmsg := encryption.EncryptDataByAESKey(cb, nonce, field_str)
			v.Field(i).SetBytes(encryptedmsg)
		}
	}
}

func DownloadImage(r *http.Request) string {
	log := f_utils.ServLog
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
	jpgpath := filepath.Join("..", "jpg")
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
	return handler.Filename
}

func OpenImage(r *http.Request, src_img string) []byte {
	currenttime := b_logger.GetCurrentTime()
	logfile := r.FormValue("HKID") + "_" + currenttime + ".log"

	b_logger.InitLogger(logfile)

	jpgpath := filepath.Join("..", "jpg")
	mapping := image.JpgToMetaData(filepath.Join(jpgpath, src_img))

	return mapping.Image
}

func Register_DownloadPatient(r *http.Request) {
	log := f_utils.ServLog
	err := r.ParseForm()
	if err != nil {
		log.Println("Parsing MultiPart Form encounters error")
		log.Println(err)
		os.Exit(1)
	}
	details := f_utils.Patient{
		FullName:      []byte(r.FormValue("FullName")),
		Gender:        []byte(r.FormValue("Gender")),
		HKID:          []byte(r.FormValue("HKID")),
		DateOfBirth:   []byte(r.FormValue("DOB")),
		Race:          []byte(r.FormValue("Race")),
		PlaceOfBirth:  []byte(r.FormValue("POB")),
		Address:       []byte(r.FormValue("Address")),
		ContactNumber: []byte(r.FormValue("ContactNumber")),
		Email:         []byte(r.FormValue("Email")),
		MaritalStatus: []byte(r.FormValue("MaritalStatus")),
	}
	currenttime := b_logger.GetCurrentTime()
	basefile := r.FormValue("HKID") + "_" + currenttime
	logfile := basefile + ".log"

	b_logger.InitLogger(logfile)

	patientpath := filepath.Join("..", "patient")
	err = os.MkdirAll(patientpath, os.ModePerm)
	if err != nil {
		log.Println("Creating json directory for patient encounters error")
		log.Println(err)
	}

	destfile := basefile + "_Patient.json"

	var attachment []string
	//Check if Key Exist
	if !f_utils.Dev {
		log.Println("In DownloadPatient, Checking Key Exist ? ", encryption.CheckKeyExist(r.FormValue("HKID")))
		if encryption.CheckKeyExist(r.FormValue("HKID")) {
			aespath := filepath.Join("..", "key", r.FormValue("HKID")+"_AES.key")
			privatekeypath := filepath.Join("..", "key", r.FormValue("HKID")+"_Private.pem")
			encrypted_aeskey := encryption.ReadAESKey(aespath)
			privatekey := encryption.ReadPrivateKey(privatekeypath)

			aes_key := encryption.DecryptKeyByPrivateKey(privatekey, encrypted_aeskey)
			cb, nonce := encryption.GetCipherAndNonce(aes_key)
			EncryptStringsInStruct(r.FormValue("HKID"), reflect.ValueOf(&details).Elem(), cb, nonce)
		} else {
			private_key, public_key := encryption.InitatePubicPrivateKey()
			aes_key := encryption.GenerateAESKey()
			cb, nonce := encryption.GetCipherAndNonce(aes_key)
			encrypted_key := encryption.EncryptKeyByPublicKey(public_key, aes_key)
			Private_Key_Name, Public_Key_Name, encryption_key_name := encryption.CreateKeyFiles(public_key, private_key, encrypted_key, r.FormValue("HKID"))
			keypath := filepath.Join("..", "key")
			prikeypath := filepath.Join(keypath, Private_Key_Name)
			pubkeypath := filepath.Join(keypath, Public_Key_Name)
			encryptkeypath := filepath.Join(keypath, encryption_key_name)
			EncryptStringsInStruct(r.FormValue("HKID"), reflect.ValueOf(&details).Elem(), cb, nonce)
			attachment = append(attachment, prikeypath, pubkeypath, encryptkeypath)
		}
	}
	jsonByte, err := json.MarshalIndent(details, "", " ")
	if err != nil {
		log.Println("When creating blockchain record, JSON encoding encounters Error")
		log.Println(err)
		os.Exit(1)
	}
	err = ioutil.WriteFile(filepath.Join(patientpath, destfile), jsonByte, 0644)
	if err != nil {
		log.Println("When creating blockchain record, writing JSON encounters Error")
		log.Println(err)
		os.Exit(1)
	}
	recipient := r.FormValue("Email")
	patient_file_path := filepath.Join(patientpath, destfile)
	attachment = append(attachment, patient_file_path)
	email.SendEmailViaAWS(recipient, attachment...)
}

func Register_DownloadAllergy(r *http.Request) {
	log := f_utils.ServLog
	err := r.ParseForm()
	if err != nil {
		log.Println("Parsing MultiPart Form encounters error")
		log.Println(err)
		os.Exit(1)
	}
	details := f_utils.Allergy{
		AllergyIndicator: []byte(r.FormValue("Allergy")),
		AllergyTo:        []byte(r.FormValue("AllergyTo")),
	}
	currenttime := b_logger.GetCurrentTime()
	basefile := r.FormValue("HKID") + "_" + currenttime
	logfile := basefile + ".log"

	b_logger.InitLogger(logfile)

	patientpath := filepath.Join("..", "allergy")
	err = os.MkdirAll(patientpath, os.ModePerm)
	if err != nil {
		log.Println("Creating json directory for patient encounters error")
		log.Println(err)
	}

	destfile := basefile + "_Allergy.json"

	var attachment []string
	//Check if Key Exist
	if !f_utils.Dev {
		log.Println("In DownloadAllergy, Checking Key Exist ? ", encryption.CheckKeyExist(r.FormValue("HKID")))
		if encryption.CheckKeyExist(r.FormValue("HKID")) {
			aespath := filepath.Join("..", "key", r.FormValue("HKID")+"_AES.key")
			privatekeypath := filepath.Join("..", "key", r.FormValue("HKID")+"_Private.pem")
			encrypted_aeskey := encryption.ReadAESKey(aespath)
			privatekey := encryption.ReadPrivateKey(privatekeypath)

			aes_key := encryption.DecryptKeyByPrivateKey(privatekey, encrypted_aeskey)
			cb, nonce := encryption.GetCipherAndNonce(aes_key)
			EncryptStringsInStruct(r.FormValue("HKID"), reflect.ValueOf(&details).Elem(), cb, nonce)
		} else {
			private_key, public_key := encryption.InitatePubicPrivateKey()
			aes_key := encryption.GenerateAESKey()
			cb, nonce := encryption.GetCipherAndNonce(aes_key)
			encrypted_key := encryption.EncryptKeyByPublicKey(public_key, aes_key)
			Private_Key_Name, Public_Key_Name, encryption_key_name := encryption.CreateKeyFiles(public_key, private_key, encrypted_key, r.FormValue("HKID"))
			keypath := filepath.Join("..", "key")
			prikeypath := filepath.Join(keypath, Private_Key_Name)
			pubkeypath := filepath.Join(keypath, Public_Key_Name)
			encryptkeypath := filepath.Join(keypath, encryption_key_name)
			EncryptStringsInStruct(r.FormValue("HKID"), reflect.ValueOf(&details).Elem(), cb, nonce)
			attachment = append(attachment, prikeypath, pubkeypath, encryptkeypath)
		}
	}
	jsonByte, err := json.MarshalIndent(details, "", " ")
	if err != nil {
		log.Println("When creating blockchain record, JSON encoding encounters Error")
		log.Println(err)
		os.Exit(1)
	}
	err = ioutil.WriteFile(filepath.Join(patientpath, destfile), jsonByte, 0644)
	if err != nil {
		log.Println("When creating blockchain record, writing JSON encounters Error")
		log.Println(err)
		os.Exit(1)
	}
	recipient := r.FormValue("Email")
	patient_file_path := filepath.Join(patientpath, destfile)
	attachment = append(attachment, patient_file_path)
	email.SendEmailViaAWS(recipient, attachment...)
}

func Vaccine_DownloadVaccination(r *http.Request) {
	log := f_utils.ServLog
	err := r.ParseForm()
	if err != nil {
		log.Println("Parsing MultiPart Form encounters error")
		log.Println(err)
		os.Exit(1)
	}
	details := f_utils.COVIDVaccination{
		Date:  []byte(r.FormValue("DOV")),
		Venue: []byte(r.FormValue("Location")),
		Name:  []byte(r.FormValue("VaccineName")),
		Dose:  []byte(r.FormValue("Dose")),
		Complication: f_utils.Complication{
			Death:         []byte(r.FormValue("Death")),
			Stroke:        []byte(r.FormValue("Stroke")),
			HeartFailure:  []byte(r.FormValue("HeartFailure")),
			OtherSymptoms: []byte(r.FormValue("OtherSymptoms")),
		},
	}
	currenttime := b_logger.GetCurrentTime()
	basefile := r.FormValue("HKID") + "_" + currenttime
	logfile := basefile + ".log"

	b_logger.InitLogger(logfile)

	patientpath := filepath.Join("..", "vaccine")
	err = os.MkdirAll(patientpath, os.ModePerm)
	if err != nil {
		log.Println("Creating json directory for patient encounters error")
		log.Println(err)
	}

	destfile := basefile + "_vaccine.json"

	var attachment []string
	//Check if Key Exist
	if !f_utils.Dev {
		log.Println("In DownloadVaccine, Checking Key Exist ? ", encryption.CheckKeyExist(r.FormValue("HKID")))
		if encryption.CheckKeyExist(r.FormValue("HKID")) {
			aespath := filepath.Join("..", "key", r.FormValue("HKID")+"_AES.key")
			privatekeypath := filepath.Join("..", "key", r.FormValue("HKID")+"_Private.pem")
			encrypted_aeskey := encryption.ReadAESKey(aespath)
			privatekey := encryption.ReadPrivateKey(privatekeypath)

			aes_key := encryption.DecryptKeyByPrivateKey(privatekey, encrypted_aeskey)
			cb, nonce := encryption.GetCipherAndNonce(aes_key)
			EncryptStringsInStruct(r.FormValue("HKID"), reflect.ValueOf(&details).Elem(), cb, nonce)
		} else {
			private_key, public_key := encryption.InitatePubicPrivateKey()
			aes_key := encryption.GenerateAESKey()
			cb, nonce := encryption.GetCipherAndNonce(aes_key)
			encrypted_key := encryption.EncryptKeyByPublicKey(public_key, aes_key)
			Private_Key_Name, Public_Key_Name, encryption_key_name := encryption.CreateKeyFiles(public_key, private_key, encrypted_key, r.FormValue("HKID"))
			keypath := filepath.Join("..", "key")
			prikeypath := filepath.Join(keypath, Private_Key_Name)
			pubkeypath := filepath.Join(keypath, Public_Key_Name)
			encryptkeypath := filepath.Join(keypath, encryption_key_name)
			EncryptStringsInStruct(r.FormValue("HKID"), reflect.ValueOf(&details).Elem(), cb, nonce)
			attachment = append(attachment, prikeypath, pubkeypath, encryptkeypath)
		}
	}
	jsonByte, err := json.MarshalIndent(details, "", " ")
	if err != nil {
		log.Println("When creating blockchain record, JSON encoding encounters Error")
		log.Println(err)
		os.Exit(1)
	}
	err = ioutil.WriteFile(filepath.Join(patientpath, destfile), jsonByte, 0644)
	if err != nil {
		log.Println("When creating blockchain record, writing JSON encounters Error")
		log.Println(err)
		os.Exit(1)
	}
	recipient := r.FormValue("Email")
	patient_file_path := filepath.Join(patientpath, destfile)
	attachment = append(attachment, patient_file_path)
	email.SendEmailViaAWS(recipient, attachment...)
}

func Diet_DownloadBodyMeasurement(r *http.Request) {
	log := f_utils.ServLog
	err := r.ParseForm()
	if err != nil {
		log.Println("Parsing MultiPart Form encounters error")
		log.Println(err)
		os.Exit(1)
	}
	details := f_utils.BodyMeasurement{
		Height:       []byte(r.FormValue("Height")),
		Weight:       []byte(r.FormValue("Weight")),
		BMI:          []byte(r.FormValue("BMI")),
		BloodGlucose: []byte(r.FormValue("BloodGlucose")),
	}
	currenttime := b_logger.GetCurrentTime()
	basefile := r.FormValue("HKID") + "_" + currenttime
	logfile := basefile + ".log"

	b_logger.InitLogger(logfile)

	patientpath := filepath.Join("..", "bodymeasurement")
	err = os.MkdirAll(patientpath, os.ModePerm)
	if err != nil {
		log.Println("Creating json directory for patient encounters error")
		log.Println(err)
	}

	destfile := basefile + "_bodymeasurement.json"

	var attachment []string
	//Check if Key Exist
	if !f_utils.Dev {
		log.Println("In Download BodyMeasurement, Checking Key Exist ? ", encryption.CheckKeyExist(r.FormValue("HKID")))
		if encryption.CheckKeyExist(r.FormValue("HKID")) {
			aespath := filepath.Join("..", "key", r.FormValue("HKID")+"_AES.key")
			privatekeypath := filepath.Join("..", "key", r.FormValue("HKID")+"_Private.pem")
			encrypted_aeskey := encryption.ReadAESKey(aespath)
			privatekey := encryption.ReadPrivateKey(privatekeypath)

			aes_key := encryption.DecryptKeyByPrivateKey(privatekey, encrypted_aeskey)
			cb, nonce := encryption.GetCipherAndNonce(aes_key)
			EncryptStringsInStruct(r.FormValue("HKID"), reflect.ValueOf(&details).Elem(), cb, nonce)
		} else {
			private_key, public_key := encryption.InitatePubicPrivateKey()
			aes_key := encryption.GenerateAESKey()
			cb, nonce := encryption.GetCipherAndNonce(aes_key)
			encrypted_key := encryption.EncryptKeyByPublicKey(public_key, aes_key)
			Private_Key_Name, Public_Key_Name, encryption_key_name := encryption.CreateKeyFiles(public_key, private_key, encrypted_key, r.FormValue("HKID"))
			keypath := filepath.Join("..", "key")
			prikeypath := filepath.Join(keypath, Private_Key_Name)
			pubkeypath := filepath.Join(keypath, Public_Key_Name)
			encryptkeypath := filepath.Join(keypath, encryption_key_name)
			EncryptStringsInStruct(r.FormValue("HKID"), reflect.ValueOf(&details).Elem(), cb, nonce)
			attachment = append(attachment, prikeypath, pubkeypath, encryptkeypath)
		}
	}
	jsonByte, err := json.MarshalIndent(details, "", " ")
	if err != nil {
		log.Println("When creating blockchain record, JSON encoding encounters Error")
		log.Println(err)
		os.Exit(1)
	}
	err = ioutil.WriteFile(filepath.Join(patientpath, destfile), jsonByte, 0644)
	if err != nil {
		log.Println("When creating blockchain record, writing JSON encounters Error")
		log.Println(err)
		os.Exit(1)
	}
	recipient := r.FormValue("Email")
	patient_file_path := filepath.Join(patientpath, destfile)
	attachment = append(attachment, patient_file_path)
	email.SendEmailViaAWS(recipient, attachment...)
}

func Diet_DownloadMealPlan(r *http.Request) {
	log := f_utils.ServLog
	err := r.ParseForm()
	if err != nil {
		log.Println("Parsing MultiPart Form encounters error")
		log.Println(err)
		os.Exit(1)
	}
	details := f_utils.MealPlan{
		Breakfast: []byte(r.FormValue("Breakfast")),
		Lunch:     []byte(r.FormValue("Lunch")),
		Dinner:    []byte(r.FormValue("Dinner")),
		Snack:     []byte(r.FormValue("Snack")),
		Remark:    []byte(r.FormValue("Remark")),
	}
	currenttime := b_logger.GetCurrentTime()
	basefile := r.FormValue("HKID") + "_" + currenttime
	logfile := basefile + ".log"

	b_logger.InitLogger(logfile)

	patientpath := filepath.Join("..", "mealplan")
	err = os.MkdirAll(patientpath, os.ModePerm)
	if err != nil {
		log.Println("Creating json directory for patient encounters error")
		log.Println(err)
	}

	destfile := basefile + "_mealplan.json"

	var attachment []string
	//Check if Key Exist
	if !f_utils.Dev {
		log.Println("In DownloadMealPlan, Checking Key Exist ? ", encryption.CheckKeyExist(r.FormValue("HKID")))
		if encryption.CheckKeyExist(r.FormValue("HKID")) {
			aespath := filepath.Join("..", "key", r.FormValue("HKID")+"_AES.key")
			privatekeypath := filepath.Join("..", "key", r.FormValue("HKID")+"_Private.pem")
			encrypted_aeskey := encryption.ReadAESKey(aespath)
			privatekey := encryption.ReadPrivateKey(privatekeypath)

			aes_key := encryption.DecryptKeyByPrivateKey(privatekey, encrypted_aeskey)
			cb, nonce := encryption.GetCipherAndNonce(aes_key)
			EncryptStringsInStruct(r.FormValue("HKID"), reflect.ValueOf(&details).Elem(), cb, nonce)
		} else {
			private_key, public_key := encryption.InitatePubicPrivateKey()
			aes_key := encryption.GenerateAESKey()
			cb, nonce := encryption.GetCipherAndNonce(aes_key)
			encrypted_key := encryption.EncryptKeyByPublicKey(public_key, aes_key)
			Private_Key_Name, Public_Key_Name, encryption_key_name := encryption.CreateKeyFiles(public_key, private_key, encrypted_key, r.FormValue("HKID"))
			keypath := filepath.Join("..", "key")
			prikeypath := filepath.Join(keypath, Private_Key_Name)
			pubkeypath := filepath.Join(keypath, Public_Key_Name)
			encryptkeypath := filepath.Join(keypath, encryption_key_name)
			EncryptStringsInStruct(r.FormValue("HKID"), reflect.ValueOf(&details).Elem(), cb, nonce)
			attachment = append(attachment, prikeypath, pubkeypath, encryptkeypath)
		}
	}
	jsonByte, err := json.MarshalIndent(details, "", " ")
	if err != nil {
		log.Println("When creating blockchain record, JSON encoding encounters Error")
		log.Println(err)
		os.Exit(1)
	}
	err = ioutil.WriteFile(filepath.Join(patientpath, destfile), jsonByte, 0644)
	if err != nil {
		log.Println("When creating blockchain record, writing JSON encounters Error")
		log.Println(err)
		os.Exit(1)
	}
	recipient := r.FormValue("Email")
	patient_file_path := filepath.Join(patientpath, destfile)
	attachment = append(attachment, patient_file_path)
	email.SendEmailViaAWS(recipient, attachment...)
}

func AE_DownloadAssessment(r *http.Request) {
	log := f_utils.ServLog
	err := r.ParseForm()
	if err != nil {
		log.Println("Parsing MultiPart Form encounters error")
		log.Println(err)
		os.Exit(1)
	}
	details := f_utils.Assessment{
		AttendanceDate:  []byte(r.FormValue("AttendanceDate")),
		Triage:          []byte(r.FormValue("triage")),
		BloodPressure:   []byte(r.FormValue("BloodPressure")),
		BodyTemperature: []byte(r.FormValue("BodyTemp")),
		Pulse:           []byte(r.FormValue("Pulse")),
		Allergy: f_utils.Allergy{
			AllergyIndicator: []byte(r.FormValue("Allergy")),
			AllergyTo:        []byte(r.FormValue("AllergyTo")),
		},
		Diagnosis:           []byte(r.FormValue("Diagnosis")),
		AttendanceSpecialty: []byte(r.FormValue("AttendanceSpecialty")),
		ClinicalNote:        []byte(r.FormValue("ClinicalNotes")),
		FollowUpPlan: f_utils.FollowUpPlan{
			DrugPrescription:     []byte(r.FormValue("Drug")),
			Examination:          []byte(r.FormValue("Examination")),
			DischargeDestination: []byte(r.FormValue("DischargeDest")),
			AdmittedTo:           []byte(r.FormValue("AdmittedDest")),
		},
	}
	currenttime := b_logger.GetCurrentTime()
	basefile := r.FormValue("HKID") + "_" + currenttime
	logfile := basefile + ".log"

	b_logger.InitLogger(logfile)

	patientpath := filepath.Join("..", "assessment")
	err = os.MkdirAll(patientpath, os.ModePerm)
	if err != nil {
		log.Println("Creating json directory for patient encounters error")
		log.Println(err)
	}

	destfile := basefile + "_assessment.json"

	var attachment []string
	//Check if Key Exist
	if !f_utils.Dev {
		log.Println("In DownloadAssessment, Checking Key Exist ? ", encryption.CheckKeyExist(r.FormValue("HKID")))
		if encryption.CheckKeyExist(r.FormValue("HKID")) {
			aespath := filepath.Join("..", "key", r.FormValue("HKID")+"_AES.key")
			privatekeypath := filepath.Join("..", "key", r.FormValue("HKID")+"_Private.pem")
			encrypted_aeskey := encryption.ReadAESKey(aespath)
			privatekey := encryption.ReadPrivateKey(privatekeypath)

			aes_key := encryption.DecryptKeyByPrivateKey(privatekey, encrypted_aeskey)
			cb, nonce := encryption.GetCipherAndNonce(aes_key)
			EncryptStringsInStruct(r.FormValue("HKID"), reflect.ValueOf(&details).Elem(), cb, nonce)
		} else {
			private_key, public_key := encryption.InitatePubicPrivateKey()
			aes_key := encryption.GenerateAESKey()
			cb, nonce := encryption.GetCipherAndNonce(aes_key)
			encrypted_key := encryption.EncryptKeyByPublicKey(public_key, aes_key)
			Private_Key_Name, Public_Key_Name, encryption_key_name := encryption.CreateKeyFiles(public_key, private_key, encrypted_key, r.FormValue("HKID"))
			keypath := filepath.Join("..", "key")
			prikeypath := filepath.Join(keypath, Private_Key_Name)
			pubkeypath := filepath.Join(keypath, Public_Key_Name)
			encryptkeypath := filepath.Join(keypath, encryption_key_name)
			EncryptStringsInStruct(r.FormValue("HKID"), reflect.ValueOf(&details).Elem(), cb, nonce)
			attachment = append(attachment, prikeypath, pubkeypath, encryptkeypath)
		}
	}
	jsonByte, err := json.MarshalIndent(details, "", " ")
	if err != nil {
		log.Println("When creating blockchain record, JSON encoding encounters Error")
		log.Println(err)
		os.Exit(1)
	}
	err = ioutil.WriteFile(filepath.Join(patientpath, destfile), jsonByte, 0644)
	if err != nil {
		log.Println("When creating blockchain record, writing JSON encounters Error")
		log.Println(err)
		os.Exit(1)
	}
	recipient := r.FormValue("Email")
	patient_file_path := filepath.Join(patientpath, destfile)
	attachment = append(attachment, patient_file_path)
	email.SendEmailViaAWS(recipient, attachment...)
}

func ChestPain_DownloadChestPain(r *http.Request) {
	log := f_utils.ServLog
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		log.Println("Parsing MultiPart Form encounters error")
		log.Println(err)
		os.Exit(1)
	}
	img_file := DownloadImage(r)
	img_byte := OpenImage(r, img_file)
	details := f_utils.ChestPainReferral{
		ReferTo:        []byte(r.FormValue("ToDest")),
		Reason:         []byte(r.FormValue("Reason")),
		Duration:       []byte(r.FormValue("Duration")),
		Severity:       []byte(r.FormValue("Severity")),
		AnginaSymptoms: []byte(r.FormValue("TAS")),
		HistoryMI:      []byte(r.FormValue("MIHistory")),
		HistoryPTCA:    []byte(r.FormValue("PTCAHistory")),
		ChestXRayImage: img_byte,
	}
	currenttime := b_logger.GetCurrentTime()
	basefile := r.FormValue("HKID") + "_" + currenttime
	logfile := basefile + ".log"

	b_logger.InitLogger(logfile)

	patientpath := filepath.Join("..", "chestpain")
	err = os.MkdirAll(patientpath, os.ModePerm)
	if err != nil {
		log.Println("Creating json directory for patient encounters error")
		log.Println(err)
	}

	destfile := basefile + "_chestpain.json"

	var attachment []string
	//Check if Key Exist
	if !f_utils.Dev {
		log.Println("In DownloadChestPain, Checking Key Exist ? ", encryption.CheckKeyExist(r.FormValue("HKID")))
		if encryption.CheckKeyExist(r.FormValue("HKID")) {
			aespath := filepath.Join("..", "key", r.FormValue("HKID")+"_AES.key")
			privatekeypath := filepath.Join("..", "key", r.FormValue("HKID")+"_Private.pem")
			encrypted_aeskey := encryption.ReadAESKey(aespath)
			privatekey := encryption.ReadPrivateKey(privatekeypath)

			aes_key := encryption.DecryptKeyByPrivateKey(privatekey, encrypted_aeskey)
			cb, nonce := encryption.GetCipherAndNonce(aes_key)
			EncryptStringsInStruct(r.FormValue("HKID"), reflect.ValueOf(&details).Elem(), cb, nonce)
		} else {
			private_key, public_key := encryption.InitatePubicPrivateKey()
			aes_key := encryption.GenerateAESKey()
			cb, nonce := encryption.GetCipherAndNonce(aes_key)
			encrypted_key := encryption.EncryptKeyByPublicKey(public_key, aes_key)
			Private_Key_Name, Public_Key_Name, encryption_key_name := encryption.CreateKeyFiles(public_key, private_key, encrypted_key, r.FormValue("HKID"))
			keypath := filepath.Join("..", "key")
			prikeypath := filepath.Join(keypath, Private_Key_Name)
			pubkeypath := filepath.Join(keypath, Public_Key_Name)
			encryptkeypath := filepath.Join(keypath, encryption_key_name)
			EncryptStringsInStruct(r.FormValue("HKID"), reflect.ValueOf(&details).Elem(), cb, nonce)
			attachment = append(attachment, prikeypath, pubkeypath, encryptkeypath)
		}
	}
	jsonByte, err := json.MarshalIndent(details, "", " ")
	if err != nil {
		log.Println("When creating blockchain record, JSON encoding encounters Error")
		log.Println(err)
		os.Exit(1)
	}
	err = ioutil.WriteFile(filepath.Join(patientpath, destfile), jsonByte, 0644)
	if err != nil {
		log.Println("When creating blockchain record, writing JSON encounters Error")
		log.Println(err)
		os.Exit(1)
	}
	recipient := r.FormValue("Email")
	patient_file_path := filepath.Join(patientpath, destfile)
	attachment = append(attachment, patient_file_path)
	email.SendEmailViaAWS(recipient, attachment...)
}

func PCI_DownloadRiskFactor(r *http.Request) {
	log := f_utils.ServLog
	err := r.ParseForm()
	if err != nil {
		log.Println("Parsing MultiPart Form encounters error")
		log.Println(err)
		os.Exit(1)
	}
	details := f_utils.RiskFactor{
		PreviousPCI:       []byte(r.FormValue("PreviousPCI")),
		VascularDiseases:  []byte(r.FormValue("VascularDiseases")),
		DiabetesMelilites: []byte(r.FormValue("DiabetesMelilites")),
		Hypertension:      []byte(r.FormValue("Hypertension")),
		SmokingHistory:    []byte(r.FormValue("SmokeHistory")),
	}
	currenttime := b_logger.GetCurrentTime()
	basefile := r.FormValue("HKID") + "_" + currenttime
	logfile := basefile + ".log"

	b_logger.InitLogger(logfile)

	patientpath := filepath.Join("..", "riskfactor")
	err = os.MkdirAll(patientpath, os.ModePerm)
	if err != nil {
		log.Println("Creating json directory for patient encounters error")
		log.Println(err)
	}

	destfile := basefile + "_riskfactor.json"

	var attachment []string
	//Check if Key Exist
	if !f_utils.Dev {
		log.Println("In Download RiskFactor, Checking Key Exist ? ", encryption.CheckKeyExist(r.FormValue("HKID")))
		if encryption.CheckKeyExist(r.FormValue("HKID")) {
			aespath := filepath.Join("..", "key", r.FormValue("HKID")+"_AES.key")
			privatekeypath := filepath.Join("..", "key", r.FormValue("HKID")+"_Private.pem")
			encrypted_aeskey := encryption.ReadAESKey(aespath)
			privatekey := encryption.ReadPrivateKey(privatekeypath)

			aes_key := encryption.DecryptKeyByPrivateKey(privatekey, encrypted_aeskey)
			cb, nonce := encryption.GetCipherAndNonce(aes_key)
			EncryptStringsInStruct(r.FormValue("HKID"), reflect.ValueOf(&details).Elem(), cb, nonce)
		} else {
			private_key, public_key := encryption.InitatePubicPrivateKey()
			aes_key := encryption.GenerateAESKey()
			cb, nonce := encryption.GetCipherAndNonce(aes_key)
			encrypted_key := encryption.EncryptKeyByPublicKey(public_key, aes_key)
			Private_Key_Name, Public_Key_Name, encryption_key_name := encryption.CreateKeyFiles(public_key, private_key, encrypted_key, r.FormValue("HKID"))
			keypath := filepath.Join("..", "key")
			prikeypath := filepath.Join(keypath, Private_Key_Name)
			pubkeypath := filepath.Join(keypath, Public_Key_Name)
			encryptkeypath := filepath.Join(keypath, encryption_key_name)
			EncryptStringsInStruct(r.FormValue("HKID"), reflect.ValueOf(&details).Elem(), cb, nonce)
			attachment = append(attachment, prikeypath, pubkeypath, encryptkeypath)
		}
	}
	jsonByte, err := json.MarshalIndent(details, "", " ")
	if err != nil {
		log.Println("When creating blockchain record, JSON encoding encounters Error")
		log.Println(err)
		os.Exit(1)
	}
	err = ioutil.WriteFile(filepath.Join(patientpath, destfile), jsonByte, 0644)
	if err != nil {
		log.Println("When creating blockchain record, writing JSON encounters Error")
		log.Println(err)
		os.Exit(1)
	}
	recipient := r.FormValue("Email")
	patient_file_path := filepath.Join(patientpath, destfile)
	attachment = append(attachment, patient_file_path)
	email.SendEmailViaAWS(recipient, attachment...)
}

func PCI_DownloadProcedure(r *http.Request) {
	log := f_utils.ServLog
	err := r.ParseForm()
	if err != nil {
		log.Println("Parsing MultiPart Form encounters error")
		log.Println(err)
		os.Exit(1)
	}
	preInt := f_utils.PreIntervention{
		AnginaType:       []byte(r.FormValue("AnginaType")),
		HeartFailure:     []byte(r.FormValue("HeartFailure")),
		EjectionFraction: []byte(r.FormValue("EjectionFraction")),
	}

	com := f_utils.Complication{
		Death:         []byte(r.FormValue("Death")),
		Stroke:        []byte(r.FormValue("Stroke")),
		HeartFailure:  []byte(r.FormValue("HeartFailure")),
		OtherSymptoms: []byte(r.FormValue("OtherSymptoms")),
	}

	details := f_utils.Procedure{
		PreIntervention: preInt,
		ProcedureDate:   []byte(r.FormValue("ProcedureDate")),
		Urgency:         []byte(r.FormValue("Urgency")),
		Result:          []byte(r.FormValue("Result")),
		Device:          []byte(r.FormValue("Device")),
		Type:            []byte(r.FormValue("Type")),
		Complication:    com,
	}
	currenttime := b_logger.GetCurrentTime()
	basefile := r.FormValue("HKID") + "_" + currenttime
	logfile := basefile + ".log"

	b_logger.InitLogger(logfile)

	patientpath := filepath.Join("..", "procedure")
	err = os.MkdirAll(patientpath, os.ModePerm)
	if err != nil {
		log.Println("Creating json directory for patient encounters error")
		log.Println(err)
	}

	destfile := basefile + "_procedure.json"

	var attachment []string
	//Check if Key Exist
	if !f_utils.Dev {
		log.Println("In Download Procedure, Checking Key Exist ? ", encryption.CheckKeyExist(r.FormValue("HKID")))
		if encryption.CheckKeyExist(r.FormValue("HKID")) {
			aespath := filepath.Join("..", "key", r.FormValue("HKID")+"_AES.key")
			privatekeypath := filepath.Join("..", "key", r.FormValue("HKID")+"_Private.pem")
			encrypted_aeskey := encryption.ReadAESKey(aespath)
			privatekey := encryption.ReadPrivateKey(privatekeypath)

			aes_key := encryption.DecryptKeyByPrivateKey(privatekey, encrypted_aeskey)
			cb, nonce := encryption.GetCipherAndNonce(aes_key)
			EncryptStringsInStruct(r.FormValue("HKID"), reflect.ValueOf(&details).Elem(), cb, nonce)
		} else {
			private_key, public_key := encryption.InitatePubicPrivateKey()
			aes_key := encryption.GenerateAESKey()
			cb, nonce := encryption.GetCipherAndNonce(aes_key)
			encrypted_key := encryption.EncryptKeyByPublicKey(public_key, aes_key)
			Private_Key_Name, Public_Key_Name, encryption_key_name := encryption.CreateKeyFiles(public_key, private_key, encrypted_key, r.FormValue("HKID"))
			keypath := filepath.Join("..", "key")
			prikeypath := filepath.Join(keypath, Private_Key_Name)
			pubkeypath := filepath.Join(keypath, Public_Key_Name)
			encryptkeypath := filepath.Join(keypath, encryption_key_name)
			EncryptStringsInStruct(r.FormValue("HKID"), reflect.ValueOf(&details).Elem(), cb, nonce)
			attachment = append(attachment, prikeypath, pubkeypath, encryptkeypath)
		}
	}
	jsonByte, err := json.MarshalIndent(details, "", " ")
	if err != nil {
		log.Println("When creating blockchain record, JSON encoding encounters Error")
		log.Println(err)
		os.Exit(1)
	}
	err = ioutil.WriteFile(filepath.Join(patientpath, destfile), jsonByte, 0644)
	if err != nil {
		log.Println("When creating blockchain record, writing JSON encounters Error")
		log.Println(err)
		os.Exit(1)
	}
	recipient := r.FormValue("Email")
	patient_file_path := filepath.Join(patientpath, destfile)
	attachment = append(attachment, patient_file_path)
	email.SendEmailViaAWS(recipient, attachment...)
}

func DownloadNurse(r *http.Request) {
	log := f_utils.ServLog
	err := r.ParseForm()
	if err != nil {
		log.Println("Parsing MultiPart Form encounters error")
		log.Println(err)
		os.Exit(1)
	}
	details := f_utils.Nurse{
		FullName: []byte(r.FormValue("NurseName")),
		ID:       []byte(r.FormValue("NurseID")),
	}
	currenttime := b_logger.GetCurrentTime()
	basefile := r.FormValue("NurseID") + "_" + currenttime
	logfile := basefile + ".log"

	b_logger.InitLogger(logfile)

	patientpath := filepath.Join("..", "nurse")
	err = os.MkdirAll(patientpath, os.ModePerm)
	if err != nil {
		log.Println("Creating json directory for patient encounters error")
		log.Println(err)
	}

	destfile := basefile + "_nurse.json"

	jsonByte, err := json.MarshalIndent(details, "", " ")
	if err != nil {
		log.Println("When creating blockchain record, JSON encoding encounters Error")
		log.Println(err)
		os.Exit(1)
	}
	err = ioutil.WriteFile(filepath.Join(patientpath, destfile), jsonByte, 0644)
	if err != nil {
		log.Println("When creating blockchain record, writing JSON encounters Error")
		log.Println(err)
		os.Exit(1)
	}
	log.Println("Successfully created a nurse record")
}

func DownloadDoctor(r *http.Request) {
	log := f_utils.ServLog
	err := r.ParseForm()
	if err != nil {
		log.Println("Parsing MultiPart Form encounters error")
		log.Println(err)
		os.Exit(1)
	}
	details := f_utils.Doctor{
		FullName:  []byte(r.FormValue("DoctorName")),
		ID:        []byte(r.FormValue("DoctorID")),
		Hospital:  []byte(r.FormValue("Hospital")),
		Specialty: []byte(r.FormValue("Specialty")),
	}
	currenttime := b_logger.GetCurrentTime()
	basefile := r.FormValue("DoctorID") + "_" + currenttime
	logfile := basefile + ".log"

	b_logger.InitLogger(logfile)

	patientpath := filepath.Join("..", "doctor")
	err = os.MkdirAll(patientpath, os.ModePerm)
	if err != nil {
		log.Println("Creating json directory for patient encounters error")
		log.Println(err)
	}

	destfile := basefile + "_doctor.json"

	jsonByte, err := json.MarshalIndent(details, "", " ")
	if err != nil {
		log.Println("When creating blockchain record, JSON encoding encounters Error")
		log.Println(err)
		os.Exit(1)
	}
	err = ioutil.WriteFile(filepath.Join(patientpath, destfile), jsonByte, 0644)
	if err != nil {
		log.Println("When creating blockchain record, writing JSON encounters Error")
		log.Println(err)
		os.Exit(1)
	}
	log.Println("Successfully created a doctor record")
}
