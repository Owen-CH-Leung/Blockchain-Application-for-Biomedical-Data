package main

import (
	"upload-view/frontend/handler"
	f_logger "upload-view/frontend/utils"

	"fmt"
	// "html/template"
	"net/http"
)

func main() {
	f_logger.InitLogger()
	http.Handle("/_style/", http.StripPrefix("/_style/", http.FileServer(http.Dir("_style"))))
	http.HandleFunc("/", handler.Selectform)
	http.HandleFunc("/ae", handler.AE_form)
	http.HandleFunc("/chestpain", handler.Chestpain_form)
	http.HandleFunc("/pci", handler.PCI_form)
	http.HandleFunc("/vaccine", handler.Vaccine_form)
	http.HandleFunc("/registration", handler.Registration_form)
	http.HandleFunc("/diet", handler.Diet_form)
	http.HandleFunc("/nurseregistration", handler.Registration_Nurse)
	http.HandleFunc("/doctorregistration", handler.Registration_Doctor)
	http.HandleFunc("/error", handler.Error)
	fmt.Println("STARTING APPLICATION ON PORT 8888")
	http.ListenAndServe(":8888", nil)
}
