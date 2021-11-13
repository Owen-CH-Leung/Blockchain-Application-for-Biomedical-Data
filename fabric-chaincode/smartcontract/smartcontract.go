package smartcontract

import (
	//"blockchain/frontend/utils"
	"encoding/json"
	"fmt"
	"os"

	//"time"
	"fabric-chaincode/utils"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

// CreateAsset issues a new asset to the world state with given details.
func (s *SmartContract) CreateNurse(ctx contractapi.TransactionContextInterface, fullname string, nurseid string) error {

	record := utils.Nurse{
		FullName: fullname,
		ID:       nurseid,
	}

	recordJSON, err := json.Marshal(record)
	if err != nil {
		fmt.Println("Making Json enocunters error")
		fmt.Println(err)
		os.Exit(1)
	}

	return ctx.GetStub().PutState(nurseid, recordJSON)
}

func (s *SmartContract) CreateDoctor(ctx contractapi.TransactionContextInterface, fullname string, docid string,
	hospital string, specialty string) error {

	record := utils.Doctor{
		FullName:  fullname,
		ID:        docid,
		Hospital:  hospital,
		Specialty: specialty,
	}

	recordJSON, err := json.Marshal(record)
	if err != nil {
		fmt.Println("Making Json enocunters error")
		fmt.Println(err)
		os.Exit(1)
	}

	return ctx.GetStub().PutState(docid, recordJSON)
}

func (s *SmartContract) CreatePatient(ctx contractapi.TransactionContextInterface, hkid_file string, fullname string, gender string,
	hkid_encrypted string, dob string, race string, pob string, address string, contactno string, email string, maritalstatus string) error {

	exists, err := s.AssetExists(ctx, hkid_file, "patient")
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the patient %s already exist", hkid_file)
	}
	record := utils.Patient{
		FullName:      fullname,
		Gender:        gender,
		HKID:          hkid_encrypted,
		DateOfBirth:   dob,
		Race:          race,
		PlaceOfBirth:  pob,
		Address:       address,
		ContactNumber: contactno,
		Email:         email,
		MaritalStatus: maritalstatus,
	}
	key := hkid_file + "_patient"
	recordJSON, err := json.Marshal(record)
	if err != nil {
		fmt.Println("Making Json enocunters error")
		fmt.Println(err)
		os.Exit(1)
	}

	return ctx.GetStub().PutState(key, recordJSON)
}

func (s *SmartContract) CreateCOVIDVaccine(ctx contractapi.TransactionContextInterface, hkid string, date string, venue string,
	name string, dose string, death string, stroke string, heartfailure string, othersymptoms string) error {

	exists, err := s.AssetExists(ctx, hkid, "vaccine")
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the vaccine for %s already exist", hkid)
	}
	complication := utils.Complication{
		Death:         death,
		Stroke:        stroke,
		HeartFailure:  heartfailure,
		OtherSymptoms: othersymptoms,
	}
	record := utils.COVIDVaccination{
		Date:         date,
		Venue:        venue,
		Name:         name,
		Dose:         dose,
		Complication: complication,
	}
	key := hkid + "_vaccine"
	recordJSON, err := json.Marshal(record)
	if err != nil {
		fmt.Println("Making Json enocunters error")
		fmt.Println(err)
		os.Exit(1)
	}

	return ctx.GetStub().PutState(key, recordJSON)
}

func (s *SmartContract) CreateAllergy(ctx contractapi.TransactionContextInterface, hkid string, indicator string, allergydetails string) error {
	exists, err := s.AssetExists(ctx, hkid, "allergy")
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the allergy for %s already exist", hkid)
	}
	record := utils.Allergy{
		AllergyIndicator: indicator,
		AllergyTo:        allergydetails,
	}
	key := hkid + "_allergy"
	recordJSON, err := json.Marshal(record)
	if err != nil {
		fmt.Println("Making Json enocunters error")
		fmt.Println(err)
		os.Exit(1)
	}

	return ctx.GetStub().PutState(key, recordJSON)
}

func (s *SmartContract) CreateBodyMeasurement(ctx contractapi.TransactionContextInterface, hkid string,
	height string, weight string, BMI string, bloodglucose string) error {

	exists, err := s.AssetExists(ctx, hkid, "bodymeasurement")
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the bodymeasurement for %s already exist", hkid)
	}
	record := utils.BodyMeasurement{
		Height:       height,
		Weight:       weight,
		BMI:          BMI,
		BloodGlucose: bloodglucose,
	}
	key := hkid + "_bodymeasurement"
	recordJSON, err := json.Marshal(record)
	if err != nil {
		fmt.Println("Making Json enocunters error")
		fmt.Println(err)
		os.Exit(1)
	}

	return ctx.GetStub().PutState(key, recordJSON)
}

