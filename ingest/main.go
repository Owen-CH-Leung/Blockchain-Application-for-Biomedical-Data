package main

import (
	"fmt"
	"ingest/handler"
	f_utils "ingest/update_utils"
	"net/http"
)

func main() {
	f_utils.InitLogger()
	http.Handle("/_style/", http.StripPrefix("/_style/", http.FileServer(http.Dir("_style"))))
	http.HandleFunc("/", handler.Selectform)
	http.HandleFunc("/selectrecord", handler.Selectrecord)
	http.HandleFunc("/ae", handler.AE_form)
	http.HandleFunc("/chestpain", handler.Chestpain_form)
	http.HandleFunc("/pci", handler.PCI_form)
	http.HandleFunc("/vaccine", handler.Vaccine_form)
	http.HandleFunc("/registration", handler.Registration_form)
	http.HandleFunc("/diet", handler.Diet_form)
	http.HandleFunc("/error", handler.Error)
	fmt.Println("STARTING APPLICATION ON PORT 7777")
	http.ListenAndServe(":7777", nil)
}
