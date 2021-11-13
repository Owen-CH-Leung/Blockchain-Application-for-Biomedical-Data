package handler

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"upload-view/frontend/download"
)

func Upload(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("nurse_registration_form.html"))
	if r.Method != http.MethodPost {
		tmpl.Execute(w, nil)
		return
	}
	//patientdetails, download_file := process.DownloadFile(r)
	//basefile, img_details := process.ReadEncryptEmail(download_file, patientdetails)
	//process.CreateBlockchainRecordJSON(basefile, img_details, patientdetails)
	//download.DownloadPatient(r)
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
		download.AE_DownloadAssessment(r)
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
		download.ChestPain_DownloadChestPain(r)
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
		download.PCI_DownloadRiskFactor(r)
		download.PCI_DownloadProcedure(r)
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
		download.Vaccine_DownloadVaccination(r)
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
		download.Register_DownloadPatient(r)
		download.Register_DownloadAllergy(r)
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
		download.Diet_DownloadBodyMeasurement(r)
		download.Diet_DownloadMealPlan(r)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/error", http.StatusSeeOther)
	}
	tmpl.Execute(w, struct{ Success bool }{true})
}

func Registration_Nurse(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("nurse_clinician_registration.html"))
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
	fmt.Println("In Nurse Registration Form")
	download.DownloadNurse(r)
	http.Redirect(w, r, "/", http.StatusSeeOther)
	tmpl.Execute(w, struct{ Success bool }{true})
}

func Registration_Doctor(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("doctor_registration_form.html"))
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
	fmt.Println("In Doctor Registration Form")
	download.DownloadDoctor(r)
	http.Redirect(w, r, "/", http.StatusSeeOther)
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

func Selectform(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("_master_form.html"))
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
	fmt.Println("Testing if this line will be printed")
	fmt.Println("Form Selected is ", r.Form["form"])
	switch r.Form["form"][0] {
	case "ae":
		http.Redirect(w, r, "/ae", http.StatusSeeOther)
	case "chestpain":
		http.Redirect(w, r, "/chestpain", http.StatusSeeOther)
	case "pci":
		http.Redirect(w, r, "/pci", http.StatusSeeOther)
	case "vaccine":
		http.Redirect(w, r, "/vaccine", http.StatusSeeOther)
	case "registration":
		http.Redirect(w, r, "/registration", http.StatusSeeOther)
	case "diet":
		http.Redirect(w, r, "/diet", http.StatusSeeOther)
	case "nurseregistration":
		http.Redirect(w, r, "/nurseregistration", http.StatusSeeOther)
	case "doctorregistration":
		http.Redirect(w, r, "/doctorregistration", http.StatusSeeOther)
	default:
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	tmpl.Execute(w, struct{ Success bool }{true})
}