func (s *SmartContract) CreateMealPlan(ctx contractapi.TransactionContextInterface, hkid string,
	breakfast string, lunch string, dinner string, snack string, remark string) error {

	exists, err := s.AssetExists(ctx, hkid, "mealplan")
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the mealplan for %s already exist", hkid)
	}
	record := utils.MealPlan{
		Breakfast: breakfast,
		Lunch:     lunch,
		Dinner:    dinner,
		Snack:     snack,
		Remark:    remark,
	}
	key := hkid + "_mealplan"
	recordJSON, err := json.Marshal(record)
	if err != nil {
		fmt.Println("Making Json enocunters error")
		fmt.Println(err)
		os.Exit(1)
	}

	return ctx.GetStub().PutState(key, recordJSON)
}

func (s *SmartContract) CreateAssessment(ctx contractapi.TransactionContextInterface, hkid string,
	indicator string, allergyto string, attendancedate string, triage string, bloodpressure string,
	bodytemperature string, pulse string, diagnosis string, attendanc_specialty string, clinical_note string, drug_prescription string,
	examination string, discharge_dest string, admit_to string) error {

	exists, err := s.AssetExists(ctx, hkid, "assessment")
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the assessment for %s already exist", hkid)
	}

	allergy := utils.Allergy{
		AllergyIndicator: indicator,
		AllergyTo:        allergyto,
	}

	followup := utils.FollowUpPlan{
		DrugPrescription:     drug_prescription,
		Examination:          examination,
		DischargeDestination: discharge_dest,
		AdmittedTo:           admit_to,
	}
	record := utils.Assessment{
		AttendanceDate:      attendancedate,
		Triage:              triage,
		BloodPressure:       bloodpressure,
		BodyTemperature:     bodytemperature,
		Pulse:               pulse,
		Allergy:             allergy,
		Diagnosis:           diagnosis,
		AttendanceSpecialty: attendanc_specialty,
		ClinicalNote:        clinical_note,
		FollowUpPlan:        followup,
	}
	key := hkid + "_assessment"
	recordJSON, err := json.Marshal(record)
	if err != nil {
		fmt.Println("Making Json enocunters error")
		fmt.Println(err)
		os.Exit(1)
	}

	return ctx.GetStub().PutState(key, recordJSON)
}

func (s *SmartContract) CreateRiskFactor(ctx contractapi.TransactionContextInterface, hkid string,
	previousPCI string, vascular_disease string, diabetes_melilites string, hypertension string, smokinghistory string) error {

	exists, err := s.AssetExists(ctx, hkid, "riskfactor")
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the riskfactor for %s already exist", hkid)
	}
	record := utils.RiskFactor{
		PreviousPCI:       previousPCI,
		VascularDiseases:  vascular_disease,
		DiabetesMelilites: diabetes_melilites,
		Hypertension:      hypertension,
		SmokingHistory:    smokinghistory,
	}
	key := hkid + "_riskfactor"
	recordJSON, err := json.Marshal(record)
	if err != nil {
		fmt.Println("Making Json enocunters error")
		fmt.Println(err)
		os.Exit(1)
	}

	return ctx.GetStub().PutState(key, recordJSON)
}

func (s *SmartContract) CreateProcedure(ctx contractapi.TransactionContextInterface, hkid string,
	angina string, heartfailure string, ejectionfraction string, proceduredate string, urgency string,
	result string, device string, pro_type string, death string, stroke string, othersymptoms string) error {

	exists, err := s.AssetExists(ctx, hkid, "procedure")
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the procedure for %s already exist", hkid)
	}
	complication := utils.Complication{
		Death:         death,
		Stroke:        stroke,
		HeartFailure:  heartfailure,
		OtherSymptoms: othersymptoms,
	}

	preintervention := utils.PreIntervention{
		AnginaType:       angina,
		HeartFailure:     heartfailure,
		EjectionFraction: ejectionfraction,
	}

	record := utils.Procedure{
		PreIntervention: preintervention,
		ProcedureDate:   proceduredate,
		Urgency:         urgency,
		Result:          result,
		Device:          device,
		Type:            pro_type,
		Complication:    complication,
	}

	key := hkid + "_procedure"
	recordJSON, err := json.Marshal(record)
	if err != nil {
		fmt.Println("Making Json enocunters error")
		fmt.Println(err)
		os.Exit(1)
	}

	return ctx.GetStub().PutState(key, recordJSON)
}

