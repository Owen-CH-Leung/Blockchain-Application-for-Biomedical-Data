package DASC7600

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

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
	fmt.Println("HKID is ", r.FormValue("HKID"))
	http.Redirect(w, r, "/", http.StatusSeeOther)
	tmpl.Execute(w, struct{ Success bool }{true})
}

func Chestpain_form(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("doctor_chestpain_form.html"))
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
	fmt.Println("In chestpain Form")
	fmt.Println("HKID is ", r.FormValue("HKID"))
	http.Redirect(w, r, "/", http.StatusSeeOther)
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
	fmt.Println("HKID is ", r.FormValue("HKID"))
	http.Redirect(w, r, "/", http.StatusSeeOther)
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
	fmt.Println("HKID is ", r.FormValue("HKID"))
	http.Redirect(w, r, "/", http.StatusSeeOther)
	tmpl.Execute(w, struct{ Success bool }{true})
}

func Registration_form(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("nurse_registration_form.html"))
	if r.Method != http.MethodPost {
		tmpl.Execute(w, nil)
		return
	}
	printout(r, "HKID")
	printout(r, "Allergy")
	printout(r, "AllergyTo")
	// err := r.ParseForm()
	// if err != nil {
	// 	log.Println("Parsing MultiPart Form encounters error")
	// 	log.Println(err)
	// 	os.Exit(1)
	// }
	// fmt.Println("In Registration Form")
	// fmt.Println("HKID is ", r.FormValue("HKID"))
	http.Redirect(w, r, "/", http.StatusSeeOther)
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
	fmt.Println("HKID is ", r.FormValue("HKID"))
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
	default:
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	tmpl.Execute(w, struct{ Success bool }{true})
}

func printout(r *http.Request, field string) {
	r.ParseForm()
	fmt.Println("The field value is ", r.FormValue(field))
}

func main() {
	http.HandleFunc("/", Selectform)
	http.HandleFunc("/ae", AE_form)
	http.HandleFunc("/chestpain", Chestpain_form)
	http.HandleFunc("/pci", PCI_form)
	http.HandleFunc("/vaccine", Vaccine_form)
	http.HandleFunc("/registration", Registration_form)
	http.HandleFunc("/diet", Diet_form)
	fmt.Println("STARTING APPLICATION ON PORT 8888")
	http.ListenAndServe(":8888", nil)
}
