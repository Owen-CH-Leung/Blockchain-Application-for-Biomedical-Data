package detect

import (
	"ingest/readFile"
	"log"
	"sync"
)

func Trigger_Upload() {
	log.Println("=======STARTING THE APPLICATION=======")
	folder_path := []string{"allergy", "assessment", "bodymeasurement", "chestpain", "mealplan", "patient", "procedure",
		"riskfactor", "vaccine"}
	//readFile.ReadJSON("patient")
	// folder_path := []string{"allergy", "assessment", "bodymeasurement", "chestpain", "mealplan", "patient", "procedure",
	// 	"riskfactor"}
	// for _, path := range folder_path {
	// 	readFile.ReadJSON(path)
	// }
	var wg sync.WaitGroup
	for _, path := range folder_path {
		wg.Add(1)
		go func(arg string) {
			defer wg.Done()
			readFile.ReadJSON(arg)
		}(path)

	}
	wg.Wait()
}

func Trigger_Update() {
	log.Println("=======STARTING THE APPLICATION=======")
	folder_path := []string{"allergy_update", "assessment_update", "bodymeasurement_update", "chestpain_update", "mealplan_update",
		"patient_update", "procedure_update", "riskfactor_update", "vaccine_update"}
	//readFile.ReadJSON("patient")
	// folder_path := []string{"allergy", "assessment", "bodymeasurement", "chestpain", "mealplan", "patient", "procedure",
	// 	"riskfactor"}
	// for _, path := range folder_path {
	// 	readFile.ReadJSON(path)
	// }
	var wg sync.WaitGroup
	for _, path := range folder_path {
		wg.Add(1)
		go func(arg string) {
			defer wg.Done()
			readFile.ReadJSONForUpdate(arg)
		}(path)

	}
	wg.Wait()
}