func (s *SmartContract) CreateChestPainReferral(ctx contractapi.TransactionContextInterface, hkid string,
	referto string, reason string, duration string, severity string, angina_symptom string,
	history_mi string, history_ptca string, chest_image string) error {

	exists, err := s.AssetExists(ctx, hkid, "chestpainreferral")
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the chestpainreferral for %s already exist", hkid)
	}
	record := utils.ChestPainReferral{
		ReferTo:        referto,
		Reason:         reason,
		Duration:       duration,
		Severity:       severity,
		AnginaSymptoms: angina_symptom,
		HistoryMI:      history_mi,
		HistoryPTCA:    history_ptca,
		ChestXRayImage: chest_image,
	}

	key := hkid + "_chestpainreferral"
	recordJSON, err := json.Marshal(record)
	if err != nil {
		fmt.Println("Making Json enocunters error")
		fmt.Println(err)
		os.Exit(1)
	}

	return ctx.GetStub().PutState(key, recordJSON)
}

func (s *SmartContract) AssetExists(ctx contractapi.TransactionContextInterface, hkid string, info string) (bool, error) {
	var key string
	switch {
	case info == "patient":
		key = hkid + "_patient"
	case info == "vaccine":
		key = hkid + "_vaccine"
	case info == "allergy":
		key = hkid + "_allergy"
	case info == "bodymeasurement":
		key = hkid + "_bodymeasurement"
	case info == "mealplan":
		key = hkid + "_mealplan"
	case info == "assessment":
		key = hkid + "_assessment"
	case info == "riskfactor":
		key = hkid + "_riskfactor"
	case info == "procedure":
		key = hkid + "_procedure"
	case info == "chestpainreferral":
		key = hkid + "_chestpainreferral"
	default:
		key = hkid + "_patient"
	}
	assetJSON, err := ctx.GetStub().GetState(key)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return assetJSON != nil, nil
}

func (s *SmartContract) UpdatePatient(ctx contractapi.TransactionContextInterface, hkid_file string, fullname string, gender string,
	hkid_encrypt string, dob string, race string, pob string, address string, contactno string, email string, maritalstatus string) error {

	exists, err := s.AssetExists(ctx, hkid_file, "patient")
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the patient %s does not exist", hkid_file)
	}
	record := utils.Patient{
		FullName:      fullname,
		Gender:        gender,
		HKID:          hkid_encrypt,
		DateOfBirth:   dob,
		Race:          race,
		PlaceOfBirth:  pob,
		Address:       address,
		ContactNumber: contactno,
		Email:         email,
		MaritalStatus: maritalstatus,
	}
	key := hkid_file + "_patient"
	recordJSON, err := json.Marshal(record)
	if err != nil {
		fmt.Println("Making Json enocunters error")
		fmt.Println(err)
		os.Exit(1)
	}

	return ctx.GetStub().PutState(key, recordJSON)
}

func (s *SmartContract) UpdateCOVIDVaccine(ctx contractapi.TransactionContextInterface, hkid string, date string, venue string,
	name string, dose string, death string, stroke string, heartfailure string, othersymptoms string) error {

	exists, err := s.AssetExists(ctx, hkid, "vaccine")
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the vaccine for %s does not exist", hkid)
	}
	complication := utils.Complication{
		Death:         death,
		Stroke:        stroke,
		HeartFailure:  heartfailure,
		OtherSymptoms: othersymptoms,
	}
	record := utils.COVIDVaccination{
		Date:         date,
		Venue:        venue,
		Name:         name,
		Dose:         dose,
		Complication: complication,
	}
	key := hkid + "_vaccine"
	recordJSON, err := json.Marshal(record)
	if err != nil {
		fmt.Println("Making Json enocunters error")
		fmt.Println(err)
		os.Exit(1)
	}

	return ctx.GetStub().PutState(key, recordJSON)
}

func (s *SmartContract) UpdateAllergy(ctx contractapi.TransactionContextInterface, hkid string, indicator string, allergydetails string) error {
	exists, err := s.AssetExists(ctx, hkid, "allergy")
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the allergy for %s does not exist", hkid)
	}
	record := utils.Allergy{
		AllergyIndicator: indicator,
		AllergyTo:        allergydetails,
	}
	key := hkid + "_allergy"
	recordJSON, err := json.Marshal(record)
	if err != nil {
		fmt.Println("Making Json enocunters error")
		fmt.Println(err)
		os.Exit(1)
	}

	return ctx.GetStub().PutState(key, recordJSON)
}

