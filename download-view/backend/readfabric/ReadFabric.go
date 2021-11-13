package readfabric

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	//"time"
	"download-view/frontend/utils"

	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

func ReadLatest(hkid string, info string) (interface{}, error) {
	contract := GetContract(info)

	funcname := "Read" + info
	result, err := contract.SubmitTransaction(funcname, hkid)
	if err != nil {
		log.Println("Unable to retrieve transaction from fabric")
		log.Println(err)
		os.Exit(1)
	}
	log.Println("Successfully retrieve record by calling API : ", funcname)

	switch {
	case info == "patient":
		var record utils.Patient
		err = json.Unmarshal(result, &record)
		if err != nil {
			fmt.Println("Unmarshalling Json enocunters error")
			fmt.Println(err)
			os.Exit(1)
		}
		log.Println("Successfully Unmarshal Patient")
		return &record, nil
	case info == "vaccine":
		var record utils.COVIDVaccination
		err = json.Unmarshal(result, &record)
		if err != nil {
			fmt.Println("Unmarshalling Json enocunters error")
			fmt.Println(err)
			os.Exit(1)
		}

		return &record, nil
	case info == "allergy":
		var record utils.Allergy
		err = json.Unmarshal(result, &record)
		if err != nil {
			fmt.Println("Unmarshalling Json enocunters error")
			fmt.Println(err)
			os.Exit(1)
		}

		return &record, nil
	case info == "bodymeasurement":
		var record utils.BodyMeasurement
		err = json.Unmarshal(result, &record)
		if err != nil {
			fmt.Println("Unmarshalling Json enocunters error")
			fmt.Println(err)
			os.Exit(1)
		}
		return &record, nil
	case info == "mealplan":
		var record utils.MealPlan
		err = json.Unmarshal(result, &record)
		if err != nil {
			fmt.Println("Unmarshalling Json enocunters error")
			fmt.Println(err)
			os.Exit(1)
		}
		return &record, nil
	case info == "assessment":
		var record utils.Assessment
		err = json.Unmarshal(result, &record)
		if err != nil {
			fmt.Println("Unmarshalling Json enocunters error")
			fmt.Println(err)
			os.Exit(1)
		}
		return &record, nil
	case info == "riskfactor":
		var record utils.RiskFactor
		err = json.Unmarshal(result, &record)
		if err != nil {
			fmt.Println("Unmarshalling Json enocunters error")
			fmt.Println(err)
			os.Exit(1)
		}
		return &record, nil
	case info == "procedure":
		var record utils.Procedure
		err = json.Unmarshal(result, &record)
		if err != nil {
			fmt.Println("Unmarshalling Json enocunters error")
			fmt.Println(err)
			os.Exit(1)
		}
		return &record, nil
	case info == "chestpainreferral":
		var record utils.ChestPainReferral
		err = json.Unmarshal(result, &record)
		if err != nil {
			fmt.Println("Unmarshalling Json enocunters error")
			fmt.Println(err)
			os.Exit(1)
		}
		return &record, nil
	default:
		var record utils.Patient
		err = json.Unmarshal(result, &record)
		if err != nil {
			fmt.Println("Unmarshalling Json enocunters error")
			fmt.Println(err)
			os.Exit(1)
		}
		return &record, nil
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

func ReadPatient_All(hkid string) (*[]utils.Patient, error) {
	contract := GetContract("patient")
	funcname := "Readpatient_All"
	result, err := contract.SubmitTransaction(funcname, hkid)
	if err != nil {
		log.Println("Unable to retrieve transaction from fabric")
		log.Println(err)
		os.Exit(1)
	}
	log.Println("Successfully retrieve record by calling API : ", funcname)
	var record []utils.Patient
	err = json.Unmarshal(result, &record)
	if err != nil {
		fmt.Println("Unmarshalling Json enocunters error")
		fmt.Println(err)
		os.Exit(1)
	}
	log.Println("Successfully Unmarshal Patient")
	return &record, nil
}

func ReadVaccine_All(hkid string) (*[]utils.COVIDVaccination, error) {
	contract := GetContract("vaccine")
	funcname := "Readvaccine_All"
	result, err := contract.SubmitTransaction(funcname, hkid)
	if err != nil {
		log.Println("Unable to retrieve transaction from fabric")
		log.Println(err)
		os.Exit(1)
	}
	log.Println("Successfully retrieve record by calling API : ", funcname)
	var record []utils.COVIDVaccination
	err = json.Unmarshal(result, &record)
	if err != nil {
		fmt.Println("Unmarshalling Json enocunters error")
		fmt.Println(err)
		os.Exit(1)
	}
	log.Println("Successfully Unmarshal Patient")
	return &record, nil
}

func ReadAllergy_All(hkid string) (*[]utils.Allergy, error) {
	contract := GetContract("allergy")
	funcname := "Readallergy_All"
	result, err := contract.SubmitTransaction(funcname, hkid)
	if err != nil {
		log.Println("Unable to retrieve transaction from fabric")
		log.Println(err)
		os.Exit(1)
	}
	log.Println("Successfully retrieve record by calling API : ", funcname)
	var record []utils.Allergy
	err = json.Unmarshal(result, &record)
	if err != nil {
		fmt.Println("Unmarshalling Json enocunters error")
		fmt.Println(err)
		os.Exit(1)
	}
	log.Println("Successfully Unmarshal Patient")
	return &record, nil
}

func ReadBodyMeasurement_All(hkid string) (*[]utils.BodyMeasurement, error) {
	contract := GetContract("Bodymeasurement")
	funcname := "Readbodymeasurement_All"
	result, err := contract.SubmitTransaction(funcname, hkid)
	if err != nil {
		log.Println("Unable to retrieve transaction from fabric")
		log.Println(err)
		os.Exit(1)
	}
	log.Println("Successfully retrieve record by calling API : ", funcname)
	var record []utils.BodyMeasurement
	err = json.Unmarshal(result, &record)
	if err != nil {
		fmt.Println("Unmarshalling Json enocunters error")
		fmt.Println(err)
		os.Exit(1)
	}
	log.Println("Successfully Unmarshal Patient")
	return &record, nil
}

func ReadMealPlan_All(hkid string) (*[]utils.MealPlan, error) {
	contract := GetContract("mealplan")
	funcname := "Readmealplan_All"
	result, err := contract.SubmitTransaction(funcname, hkid)
	if err != nil {
		log.Println("Unable to retrieve transaction from fabric")
		log.Println(err)
		os.Exit(1)
	}
	log.Println("Successfully retrieve record by calling API : ", funcname)
	var record []utils.MealPlan
	err = json.Unmarshal(result, &record)
	if err != nil {
		fmt.Println("Unmarshalling Json enocunters error")
		fmt.Println(err)
		os.Exit(1)
	}
	log.Println("Successfully Unmarshal Patient")
	return &record, nil
}

func ReadAssessment_All(hkid string) (*[]utils.Assessment, error) {
	contract := GetContract("assessment")
	funcname := "Readassessment_All"
	result, err := contract.SubmitTransaction(funcname, hkid)
	if err != nil {
		log.Println("Unable to retrieve transaction from fabric")
		log.Println(err)
		os.Exit(1)
	}
	log.Println("Successfully retrieve record by calling API : ", funcname)
	var record []utils.Assessment
	err = json.Unmarshal(result, &record)
	if err != nil {
		fmt.Println("Unmarshalling Json enocunters error")
		fmt.Println(err)
		os.Exit(1)
	}
	log.Println("Successfully Unmarshal Patient")
	return &record, nil
}

func ReadRiskFactor_All(hkid string) (*[]utils.RiskFactor, error) {
	contract := GetContract("riskfactor")
	funcname := "Readriskfactor_All"
	result, err := contract.SubmitTransaction(funcname, hkid)
	if err != nil {
		log.Println("Unable to retrieve transaction from fabric")
		log.Println(err)
		os.Exit(1)
	}
	log.Println("Successfully retrieve record by calling API : ", funcname)
	var record []utils.RiskFactor
	err = json.Unmarshal(result, &record)
	if err != nil {
		fmt.Println("Unmarshalling Json enocunters error")
		fmt.Println(err)
		os.Exit(1)
	}
	log.Println("Successfully Unmarshal Patient")
	return &record, nil
}

func ReadProcedure_All(hkid string) (*[]utils.Procedure, error) {
	contract := GetContract("procedure")
	funcname := "Readprocedure_All"
	result, err := contract.SubmitTransaction(funcname, hkid)
	if err != nil {
		log.Println("Unable to retrieve transaction from fabric")
		log.Println(err)
		os.Exit(1)
	}
	log.Println("Successfully retrieve record by calling API : ", funcname)
	var record []utils.Procedure
	err = json.Unmarshal(result, &record)
	if err != nil {
		fmt.Println("Unmarshalling Json enocunters error")
		fmt.Println(err)
		os.Exit(1)
	}
	log.Println("Successfully Unmarshal Patient")
	return &record, nil
}

func ReadChestPain_All(hkid string) (*[]utils.ChestPainReferral, error) {
	contract := GetContract("chestpain")
	funcname := "Readchestpainreferral_All"
	result, err := contract.SubmitTransaction(funcname, hkid)
	if err != nil {
		log.Println("Unable to retrieve transaction from fabric")
		log.Println(err)
		os.Exit(1)
	}
	log.Println("Successfully retrieve record by calling API : ", funcname)
	var record []utils.ChestPainReferral
	err = json.Unmarshal(result, &record)
	if err != nil {
		fmt.Println("Unmarshalling Json enocunters error")
		fmt.Println(err)
		os.Exit(1)
	}
	log.Println("Successfully Unmarshal Patient")
	return &record, nil
}

func GetContract(info string) *gateway.Contract {
	log.Printf("============ Reading Latest Patient Record : %v============", info)

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
	return contract
}

// func ReadTEST(hkid string, info string) (*utils.Patient, error) {
// 	log.Println("============ Reading Latest Patient Record ============")

// 	err := os.Setenv("DISCOVERY_AS_LOCALHOST", "true")
// 	if err != nil {
// 		log.Fatalf("Error setting DISCOVERY_AS_LOCALHOST environemnt variable: %v", err)
// 	}

// 	wallet, err := gateway.NewFileSystemWallet("wallet")
// 	if err != nil {
// 		log.Fatalf("Failed to create wallet: %v", err)
// 	}

// 	if !wallet.Exists("appUser") {
// 		err = populateWallet(wallet)
// 		if err != nil {
// 			log.Fatalf("Failed to populate wallet contents: %v", err)
// 		}
// 	}

// 	ccpPath := filepath.Join(
// 		"..",
// 		"fabric-network",
// 		"crypto-config",
// 		"peerOrganizations",
// 		"WalesHospital.com",
// 		"connection-org1.yaml",
// 	)

// 	gw, err := gateway.Connect(
// 		gateway.WithConfig(config.FromFile(filepath.Clean(ccpPath))),
// 		gateway.WithIdentity(wallet, "appUser"),
// 	)
// 	if err != nil {
// 		log.Fatalf("Failed to connect to gateway: %v", err)
// 	}
// 	defer gw.Close()

// 	network, err := gw.GetNetwork("hospitalhk")
// 	if err != nil {
// 		log.Fatalf("Failed to get network: %v", err)
// 	}

// 	contract := network.GetContract("hospitalCC")

// 	result, err := contract.SubmitTransaction("ReadPatient", hkid, info)
// 	log.Println("The Record from Fabric is :")
// 	log.Println(string(result))

// 	var record utils.Patient
// 	err = json.Unmarshal(result, &record)
// 	if err != nil {
// 		log.Printf("error decoding response: %v", err)
// 		if e, ok := err.(*json.SyntaxError); ok {
// 			log.Printf("syntax error at byte offset %d", e.Offset)
// 		}
// 		log.Printf("response: %q", result)
// 		os.Exit(1)
// 	}
// 	log.Println("Successfully Unmarshal Patient")
// 	return &record, nil
// }
