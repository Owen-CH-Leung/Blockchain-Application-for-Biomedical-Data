package handler

import (
	aws_email "download-view/backend/email"
	"download-view/frontend/process"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func SelectForm(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("_master_form.html"))
	if r.Method != http.MethodPost {
		tmpl.Execute(w, nil)
		return
	}
	switch r.FormValue("role") {
	case "nurse":
		http.Redirect(w, r, "/nurse", http.StatusSeeOther)
	case "doctor":
		http.Redirect(w, r, "/doctor", http.StatusSeeOther)
	case "patient":
		http.Redirect(w, r, "/patient", http.StatusSeeOther)
	case "researcher":
		http.Redirect(w, r, "/researcher", http.StatusSeeOther)
	case "insurance":
		http.Redirect(w, r, "/insurance", http.StatusSeeOther)
	case "emergency":
		http.Redirect(w, r, "/emergency", http.StatusSeeOther)
	}
	fmt.Println("Form Value is ", r.FormValue("role"))
	// a function that takes AESkeypath and privatekey path, return the key content
	// a function that takes hkid as arg, decrypt the msg
	// basefile, img_details := process.ReadEncryptEmail(download_file, patientdetails)
	// process.CreateBlockchainRecordJSON(basefile, img_details, patientdetails)
	tmpl.Execute(w, struct{ Success bool }{true})
}

func Patient(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("patient.html"))
	if r.Method != http.MethodPost {
		tmpl.Execute(w, nil)
		return
	}

	email, hkid, privatekeypath := process.DownloadKeyFile(r)
	log.Println("Private Key Path is :", privatekeypath)
	aespath := process.GetAESKeyFiles(hkid)

	attachment := []string{}
	log.Println("The Selected Choices are :", r.Form["choice"])
	log.Println("The Selected data type is :", r.FormValue("type"))
	switch r.FormValue("type") {
	case "latest":
		for _, choice := range r.Form["choice"] {
			log.Println("Handling : ", choice)
			file := process.DecryptAndSaveAsJson(hkid, choice, aespath, privatekeypath)
			attachment = append(attachment, file...)
		}
	case "history":
		for _, choice := range r.Form["choice"] {
			log.Println("Handling : ", choice)
			file := process.DecryptAllAndSaveAsJson(hkid, choice, aespath, privatekeypath)
			attachment = append(attachment, file...)
		}
	}

	aws_email.SendEmailViaAWS(email, attachment...)
	http.Redirect(w, r, "/", http.StatusSeeOther)
	tmpl.Execute(w, struct{ Success bool }{true})
}

func Nurse(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("nurse.html"))
	if r.Method != http.MethodPost {
		tmpl.Execute(w, nil)
		return
	}
	if ClinicianExist("nurse", r.FormValue("NurseID")) {
		email, hkid, privatekeypath := process.DownloadKeyFile(r)
		log.Println("Private Key Path is :", privatekeypath)
		aespath := process.GetAESKeyFiles(hkid)

		attachment := []string{}
		log.Println("The Selected Choices are :", r.Form["choice"])
		log.Println("The Selected data type is :", r.FormValue("type"))
		switch r.FormValue("type") {
		case "latest":
			for _, choice := range r.Form["choice"] {
				log.Println("Handling : ", choice)
				file := process.DecryptAndSaveAsJson(hkid, choice, aespath, privatekeypath)
				attachment = append(attachment, file...)
			}
		case "history":
			for _, choice := range r.Form["choice"] {
				log.Println("Handling : ", choice)
				file := process.DecryptAllAndSaveAsJson(hkid, choice, aespath, privatekeypath)
				attachment = append(attachment, file...)
			}
		}

		aws_email.SendEmailViaAWS(email, attachment...)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/error", http.StatusSeeOther)
	}

	tmpl.Execute(w, struct{ Success bool }{true})
}

func Doctor(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("doctor.html"))
	if r.Method != http.MethodPost {
		tmpl.Execute(w, nil)
		return
	}
	if ClinicianExist("doctor", r.FormValue("DoctorID")) {
		email, hkid, privatekeypath := process.DownloadKeyFile(r)
		log.Println("Private Key Path is :", privatekeypath)
		aespath := process.GetAESKeyFiles(hkid)

		attachment := []string{}
		log.Println("The Selected Choices are :", r.Form["choice"])
		log.Println("The Selected data type is :", r.FormValue("type"))
		switch r.FormValue("type") {
		case "latest":
			for _, choice := range r.Form["choice"] {
				log.Println("Handling : ", choice)
				file := process.DecryptAndSaveAsJson(hkid, choice, aespath, privatekeypath)
				attachment = append(attachment, file...)
			}
		case "history":
			for _, choice := range r.Form["choice"] {
				log.Println("Handling : ", choice)
				file := process.DecryptAllAndSaveAsJson(hkid, choice, aespath, privatekeypath)
				attachment = append(attachment, file...)
			}
		}

		aws_email.SendEmailViaAWS(email, attachment...)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/error", http.StatusSeeOther)
	}

	tmpl.Execute(w, struct{ Success bool }{true})
}