func (s *SmartContract) UpdateBodyMeasurement(ctx contractapi.TransactionContextInterface, hkid string,
	height string, weight string, BMI string, bloodglucose string) error {

	exists, err := s.AssetExists(ctx, hkid, "bodymeasurement")
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the bodymeasurement for %s does not exist", hkid)
	}
	record := utils.BodyMeasurement{
		Height:       height,
		Weight:       weight,
		BMI:          BMI,
		BloodGlucose: bloodglucose,
	}
	key := hkid + "_bodymeasurement"
	recordJSON, err := json.Marshal(record)
	if err != nil {
		fmt.Println("Making Json enocunters error")
		fmt.Println(err)
		os.Exit(1)
	}

	return ctx.GetStub().PutState(key, recordJSON)
}

func (s *SmartContract) UpdateMealPlan(ctx contractapi.TransactionContextInterface, hkid string,
	breakfast string, lunch string, dinner string, snack string, remark string) error {

	exists, err := s.AssetExists(ctx, hkid, "mealplan")
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the mealplan for %s does not exist", hkid)
	}
	record := utils.MealPlan{
		Breakfast: breakfast,
		Lunch:     lunch,
		Dinner:    dinner,
		Snack:     snack,
		Remark:    remark,
	}
	key := hkid + "_mealplan"
	recordJSON, err := json.Marshal(record)
	if err != nil {
		fmt.Println("Making Json enocunters error")
		fmt.Println(err)
		os.Exit(1)
	}

	return ctx.GetStub().PutState(key, recordJSON)
}

func (s *SmartContract) UpdateAssessment(ctx contractapi.TransactionContextInterface, hkid string,
	indicator string, allergyto string, attendancedate string, triage string, bloodpressure string,
	bodytemperature string, pulse string, diagnosis string, attendanc_specialty string, clinical_note string, drug_prescription string,
	examination string, discharge_dest string, admit_to string) error {

	exists, err := s.AssetExists(ctx, hkid, "assessment")
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the assessment for %s does not exist", hkid)
	}

	allergy := utils.Allergy{
		AllergyIndicator: indicator,
		AllergyTo:        allergyto,
	}

	followup := utils.FollowUpPlan{
		DrugPrescription:     drug_prescription,
		Examination:          examination,
		DischargeDestination: discharge_dest,
		AdmittedTo:           admit_to,
	}
	record := utils.Assessment{
		AttendanceDate:      attendancedate,
		Triage:              triage,
		BloodPressure:       bloodpressure,
		BodyTemperature:     bodytemperature,
		Pulse:               pulse,
		Allergy:             allergy,
		Diagnosis:           diagnosis,
		AttendanceSpecialty: attendanc_specialty,
		ClinicalNote:        clinical_note,
		FollowUpPlan:        followup,
	}
	key := hkid + "_assessment"
	recordJSON, err := json.Marshal(record)
	if err != nil {
		fmt.Println("Making Json enocunters error")
		fmt.Println(err)
		os.Exit(1)
	}

	return ctx.GetStub().PutState(key, recordJSON)
}

func (s *SmartContract) UpdateRiskFactor(ctx contractapi.TransactionContextInterface, hkid string,
	previousPCI string, vascular_disease string, diabetes_melilites string, hypertension string, smokinghistory string) error {

	exists, err := s.AssetExists(ctx, hkid, "riskfactor")
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the riskfactor for %s does not exist", hkid)
	}
	record := utils.RiskFactor{
		PreviousPCI:       previousPCI,
		VascularDiseases:  vascular_disease,
		DiabetesMelilites: diabetes_melilites,
		Hypertension:      hypertension,
		SmokingHistory:    smokinghistory,
	}
	key := hkid + "_riskfactor"
	recordJSON, err := json.Marshal(record)
	if err != nil {
		fmt.Println("Making Json enocunters error")
		fmt.Println(err)
		os.Exit(1)
	}

	return ctx.GetStub().PutState(key, recordJSON)
}

