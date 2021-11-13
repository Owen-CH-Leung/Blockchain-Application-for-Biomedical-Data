package utils

import "encoding/json"

type Patient struct {
	FullName      string
	Gender        string
	HKID          string
	DateOfBirth   string
	Race          string
	PlaceOfBirth  string
	Address       string
	ContactNumber string
	Email         string
	MaritalStatus string
}

type Nurse struct {
	FullName string
	ID       string
}

type Doctor struct {
	FullName  string
	ID        string
	Hospital  string
	Specialty string
}

type COVIDVaccination struct {
	Date         string
	Venue        string
	Name         string
	Dose         string
	Complication Complication
}

type Complication struct {
	Death         string
	Stroke        string
	HeartFailure  string
	OtherSymptoms string
}

type Allergy struct {
	AllergyIndicator string
	AllergyTo        string
}

type BodyMeasurement struct {
	Height       string
	Weight       string
	BMI          string
	BloodGlucose string
}

type MealPlan struct {
	Breakfast string
	Lunch     string
	Dinner    string
	Snack     string
	Remark    string
}

type Assessment struct {
	AttendanceDate      string
	Triage              string
	BloodPressure       string
	BodyTemperature     string
	Pulse               string
	Allergy             Allergy
	Diagnosis           string
	AttendanceSpecialty string
	ClinicalNote        string
	FollowUpPlan        FollowUpPlan
}

type FollowUpPlan struct {
	DrugPrescription     string
	Examination          string
	DischargeDestination string
	AdmittedTo           string
}

type RiskFactor struct {
	PreviousPCI       string
	VascularDiseases  string
	DiabetesMelilites string
	Hypertension      string
	SmokingHistory    string
}

type PreIntervention struct {
	AnginaType       string
	HeartFailure     string
	EjectionFraction string
}

type Procedure struct {
	PreIntervention PreIntervention
	ProcedureDate   string
	Urgency         string
	Result          string
	Device          string
	Type            string
	Complication    Complication
}

type ChestPainReferral struct {
	ReferTo        string
	Reason         string
	Duration       string
	Severity       string
	AnginaSymptoms string
	HistoryMI      string
	HistoryPTCA    string
	ChestXRayImage string
}

func (cpr ChestPainReferral) MarshalJSON() ([]byte, error) {
	var tmp struct {
		ReferTo        string
		Reason         string
		Duration       string
		Severity       string
		AnginaSymptoms string
		HistoryMI      string
		HistoryPTCA    string
	}
	tmp.ReferTo = cpr.ReferTo
	tmp.Reason = cpr.Reason
	tmp.Duration = cpr.Duration
	tmp.Severity = cpr.Severity
	tmp.AnginaSymptoms = cpr.AnginaSymptoms
	tmp.HistoryMI = cpr.HistoryMI
	tmp.HistoryPTCA = cpr.HistoryPTCA
	return json.Marshal(&tmp)
}
