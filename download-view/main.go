package main

import (
	"download-view/frontend/handler"
	f_logger "download-view/frontend/utils"

	"fmt"
	// "html/template"
	"net/http"
)

func main() {
	f_logger.InitLogger()
	http.Handle("/_style/", http.StripPrefix("/_style/", http.FileServer(http.Dir("_style"))))
	http.HandleFunc("/", handler.SelectForm)
	http.HandleFunc("/nurse", handler.Nurse)
	http.HandleFunc("/doctor", handler.Doctor)
	http.HandleFunc("/patient", handler.Patient)
	http.HandleFunc("/researcher", handler.Researcher)
	http.HandleFunc("/insurance", handler.Insurance)
	http.HandleFunc("/emergency", handler.Emergency)
	http.HandleFunc("/error", handler.Error)
	fmt.Println("STARTING APPLICATION ON PORT 9999")
	http.ListenAndServe(":9999", nil)
}