func (s *SmartContract) UpdateProcedure(ctx contractapi.TransactionContextInterface, hkid string,
	angina string, heartfailure string, ejectionfraction string, proceduredate string, urgency string,
	result string, device string, pro_type string, death string, stroke string, othersymptoms string) error {

	exists, err := s.AssetExists(ctx, hkid, "procedure")
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the procedure for %s does not exist", hkid)
	}
	complication := utils.Complication{
		Death:         death,
		Stroke:        stroke,
		HeartFailure:  heartfailure,
		OtherSymptoms: othersymptoms,
	}

	preintervention := utils.PreIntervention{
		AnginaType:       angina,
		HeartFailure:     heartfailure,
		EjectionFraction: ejectionfraction,
	}

	record := utils.Procedure{
		PreIntervention: preintervention,
		ProcedureDate:   proceduredate,
		Urgency:         urgency,
		Result:          result,
		Device:          device,
		Type:            pro_type,
		Complication:    complication,
	}

	key := hkid + "_procedure"
	recordJSON, err := json.Marshal(record)
	if err != nil {
		fmt.Println("Making Json enocunters error")
		fmt.Println(err)
		os.Exit(1)
	}

	return ctx.GetStub().PutState(key, recordJSON)
}

func (s *SmartContract) UpdateChestPainReferral(ctx contractapi.TransactionContextInterface, hkid string,
	referto string, reason string, duration string, severity string, angina_symptom string,
	history_mi string, history_ptca string, chest_image string) error {

	exists, err := s.AssetExists(ctx, hkid, "chestpainreferral")
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the chestpainreferral for %s does not exist", hkid)
	}
	record := utils.ChestPainReferral{
		ReferTo:        referto,
		Reason:         reason,
		Duration:       duration,
		Severity:       severity,
		AnginaSymptoms: angina_symptom,
		HistoryMI:      history_mi,
		HistoryPTCA:    history_ptca,
		ChestXRayImage: chest_image,
	}

	key := hkid + "_chestpainreferral"
	recordJSON, err := json.Marshal(record)
	if err != nil {
		fmt.Println("Making Json enocunters error")
		fmt.Println(err)
		os.Exit(1)
	}

	return ctx.GetStub().PutState(key, recordJSON)
}

func (s *SmartContract) Readpatient_All(ctx contractapi.TransactionContextInterface, hkid string) ([]*utils.Patient, error) {
	var listofRecord []*utils.Patient
	key := hkid + "_patient"
	resultiterator, err := ctx.GetStub().GetHistoryForKey(key)
	if err != nil {
		fmt.Println("Reading from fabric encounters error")
		fmt.Println(err)
		os.Exit(1)
	}
	defer resultiterator.Close()
	for resultiterator.HasNext() {
		var record utils.Patient
		response, err := resultiterator.Next()
		if err != nil {
			fmt.Println("Inside Loop there's an error")
			fmt.Println(err)
			os.Exit(1)
		}
		err = json.Unmarshal(response.Value, &record)
		if err != nil {
			fmt.Println("Unmarshalling Json enocunters error")
			fmt.Println(err)
			os.Exit(1)
		}
		listofRecord = append(listofRecord, &record)
	}
	return listofRecord, nil
}

func (s *SmartContract) Readvaccine_All(ctx contractapi.TransactionContextInterface, hkid string) ([]*utils.COVIDVaccination, error) {
	key := hkid + "_vaccine"
	var listofRecord []*utils.COVIDVaccination
	resultiterator, err := ctx.GetStub().GetHistoryForKey(key)
	if err != nil {
		fmt.Println("Reading from fabric encounters error")
		fmt.Println(err)
		os.Exit(1)
	}
	defer resultiterator.Close()
	for resultiterator.HasNext() {
		var record utils.COVIDVaccination
		response, err := resultiterator.Next()
		if err != nil {
			fmt.Println("Inside Loop there's an error")
			fmt.Println(err)
			os.Exit(1)
		}
		err = json.Unmarshal(response.Value, &record)
		if err != nil {
			fmt.Println("Unmarshalling Json enocunters error")
			fmt.Println(err)
			os.Exit(1)
		}
		listofRecord = append(listofRecord, &record)
	}
	return listofRecord, nil
}

func (s *SmartContract) Readallergy_All(ctx contractapi.TransactionContextInterface, hkid string) ([]*utils.Allergy, error) {
	key := hkid + "_allergy"
	var listofRecord []*utils.Allergy
	resultiterator, err := ctx.GetStub().GetHistoryForKey(key)
	if err != nil {
		fmt.Println("Reading from fabric encounters error")
		fmt.Println(err)
		os.Exit(1)
	}
	defer resultiterator.Close()
	for resultiterator.HasNext() {
		var record utils.Allergy
		response, err := resultiterator.Next()
		if err != nil {
			fmt.Println("Inside Loop there's an error")
			fmt.Println(err)
			os.Exit(1)
		}
		err = json.Unmarshal(response.Value, &record)
		if err != nil {
			fmt.Println("Unmarshalling Json enocunters error")
			fmt.Println(err)
			os.Exit(1)
		}
		listofRecord = append(listofRecord, &record)
	}
	return listofRecord, nil
}