func Insurance(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("insurance.html"))
	if r.Method != http.MethodPost {
		tmpl.Execute(w, nil)
		return
	}
	email, hkid, privatekeypath := process.DownloadKeyFile(r)
	log.Println("Private Key Path is :", privatekeypath)
	aespath := process.GetAESKeyFiles(hkid)

	attachment := []string{}
	log.Println("The Selected Choices are :", r.Form["choice"])
	log.Println("The Selected data type is :", r.FormValue("type"))
	switch r.FormValue("type") {
	case "latest":
		for _, choice := range r.Form["choice"] {
			log.Println("Handling : ", choice)
			file := process.DecryptAndSaveAsJson(hkid, choice, aespath, privatekeypath)
			attachment = append(attachment, file...)
		}
	case "history":
		for _, choice := range r.Form["choice"] {
			log.Println("Handling : ", choice)
			file := process.DecryptAllAndSaveAsJson(hkid, choice, aespath, privatekeypath)
			attachment = append(attachment, file...)
		}
	}

	aws_email.SendEmailViaAWS(email, attachment...)
	http.Redirect(w, r, "/", http.StatusSeeOther)
	tmpl.Execute(w, struct{ Success bool }{true})
}

func ClinicianExist(kind string, ID string) bool {
	switch kind {
	case "doctor":
		path := filepath.Join("..", "doctor")
		dir_files, err := ioutil.ReadDir(filepath.Clean(path))
		if err != nil {
			log.Fatal(err)
		}
		var exist bool = false
		for _, file := range dir_files {
			docID := strings.Split(file.Name(), "_")[0]
			log.Println("DocID in doctor folder is ", docID)
			log.Println("ID received is ", ID)
			if ID == docID {
				exist = true
				return exist
			}
		}
		return false
	case "nurse":
		path := filepath.Join("..", "nurse")
		dir_files, err := ioutil.ReadDir(filepath.Clean(path))
		if err != nil {
			log.Fatal(err)
		}
		var exist bool = false
		for _, file := range dir_files {
			nurseID := strings.Split(file.Name(), "_")[0]
			log.Println("nurseID in doctor folder is ", nurseID)
			log.Println("ID received is ", ID)
			if ID == nurseID {
				exist = true
				return exist
			}
		}
		return false
	default:
		ClinicianExist("nurse", ID)
	}
	return false
}

func Error(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("error.html"))
	if r.Method != http.MethodPost {
		tmpl.Execute(w, nil)
		return
	}
	err := r.ParseForm()
	if err != nil {
		log.Println("Parsing MultiPart Form encounters error")
		log.Println(err)
		os.Exit(1)
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
	tmpl.Execute(w, struct{ Success bool }{true})
}

func Researcher(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("researcher.html"))
	if r.Method != http.MethodPost {
		tmpl.Execute(w, nil)
		return
	}
	patient, vaccine, chestpain := process.Get_Patient_Vaccine_ChestPain_Info(r)
	details := Aggregate{
		TotalPatient:   patient,
		TotalVaccine:   vaccine,
		TotalChestPain: chestpain,
	}
	jsonByte, err := json.MarshalIndent(details, "", " ")
	if err != nil {
		log.Println("When creating blockchain record, JSON encoding encounters Error")
		log.Println(err)
		os.Exit(1)
	}
	researchpath := filepath.Join("..", "researcher_sent")
	err = os.MkdirAll(researchpath, os.ModePerm)
	if err != nil {
		log.Println("Creating json directory for patient encounters error")
		log.Println(err)
	}

	destfile := "Aggregate_Masked_List.json"
	err = ioutil.WriteFile(filepath.Join(researchpath, destfile), jsonByte, 0644)
	if err != nil {
		log.Println("When creating blockchain record, writing JSON encounters Error")
		log.Println(err)
		os.Exit(1)
	}
	recipient := r.FormValue("Email")
	file_path := filepath.Join(researchpath, destfile)
	attachment := []string{file_path}
	aws_email.SendEmailViaAWS(recipient, attachment...)
	http.Redirect(w, r, "/", http.StatusSeeOther)
	tmpl.Execute(w, struct{ Success bool }{true})
}

type Aggregate struct {
	TotalPatient   []string
	TotalVaccine   []string
	TotalChestPain []string
}

func Emergency(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("emergency.html"))
	if r.Method != http.MethodPost {
		tmpl.Execute(w, nil)
		return
	}
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		log.Println("Parsing MultiPart Form encounters error")
		log.Println(err)
		os.Exit(1)
	}
	privatekeypath := process.GetPrivateKeyFiles(r.FormValue("HKID"))
	aespath := process.GetAESKeyFiles(r.FormValue("HKID"))
	attachment := []string{}
	email := r.FormValue("Email")
	if ClinicianExist("doctor", r.FormValue("DoctorID")) {
		arr_choices := []string{"patient", "allergy", "assessment", "bodymeasurement", "chestpainreferral", "vaccine", "mealplan", "procedure", "riskfactor"}
		for _, choice := range arr_choices {
			log.Println("Handling : ", choice)
			file := process.DecryptAllAndSaveAsJson(r.FormValue("HKID"), choice, aespath, privatekeypath)
			attachment = append(attachment, file...)
		}
	} else {
		http.Redirect(w, r, "/error", http.StatusSeeOther)
	}

	aws_email.SendEmailViaAWS(email, attachment...)
	http.Redirect(w, r, "/", http.StatusSeeOther)
	tmpl.Execute(w, struct{ Success bool }{true})
}
