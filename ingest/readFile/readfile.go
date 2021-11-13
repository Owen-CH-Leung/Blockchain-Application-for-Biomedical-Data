package readFile

import (
	"encoding/json"
	"ingest/connectfabric"
	"ingest/utils"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func ReadJSON(path string) {
	log.Printf("--Looking at %v now--", path)
	keypath := filepath.Join("..", path)

	dir_files, err := ioutil.ReadDir(filepath.Clean(keypath))
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range dir_files {
		if !strings.Contains(file.Name(), "UPLOADED") {
			target_path := filepath.Clean(filepath.Join(keypath, file.Name()))
			jsonfile, err := os.Open(target_path)
			if err != nil {
				log.Println("Opening json file encounters Error")
				log.Println(err)
			}

			jsonByte, err := ioutil.ReadAll(jsonfile)
			if err != nil {
				log.Println("Read json as Byte encounters Error")
				log.Println(err)
				os.Exit(1)
			}

			switch path {
			case "allergy":
				json_map := utils.Allergy{}
				err := json.Unmarshal(jsonByte, &json_map)
				if err != nil {
					log.Println("Unmarshal JSON File stats encounters Error")
					log.Println(err)
					os.Exit(1)
				}
				hkid := strings.Split(file.Name(), "_")[0]
				connectfabric.SubmitToFabric(hkid, json_map)
			case "assessment":
				json_map := utils.Assessment{}
				err := json.Unmarshal(jsonByte, &json_map)
				if err != nil {
					log.Println("Unmarshal JSON File stats encounters Error")
					log.Println(err)
					os.Exit(1)
				}
				hkid := strings.Split(file.Name(), "_")[0]
				connectfabric.SubmitToFabric(hkid, json_map)
			case "bodymeasurement":
				json_map := utils.BodyMeasurement{}
				err := json.Unmarshal(jsonByte, &json_map)
				if err != nil {
					log.Println("Unmarshal JSON File stats encounters Error")
					log.Println(err)
					os.Exit(1)
				}
				hkid := strings.Split(file.Name(), "_")[0]
				connectfabric.SubmitToFabric(hkid, json_map)
			case "chestpain":
				json_map := utils.ChestPainReferral{}
				err := json.Unmarshal(jsonByte, &json_map)
				if err != nil {
					log.Println("Unmarshal JSON File stats encounters Error")
					log.Println(err)
					os.Exit(1)
				}
				hkid := strings.Split(file.Name(), "_")[0]
				connectfabric.SubmitToFabric(hkid, json_map)
			case "mealplan":
				json_map := utils.MealPlan{}
				err := json.Unmarshal(jsonByte, &json_map)
				if err != nil {
					log.Println("Unmarshal JSON File stats encounters Error")
					log.Println(err)
					os.Exit(1)
				}
				hkid := strings.Split(file.Name(), "_")[0]
				connectfabric.SubmitToFabric(hkid, json_map)
			case "patient":
				json_map := utils.Patient{}
				err := json.Unmarshal(jsonByte, &json_map)
				if err != nil {
					log.Println("Unmarshal JSON File stats encounters Error")
					log.Println(err)
					os.Exit(1)
				}
				hkid := strings.Split(file.Name(), "_")[0]
				connectfabric.SubmitToFabric(hkid, json_map)
			case "procedure":
				json_map := utils.Procedure{}
				err := json.Unmarshal(jsonByte, &json_map)
				if err != nil {
					log.Println("Unmarshal JSON File stats encounters Error")
					log.Println(err)
					os.Exit(1)
				}
				hkid := strings.Split(file.Name(), "_")[0]
				connectfabric.SubmitToFabric(hkid, json_map)
			case "riskfactor":
				json_map := utils.RiskFactor{}
				err := json.Unmarshal(jsonByte, &json_map)
				if err != nil {
					log.Println("Unmarshal JSON File stats encounters Error")
					log.Println(err)
					os.Exit(1)
				}
				hkid := strings.Split(file.Name(), "_")[0]
				connectfabric.SubmitToFabric(hkid, json_map)
			case "vaccine":
				json_map := utils.COVIDVaccination{}
				// log.Printf("Checking its type %T\n", json_map)
				// log.Println("In readfile.go, Before json.Unmarshal, Record Type name is ", reflect.TypeOf(json_map).Name())

				err := json.Unmarshal(jsonByte, &json_map)
				if err != nil {
					log.Println("Unmarshal JSON File stats encounters Error")
					log.Println(err)
					os.Exit(1)
				}
				hkid := strings.Split(file.Name(), "_")[0]
				// log.Println("In readfile.go, After json.Unmarshal, Record Type name is ", reflect.TypeOf(json_map).Name())
				// log.Println("Printing out the json_map ", json_map)
				connectfabric.SubmitToFabric(hkid, json_map)
			}
			jsonfile.Close()

			new_file := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name())) + "_UPLOADED.json"
			err = os.Rename(target_path, filepath.Clean(filepath.Join(keypath, new_file)))
			if err != nil {
				log.Println("Renaming file encounters error")
				log.Println(err)
			}
		} else {
			log.Printf("%s has been processed", file.Name())
		}
	}
}