func (s *SmartContract) Readbodymeasurement_All(ctx contractapi.TransactionContextInterface, hkid string) ([]*utils.BodyMeasurement, error) {
	key := hkid + "_bodymeasurement"
	var listofRecord []*utils.BodyMeasurement
	resultiterator, err := ctx.GetStub().GetHistoryForKey(key)
	if err != nil {
		fmt.Println("Reading from fabric encounters error")
		fmt.Println(err)
		os.Exit(1)
	}
	defer resultiterator.Close()
	for resultiterator.HasNext() {
		var record utils.BodyMeasurement
		response, err := resultiterator.Next()
		if err != nil {
			fmt.Println("Inside Loop there's an error")
			fmt.Println(err)
			os.Exit(1)
		}
		err = json.Unmarshal(response.Value, &record)
		if err != nil {
			fmt.Println("Unmarshalling Json enocunters error")
			fmt.Println(err)
			os.Exit(1)
		}
		listofRecord = append(listofRecord, &record)
	}
	return listofRecord, nil
}

func (s *SmartContract) Readmealplan_All(ctx contractapi.TransactionContextInterface, hkid string) ([]*utils.MealPlan, error) {
	key := hkid + "_mealplan"
	var listofRecord []*utils.MealPlan
	resultiterator, err := ctx.GetStub().GetHistoryForKey(key)
	if err != nil {
		fmt.Println("Reading from fabric encounters error")
		fmt.Println(err)
		os.Exit(1)
	}
	defer resultiterator.Close()
	for resultiterator.HasNext() {
		var record utils.MealPlan
		response, err := resultiterator.Next()
		if err != nil {
			fmt.Println("Inside Loop there's an error")
			fmt.Println(err)
			os.Exit(1)
		}
		err = json.Unmarshal(response.Value, &record)
		if err != nil {
			fmt.Println("Unmarshalling Json enocunters error")
			fmt.Println(err)
			os.Exit(1)
		}
		listofRecord = append(listofRecord, &record)
	}
	return listofRecord, nil
}

func (s *SmartContract) Readassessment_All(ctx contractapi.TransactionContextInterface, hkid string) ([]*utils.Assessment, error) {
	key := hkid + "_assessment"
	var listofRecord []*utils.Assessment
	resultiterator, err := ctx.GetStub().GetHistoryForKey(key)
	if err != nil {
		fmt.Println("Reading from fabric encounters error")
		fmt.Println(err)
		os.Exit(1)
	}
	defer resultiterator.Close()
	for resultiterator.HasNext() {
		var record utils.Assessment
		response, err := resultiterator.Next()
		if err != nil {
			fmt.Println("Inside Loop there's an error")
			fmt.Println(err)
			os.Exit(1)
		}
		err = json.Unmarshal(response.Value, &record)
		if err != nil {
			fmt.Println("Unmarshalling Json enocunters error")
			fmt.Println(err)
			os.Exit(1)
		}
		listofRecord = append(listofRecord, &record)
	}
	return listofRecord, nil
}

func (s *SmartContract) Readriskfactor_All(ctx contractapi.TransactionContextInterface, hkid string) ([]*utils.RiskFactor, error) {
	key := hkid + "_riskfactor"
	var listofRecord []*utils.RiskFactor
	resultiterator, err := ctx.GetStub().GetHistoryForKey(key)
	if err != nil {
		fmt.Println("Reading from fabric encounters error")
		fmt.Println(err)
		os.Exit(1)
	}
	defer resultiterator.Close()
	for resultiterator.HasNext() {
		var record utils.RiskFactor
		response, err := resultiterator.Next()
		if err != nil {
			fmt.Println("Inside Loop there's an error")
			fmt.Println(err)
			os.Exit(1)
		}
		err = json.Unmarshal(response.Value, &record)
		if err != nil {
			fmt.Println("Unmarshalling Json enocunters error")
			fmt.Println(err)
			os.Exit(1)
		}
		listofRecord = append(listofRecord, &record)
	}
	return listofRecord, nil
}

