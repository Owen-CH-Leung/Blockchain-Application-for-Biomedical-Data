package handler

import (
	"fmt"
	"html/template"
	"ingest/detect"
	"ingest/update"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func Selectform(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("main.html"))
	if r.Method != http.MethodPost {
		tmpl.Execute(w, nil)
		return
	}
	fmt.Println("r.URL.PATH = ", r.URL.Path)
	err := r.ParseForm()
	if err != nil {
		log.Println("Parsing MultiPart Form encounters error")
		log.Println(err)
		os.Exit(1)
	}
	fmt.Println("Selected Action is ", r.FormValue("type"))
	switch r.FormValue("type") {
	case "create":
		detect.Trigger_Upload()
	case "update":
		http.Redirect(w, r, "/selectrecord", http.StatusSeeOther)
	case "push_update":
		detect.Trigger_Update()
	}
	tmpl.Execute(w, struct{ Success bool }{true})
}

func Selectrecord(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("recordtype.html"))
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
	fmt.Println("Selected Record Type is ", r.FormValue("recordtype"))
	switch r.FormValue("recordtype") {
	case "Patient/Allergy":
		http.Redirect(w, r, "/registration", http.StatusSeeOther)
	case "COVIDVaccine":
		http.Redirect(w, r, "/vaccine", http.StatusSeeOther)
	case "BodyMeasurement/MealPlan":
		http.Redirect(w, r, "/diet", http.StatusSeeOther)
	case "Assessment":
		http.Redirect(w, r, "/ae", http.StatusSeeOther)
	case "Procedure/RiskFactor":
		http.Redirect(w, r, "/pci", http.StatusSeeOther)
	case "ChestPainReferral":
		http.Redirect(w, r, "/chestpain", http.StatusSeeOther)
	}
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

func AE_form(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("doctor_ae_form.html"))
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
	fmt.Println("In A&E Form")
	if ClinicianExist("doctor", r.FormValue("DoctorID")) {
		update.AE_DownloadAssessment(r)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/error", http.StatusSeeOther)
	}
	tmpl.Execute(w, struct{ Success bool }{true})
}

func Chestpain_form(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("doctor_chestpain_form.html"))
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
	fmt.Println("In chestpain Form")
	if ClinicianExist("doctor", r.FormValue("DoctorID")) {
		update.ChestPain_DownloadChestPain(r)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/error", http.StatusSeeOther)
	}
	tmpl.Execute(w, struct{ Success bool }{true})
}

func PCI_form(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("doctor_pci_form.html"))
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
	fmt.Println("In PCI Form")
	if ClinicianExist("doctor", r.FormValue("DoctorID")) {
		update.PCI_DownloadRiskFactor(r)
		update.PCI_DownloadProcedure(r)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/error", http.StatusSeeOther)
	}
	tmpl.Execute(w, struct{ Success bool }{true})
}

func Vaccine_form(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("nurse_vaccine_form.html"))
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
	fmt.Println("In Vaccine Form")
	if ClinicianExist("nurse", r.FormValue("NurseID")) {
		update.Vaccine_DownloadVaccination(r)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/error", http.StatusSeeOther)
	}
	tmpl.Execute(w, struct{ Success bool }{true})
}

func Registration_form(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("nurse_registration_form.html"))
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
	if ClinicianExist("nurse", r.FormValue("NurseID")) {
		update.Register_DownloadPatient(r)
		update.Register_DownloadAllergy(r)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/error", http.StatusSeeOther)
	}
	tmpl.Execute(w, struct{ Success bool }{true})
}

func Diet_form(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("ah_diet_form.html"))
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
	fmt.Println("In Diet Form")
	if ClinicianExist("nurse", r.FormValue("NurseID")) {
		update.Diet_DownloadBodyMeasurement(r)
		update.Diet_DownloadMealPlan(r)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/error", http.StatusSeeOther)
	}
	tmpl.Execute(w, struct{ Success bool }{true})
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
