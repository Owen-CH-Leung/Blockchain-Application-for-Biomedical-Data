package connectfabric

import (
	"fmt"
	"ingest/utils"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"

	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

func SubmitToFabric(hkid string, record interface{}) {
	//record := ReadJSON()
	log.Println("============ application-golang starts ============")

	err := os.Setenv("DISCOVERY_AS_LOCALHOST", "true")
	if err != nil {
		log.Fatalf("Error setting DISCOVERY_AS_LOCALHOST environemnt variable: %v", err)
	}

	wallet, err := gateway.NewFileSystemWallet("wallet")
	if err != nil {
		log.Fatalf("Failed to create wallet: %v", err)
	}

	if !wallet.Exists("appUser") {
		err = populateWallet(wallet)
		if err != nil {
			log.Fatalf("Failed to populate wallet contents: %v", err)
		}
	}

	ccpPath := filepath.Join(
		"..",
		"fabric-network",
		"crypto-config",
		"peerOrganizations",
		"WalesHospital.com",
		"connection-org1.yaml",
	)

	gw, err := gateway.Connect(
		gateway.WithConfig(config.FromFile(filepath.Clean(ccpPath))),
		gateway.WithIdentity(wallet, "appUser"),
	)
	if err != nil {
		log.Fatalf("Failed to connect to gateway: %v", err)
	}
	defer gw.Close()

	network, err := gw.GetNetwork("hospitalhk")
	if err != nil {
		log.Fatalf("Failed to get network: %v", err)
	}

	contract := network.GetContract("hospitalCC")
	log.Println("In connect.go, Record Type name is ", reflect.TypeOf(record).Name())
	switch reflect.TypeOf(record).Name() {
	case "Patient":
		log.Println("--> Submit Transaction: CreatePatient")
		new_record := record.(utils.Patient)
		result, err := contract.SubmitTransaction("CreatePatient", hkid, new_record.FullName, new_record.Gender, new_record.HKID,
			new_record.DateOfBirth, new_record.Race, new_record.PlaceOfBirth, new_record.Address, new_record.ContactNumber,
			new_record.Email, new_record.MaritalStatus)

		if err != nil {
			log.Fatalf("When creating patient, Failed to Submit transaction: %v", err)
		}
		log.Println(string(result))

		log.Println("============ application-golang ends ============")

	case "COVIDVaccination":
		log.Println("--> Submit Transaction: Create COVIDVaccination")
		new_record := record.(utils.COVIDVaccination)
		result, err := contract.SubmitTransaction("CreateCOVIDVaccine", hkid, new_record.Date, new_record.Venue,
			new_record.Name, new_record.Dose, new_record.Complication.Death, new_record.Complication.Stroke,
			new_record.Complication.HeartFailure, new_record.Complication.OtherSymptoms)

		if err != nil {
			log.Fatalf("When creating COVIDVaccination, Failed to Submit transaction: %v", err)
		}
		log.Println(string(result))

		log.Println("============ application-golang ends ============")
	case "Allergy":
		log.Println("--> Submit Transaction: Create Allergy")
		new_record := record.(utils.Allergy)
		result, err := contract.SubmitTransaction("CreateAllergy", hkid, new_record.AllergyIndicator, new_record.AllergyTo)

		if err != nil {
			log.Fatalf("When creating COVIDVaccination, Failed to Submit transaction: %v", err)
		}
		log.Println(string(result))

		log.Println("============ application-golang ends ============")
	case "BodyMeasurement":
		log.Println("--> Submit Transaction: Create BodyMeasurement")
		new_record := record.(utils.BodyMeasurement)
		result, err := contract.SubmitTransaction("CreateBodyMeasurement", hkid, new_record.Height, new_record.Weight,
			new_record.BMI, new_record.BloodGlucose)

		if err != nil {
			log.Fatalf("When creating BodyMeasurement, Failed to Submit transaction: %v", err)
		}
		log.Println(string(result))

		log.Println("============ application-golang ends ============")
	case "MealPlan":
		log.Println("--> Submit Transaction: Create MealPlan")
		new_record := record.(utils.MealPlan)
		result, err := contract.SubmitTransaction("CreateMealPlan", hkid, new_record.Breakfast, new_record.Lunch,
			new_record.Dinner, new_record.Snack, new_record.Remark)

		if err != nil {
			log.Fatalf("When creating MealPlan, Failed to Submit transaction: %v", err)
		}
		log.Println(string(result))

		log.Println("============ application-golang ends ============")
	case "Assessment":
		log.Println("--> Submit Transaction: Create Assessment")
		new_record := record.(utils.Assessment)
		result, err := contract.SubmitTransaction("CreateAssessment", hkid, new_record.Allergy.AllergyIndicator, new_record.Allergy.AllergyTo,
			new_record.AttendanceDate, new_record.Triage, new_record.BloodPressure, new_record.BodyTemperature,
			new_record.Pulse, new_record.Diagnosis, new_record.AttendanceSpecialty, new_record.ClinicalNote,
			new_record.FollowUpPlan.DrugPrescription, new_record.FollowUpPlan.Examination, new_record.FollowUpPlan.DischargeDestination,
			new_record.FollowUpPlan.AdmittedTo)

		if err != nil {
			log.Fatalf("When creating Assessment, Failed to Submit transaction: %v", err)
		}
		log.Println(string(result))

		log.Println("============ application-golang ends ============")
	case "RiskFactor":
		log.Println("--> Submit Transaction: Create RiskFactor")
		new_record := record.(utils.RiskFactor)
		result, err := contract.SubmitTransaction("CreateRiskFactor", hkid, new_record.PreviousPCI, new_record.VascularDiseases,
			new_record.DiabetesMelilites, new_record.Hypertension, new_record.SmokingHistory)

		if err != nil {
			log.Fatalf("When creating RiskFactor, Failed to Submit transaction: %v", err)
		}
		log.Println(string(result))

		log.Println("============ application-golang ends ============")
	case "Procedure":
		log.Println("--> Submit Transaction: Create Procedure")
		new_record := record.(utils.Procedure)
		result, err := contract.SubmitTransaction("CreateProcedure", hkid, new_record.PreIntervention.AnginaType, new_record.PreIntervention.HeartFailure,
			new_record.PreIntervention.EjectionFraction, new_record.ProcedureDate, new_record.Urgency, new_record.Result,
			new_record.Device, new_record.Type, new_record.Complication.Death, new_record.Complication.Stroke, new_record.Complication.OtherSymptoms)

		if err != nil {
			log.Fatalf("When creating Procedure, Failed to Submit transaction: %v", err)
		}
		log.Println(string(result))
		log.Println("============ application-golang ends ============")
	case "ChestPainReferral":
		log.Println("--> Submit Transaction: Create ChestPainReferral")
		new_record := record.(utils.ChestPainReferral)
		result, err := contract.SubmitTransaction("CreateChestPainReferral", hkid, new_record.ReferTo, new_record.Reason,
			new_record.Duration, new_record.Severity, new_record.AnginaSymptoms, new_record.HistoryMI,
			new_record.HistoryPTCA, new_record.ChestXRayImage)

		if err != nil {
			log.Fatalf("When creating ChestPainReferral, Failed to Submit transaction: %v", err)
		}
		log.Println(string(result))

		log.Println("============ application-golang ends ============")
	}
}

func populateWallet(wallet *gateway.Wallet) error {
	log.Println("============ Populating wallet ============")
	credPath := filepath.Join(
		"..",
		"fabric-network",
		"crypto-config",
		"peerOrganizations",
		"WalesHospital.com",
		"users",
		"User1@WalesHospital.com",
		"msp",
	)
	log.Println("credPath: ", credPath)
	certPath := filepath.Join(credPath, "signcerts", "User1@WalesHospital.com-cert.pem")
	log.Println("certPath: ", certPath)
	// read the certificate pem
	cert, err := ioutil.ReadFile(filepath.Clean(certPath))
	if err != nil {
		log.Println("Error Reading CertPath:")
		log.Println(err)
		return err
	}

	keyDir := filepath.Join(credPath, "keystore")
	log.Println("keyDir: ", keyDir)
	// there's a single file in this dir containing the private key
	files, err := ioutil.ReadDir(keyDir)
	if err != nil {
		log.Println("Error reading KeyDir:")
		log.Println(err)
		return err
	}
	if len(files) != 1 {
		return fmt.Errorf("keystore folder should have contain one file")
	}
	keyPath := filepath.Join(keyDir, files[0].Name())
	log.Println("keyPath: ", keyPath)
	key, err := ioutil.ReadFile(filepath.Clean(keyPath))
	if err != nil {
		log.Println("Error reading keyPath:")
		log.Println(err)
		return err
	}

	identity := gateway.NewX509Identity("WalesHospitalMSP", string(cert), string(key))

	return wallet.Put("appUser", identity)
}

func UpdateToFabric(hkid string, record interface{}) {
	//record := ReadJSON()
	log.Println("============ application-golang starts ============")

	err := os.Setenv("DISCOVERY_AS_LOCALHOST", "true")
	if err != nil {
		log.Fatalf("Error setting DISCOVERY_AS_LOCALHOST environemnt variable: %v", err)
	}

	wallet, err := gateway.NewFileSystemWallet("wallet")
	if err != nil {
		log.Fatalf("Failed to create wallet: %v", err)
	}

	if !wallet.Exists("appUser") {
		err = populateWallet(wallet)
		if err != nil {
			log.Fatalf("Failed to populate wallet contents: %v", err)
		}
	}

	ccpPath := filepath.Join(
		"..",
		"fabric-network",
		"crypto-config",
		"peerOrganizations",
		"WalesHospital.com",
		"connection-org1.yaml",
	)

	gw, err := gateway.Connect(
		gateway.WithConfig(config.FromFile(filepath.Clean(ccpPath))),
		gateway.WithIdentity(wallet, "appUser"),
	)
	if err != nil {
		log.Fatalf("Failed to connect to gateway: %v", err)
	}
	defer gw.Close()

	network, err := gw.GetNetwork("hospitalhk")
	if err != nil {
		log.Fatalf("Failed to get network: %v", err)
	}

	contract := network.GetContract("hospitalCC")
	log.Println("In connect.go, Record Type name is ", reflect.TypeOf(record).Name())
	switch reflect.TypeOf(record).Name() {
	case "Patient":
		log.Println("--> Submit Transaction: UpdatePatient")
		new_record := record.(utils.Patient)
		result, err := contract.SubmitTransaction("UpdatePatient", hkid, new_record.FullName, new_record.Gender, new_record.HKID,
			new_record.DateOfBirth, new_record.Race, new_record.PlaceOfBirth, new_record.Address, new_record.ContactNumber,
			new_record.Email, new_record.MaritalStatus)

		if err != nil {
			log.Fatalf("When updating patient, Failed to Submit transaction: %v", err)
		}
		log.Println(string(result))

		log.Println("============ application-golang ends ============")

	case "COVIDVaccination":
		log.Println("--> Submit Transaction: Update COVIDVaccination")
		new_record := record.(utils.COVIDVaccination)
		result, err := contract.SubmitTransaction("UpdateCOVIDVaccine", hkid, new_record.Date, new_record.Venue,
			new_record.Name, new_record.Dose, new_record.Complication.Death, new_record.Complication.Stroke,
			new_record.Complication.HeartFailure, new_record.Complication.OtherSymptoms)

		if err != nil {
			log.Fatalf("When updating COVIDVaccination, Failed to Submit transaction: %v", err)
		}
		log.Println(string(result))

		log.Println("============ application-golang ends ============")
	case "Allergy":
		log.Println("--> Submit Transaction: Update Allergy")
		new_record := record.(utils.Allergy)
		result, err := contract.SubmitTransaction("UpdateAllergy", hkid, new_record.AllergyIndicator, new_record.AllergyTo)

		if err != nil {
			log.Fatalf("When updating COVIDVaccination, Failed to Submit transaction: %v", err)
		}
		log.Println(string(result))

		log.Println("============ application-golang ends ============")
	case "BodyMeasurement":
		log.Println("--> Submit Transaction: Update BodyMeasurement")
		new_record := record.(utils.BodyMeasurement)
		result, err := contract.SubmitTransaction("UpdateBodyMeasurement", hkid, new_record.Height, new_record.Weight,
			new_record.BMI, new_record.BloodGlucose)

		if err != nil {
			log.Fatalf("When updating BodyMeasurement, Failed to Submit transaction: %v", err)
		}
		log.Println(string(result))

		log.Println("============ application-golang ends ============")
	case "MealPlan":
		log.Println("--> Submit Transaction: Update MealPlan")
		new_record := record.(utils.MealPlan)
		result, err := contract.SubmitTransaction("UpdateMealPlan", hkid, new_record.Breakfast, new_record.Lunch,
			new_record.Dinner, new_record.Snack, new_record.Remark)

		if err != nil {
			log.Fatalf("When updating MealPlan, Failed to Submit transaction: %v", err)
		}
		log.Println(string(result))

		log.Println("============ application-golang ends ============")
	case "Assessment":
		log.Println("--> Submit Transaction: Update Assessment")
		new_record := record.(utils.Assessment)
		result, err := contract.SubmitTransaction("UpdateAssessment", hkid, new_record.Allergy.AllergyIndicator, new_record.Allergy.AllergyTo,
			new_record.AttendanceDate, new_record.Triage, new_record.BloodPressure, new_record.BodyTemperature,
			new_record.Pulse, new_record.Diagnosis, new_record.AttendanceSpecialty, new_record.ClinicalNote,
			new_record.FollowUpPlan.DrugPrescription, new_record.FollowUpPlan.Examination, new_record.FollowUpPlan.DischargeDestination,
			new_record.FollowUpPlan.AdmittedTo)

		if err != nil {
			log.Fatalf("When updating Assessment, Failed to Submit transaction: %v", err)
		}
		log.Println(string(result))

		log.Println("============ application-golang ends ============")
	case "RiskFactor":
		log.Println("--> Submit Transaction: Update RiskFactor")
		new_record := record.(utils.RiskFactor)
		result, err := contract.SubmitTransaction("UpdateRiskFactor", hkid, new_record.PreviousPCI, new_record.VascularDiseases,
			new_record.DiabetesMelilites, new_record.Hypertension, new_record.SmokingHistory)

		if err != nil {
			log.Fatalf("When updating RiskFactor, Failed to Submit transaction: %v", err)
		}
		log.Println(string(result))

		log.Println("============ application-golang ends ============")
	case "Procedure":
		log.Println("--> Submit Transaction: Update Procedure")
		new_record := record.(utils.Procedure)
		result, err := contract.SubmitTransaction("UpdateProcedure", hkid, new_record.PreIntervention.AnginaType, new_record.PreIntervention.HeartFailure,
			new_record.PreIntervention.EjectionFraction, new_record.ProcedureDate, new_record.Urgency, new_record.Result,
			new_record.Device, new_record.Type, new_record.Complication.Death, new_record.Complication.Stroke, new_record.Complication.OtherSymptoms)

		if err != nil {
			log.Fatalf("When updating Procedure, Failed to Submit transaction: %v", err)
		}
		log.Println(string(result))
		log.Println("============ application-golang ends ============")
	case "ChestPainReferral":
		log.Println("--> Submit Transaction: Update ChestPainReferral")
		new_record := record.(utils.ChestPainReferral)
		result, err := contract.SubmitTransaction("UpdateChestPainReferral", hkid, new_record.ReferTo, new_record.Reason,
			new_record.Duration, new_record.Severity, new_record.AnginaSymptoms, new_record.HistoryMI,
			new_record.HistoryPTCA, new_record.ChestXRayImage)

		if err != nil {
			log.Fatalf("When updating ChestPainReferral, Failed to Submit transaction: %v", err)
		}
		log.Println(string(result))

		log.Println("============ application-golang ends ============")
	}
}