func (s *SmartContract) Readprocedure_All(ctx contractapi.TransactionContextInterface, hkid string) ([]*utils.Procedure, error) {
	key := hkid + "_procedure"
	var listofRecord []*utils.Procedure
	resultiterator, err := ctx.GetStub().GetHistoryForKey(key)
	if err != nil {
		fmt.Println("Reading from fabric encounters error")
		fmt.Println(err)
		os.Exit(1)
	}
	defer resultiterator.Close()
	for resultiterator.HasNext() {
		var record utils.Procedure
		response, err := resultiterator.Next()
		if err != nil {
			fmt.Println("Inside Loop there's an error")
			fmt.Println(err)
			os.Exit(1)
		}
		err = json.Unmarshal(response.Value, &record)
		if err != nil {
			fmt.Println("Unmarshalling Json enocunters error")
			fmt.Println(err)
			os.Exit(1)
		}
		listofRecord = append(listofRecord, &record)
	}
	return listofRecord, nil
}

func (s *SmartContract) Readchestpainreferral_All(ctx contractapi.TransactionContextInterface, hkid string) ([]*utils.ChestPainReferral, error) {
	key := hkid + "_chestpainreferral"
	var listofRecord []*utils.ChestPainReferral
	resultiterator, err := ctx.GetStub().GetHistoryForKey(key)
	if err != nil {
		fmt.Println("Reading from fabric encounters error")
		fmt.Println(err)
		os.Exit(1)
	}
	defer resultiterator.Close()
	for resultiterator.HasNext() {
		var record utils.ChestPainReferral
		response, err := resultiterator.Next()
		if err != nil {
			fmt.Println("Inside Loop there's an error")
			fmt.Println(err)
			os.Exit(1)
		}
		err = json.Unmarshal(response.Value, &record)
		if err != nil {
			fmt.Println("Unmarshalling Json enocunters error")
			fmt.Println(err)
			os.Exit(1)
		}
		listofRecord = append(listofRecord, &record)
	}
	return listofRecord, nil
}

func (s *SmartContract) Readpatient(ctx contractapi.TransactionContextInterface, hkid string) (*utils.Patient, error) {
	var record utils.Patient
	key := hkid + "_patient"
	assetJSON, err := ctx.GetStub().GetState(key)
	if err != nil {
		fmt.Println("Reading from fabric encounters error")
		fmt.Println(err)
		os.Exit(1)
	}
	if assetJSON == nil {
		fmt.Println("Such record didn't exist")
		os.Exit(1)
	}
	err = json.Unmarshal(assetJSON, &record)
	if err != nil {
		fmt.Println("Unmarshalling Json enocunters error")
		fmt.Println(err)
		os.Exit(1)
	}

	return &record, nil
}

func (s *SmartContract) Readvaccine(ctx contractapi.TransactionContextInterface, hkid string) (*utils.COVIDVaccination, error) {
	var record utils.COVIDVaccination
	key := hkid + "_vaccine"
	assetJSON, err := ctx.GetStub().GetState(key)
	if err != nil {
		fmt.Println("Reading from fabric encounters error")
		fmt.Println(err)
		os.Exit(1)
	}
	if assetJSON == nil {
		fmt.Println("Such record didn't exist")
		os.Exit(1)
	}
	err = json.Unmarshal(assetJSON, &record)
	if err != nil {
		fmt.Println("Unmarshalling Json enocunters error")
		fmt.Println(err)
		os.Exit(1)
	}

	return &record, nil
}

func (s *SmartContract) Readallergy(ctx contractapi.TransactionContextInterface, hkid string) (*utils.Allergy, error) {
	var record utils.Allergy
	key := hkid + "_allergy"
	assetJSON, err := ctx.GetStub().GetState(key)
	if err != nil {
		fmt.Println("Reading from fabric encounters error")
		fmt.Println(err)
		os.Exit(1)
	}
	if assetJSON == nil {
		fmt.Println("Such record didn't exist")
		os.Exit(1)
	}
	err = json.Unmarshal(assetJSON, &record)
	if err != nil {
		fmt.Println("Unmarshalling Json enocunters error")
		fmt.Println(err)
		os.Exit(1)
	}

	return &record, nil
}

func (s *SmartContract) Readbodymeasurement(ctx contractapi.TransactionContextInterface, hkid string) (*utils.BodyMeasurement, error) {
	var record utils.BodyMeasurement
	key := hkid + "_bodymeasurement"
	assetJSON, err := ctx.GetStub().GetState(key)
	if err != nil {
		fmt.Println("Reading from fabric encounters error")
		fmt.Println(err)
		os.Exit(1)
	}
	if assetJSON == nil {
		fmt.Println("Such record didn't exist")
		os.Exit(1)
	}
	err = json.Unmarshal(assetJSON, &record)
	if err != nil {
		fmt.Println("Unmarshalling Json enocunters error")
		fmt.Println(err)
		os.Exit(1)
	}

	return &record, nil
}