func ReadJSONForUpdate(path string) {
	log.Printf("--Looking at %v now--", path)
	keypath := filepath.Join("..", path)

	dir_files, err := ioutil.ReadDir(filepath.Clean(keypath))
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range dir_files {
		if !strings.Contains(file.Name(), "UPLOADED") {
			target_path := filepath.Clean(filepath.Join(keypath, file.Name()))
			jsonfile, err := os.Open(target_path)
			if err != nil {
				log.Println("Opening json file encounters Error")
				log.Println(err)
			}

			jsonByte, err := ioutil.ReadAll(jsonfile)
			if err != nil {
				log.Println("Read json as Byte encounters Error")
				log.Println(err)
				os.Exit(1)
			}

			switch path {
			case "allergy_update":
				json_map := utils.Allergy{}
				err := json.Unmarshal(jsonByte, &json_map)
				if err != nil {
					log.Println("Unmarshal JSON File stats encounters Error")
					log.Println(err)
					os.Exit(1)
				}
				hkid := strings.Split(file.Name(), "_")[0]
				connectfabric.UpdateToFabric(hkid, json_map)
			case "assessment_update":
				json_map := utils.Assessment{}
				err := json.Unmarshal(jsonByte, &json_map)
				if err != nil {
					log.Println("Unmarshal JSON File stats encounters Error")
					log.Println(err)
					os.Exit(1)
				}
				hkid := strings.Split(file.Name(), "_")[0]
				connectfabric.UpdateToFabric(hkid, json_map)
			case "bodymeasurement_update":
				json_map := utils.BodyMeasurement{}
				err := json.Unmarshal(jsonByte, &json_map)
				if err != nil {
					log.Println("Unmarshal JSON File stats encounters Error")
					log.Println(err)
					os.Exit(1)
				}
				hkid := strings.Split(file.Name(), "_")[0]
				connectfabric.UpdateToFabric(hkid, json_map)
			case "chestpain_update":
				json_map := utils.ChestPainReferral{}
				err := json.Unmarshal(jsonByte, &json_map)
				if err != nil {
					log.Println("Unmarshal JSON File stats encounters Error")
					log.Println(err)
					os.Exit(1)
				}
				hkid := strings.Split(file.Name(), "_")[0]
				connectfabric.UpdateToFabric(hkid, json_map)
			case "mealplan_update":
				json_map := utils.MealPlan{}
				err := json.Unmarshal(jsonByte, &json_map)
				if err != nil {
					log.Println("Unmarshal JSON File stats encounters Error")
					log.Println(err)
					os.Exit(1)
				}
				hkid := strings.Split(file.Name(), "_")[0]
				connectfabric.UpdateToFabric(hkid, json_map)
			case "patient_update":
				json_map := utils.Patient{}
				err := json.Unmarshal(jsonByte, &json_map)
				if err != nil {
					log.Println("Unmarshal JSON File stats encounters Error")
					log.Println(err)
					os.Exit(1)
				}
				hkid := strings.Split(file.Name(), "_")[0]
				connectfabric.UpdateToFabric(hkid, json_map)
			case "procedure_update":
				json_map := utils.Procedure{}
				err := json.Unmarshal(jsonByte, &json_map)
				if err != nil {
					log.Println("Unmarshal JSON File stats encounters Error")
					log.Println(err)
					os.Exit(1)
				}
				hkid := strings.Split(file.Name(), "_")[0]
				connectfabric.UpdateToFabric(hkid, json_map)
			case "riskfactor_update":
				json_map := utils.RiskFactor{}
				err := json.Unmarshal(jsonByte, &json_map)
				if err != nil {
					log.Println("Unmarshal JSON File stats encounters Error")
					log.Println(err)
					os.Exit(1)
				}
				hkid := strings.Split(file.Name(), "_")[0]
				connectfabric.UpdateToFabric(hkid, json_map)
			case "vaccine_update":
				json_map := utils.COVIDVaccination{}
				// log.Printf("Checking its type %T\n", json_map)
				// log.Println("In readfile.go, Before json.Unmarshal, Record Type name is ", reflect.TypeOf(json_map).Name())

				err := json.Unmarshal(jsonByte, &json_map)
				if err != nil {
					log.Println("Unmarshal JSON File stats encounters Error")
					log.Println(err)
					os.Exit(1)
				}
				hkid := strings.Split(file.Name(), "_")[0]
				// log.Println("In readfile.go, After json.Unmarshal, Record Type name is ", reflect.TypeOf(json_map).Name())
				// log.Println("Printing out the json_map ", json_map)
				connectfabric.UpdateToFabric(hkid, json_map)
			}
			jsonfile.Close()

			new_file := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name())) + "_UPLOADED.json"
			err = os.Rename(target_path, filepath.Clean(filepath.Join(keypath, new_file)))
			if err != nil {
				log.Println("Renaming file encounters error")
				log.Println(err)
			}
		} else {
			log.Printf("%s has been processed", file.Name())
		}
	}
}