func (s *SmartContract) Readmealplan(ctx contractapi.TransactionContextInterface, hkid string) (*utils.MealPlan, error) {
	var record utils.MealPlan
	key := hkid + "_mealplan"
	assetJSON, err := ctx.GetStub().GetState(key)
	if err != nil {
		fmt.Println("Reading from fabric encounters error")
		fmt.Println(err)
		os.Exit(1)
	}
	if assetJSON == nil {
		fmt.Println("Such record didn't exist")
		os.Exit(1)
	}
	err = json.Unmarshal(assetJSON, &record)
	if err != nil {
		fmt.Println("Unmarshalling Json enocunters error")
		fmt.Println(err)
		os.Exit(1)
	}

	return &record, nil
}

func (s *SmartContract) Readassessment(ctx contractapi.TransactionContextInterface, hkid string) (*utils.Assessment, error) {
	var record utils.Assessment
	key := hkid + "_assessment"
	assetJSON, err := ctx.GetStub().GetState(key)
	if err != nil {
		fmt.Println("Reading from fabric encounters error")
		fmt.Println(err)
		os.Exit(1)
	}
	if assetJSON == nil {
		fmt.Println("Such record didn't exist")
		os.Exit(1)
	}
	err = json.Unmarshal(assetJSON, &record)
	if err != nil {
		fmt.Println("Unmarshalling Json enocunters error")
		fmt.Println(err)
		os.Exit(1)
	}

	return &record, nil
}

func (s *SmartContract) Readriskfactor(ctx contractapi.TransactionContextInterface, hkid string) (*utils.RiskFactor, error) {
	var record utils.RiskFactor
	key := hkid + "_riskfactor"
	assetJSON, err := ctx.GetStub().GetState(key)
	if err != nil {
		fmt.Println("Reading from fabric encounters error")
		fmt.Println(err)
		os.Exit(1)
	}
	if assetJSON == nil {
		fmt.Println("Such record didn't exist")
		os.Exit(1)
	}
	err = json.Unmarshal(assetJSON, &record)
	if err != nil {
		fmt.Println("Unmarshalling Json enocunters error")
		fmt.Println(err)
		os.Exit(1)
	}

	return &record, nil
}

func (s *SmartContract) Readprocedure(ctx contractapi.TransactionContextInterface, hkid string) (*utils.Procedure, error) {
	var record utils.Procedure
	key := hkid + "_procedure"
	assetJSON, err := ctx.GetStub().GetState(key)
	if err != nil {
		fmt.Println("Reading from fabric encounters error")
		fmt.Println(err)
		os.Exit(1)
	}
	if assetJSON == nil {
		fmt.Println("Such record didn't exist")
		os.Exit(1)
	}
	err = json.Unmarshal(assetJSON, &record)
	if err != nil {
		fmt.Println("Unmarshalling Json enocunters error")
		fmt.Println(err)
		os.Exit(1)
	}

	return &record, nil
}

func (s *SmartContract) Readchestpainreferral(ctx contractapi.TransactionContextInterface, hkid string) (*utils.ChestPainReferral, error) {
	var record utils.ChestPainReferral
	key := hkid + "_chestpainreferral"
	assetJSON, err := ctx.GetStub().GetState(key)
	if err != nil {
		fmt.Println("Reading from fabric encounters error")
		fmt.Println(err)
		os.Exit(1)
	}
	if assetJSON == nil {
		fmt.Println("Such record didn't exist")
		os.Exit(1)
	}
	err = json.Unmarshal(assetJSON, &record)
	if err != nil {
		fmt.Println("Unmarshalling Json enocunters error")
		fmt.Println(err)
		os.Exit(1)
	}

	return &record, nil
}

// func (s *SmartContract) ReadPatientTEST(ctx contractapi.TransactionContextInterface, hkid string) (*utils.Patient, error) {
// 	var key string
// 	var record utils.Patient
// 	key = hkid + "_patient"
// 	assetJSON, err := ctx.GetStub().GetState(key)
// 	if err != nil {
// 		fmt.Println("Reading from fabric encounters error")
// 		fmt.Println(err)
// 		os.Exit(1)
// 	}
// 	if assetJSON == nil {
// 		fmt.Println("Such record didn't exist")
// 		os.Exit(1)
// 	}
// 	err = json.Unmarshal(assetJSON, &record)
// 	if err != nil {
// 		fmt.Println("Unmarshalling Json enocunters error")
// 		fmt.Println(err)
// 		os.Exit(1)
// 	}

// 	return &record, nil
// }
