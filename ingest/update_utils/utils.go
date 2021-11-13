package utils

var Dev bool = false

type Patient struct {
	FullName      []byte
	Gender        []byte
	HKID          []byte
	DateOfBirth   []byte
	Race          []byte
	PlaceOfBirth  []byte
	Address       []byte
	ContactNumber []byte
	Email         []byte
	MaritalStatus []byte
}

type Nurse struct {
	FullName []byte
	ID       []byte
}

type Doctor struct {
	FullName  []byte
	ID        []byte
	Hospital  []byte
	Specialty []byte
}

type COVIDVaccination struct {
	Date         []byte
	Venue        []byte
	Name         []byte
	Dose         []byte
	Complication Complication
}

type Complication struct {
	Death         []byte
	Stroke        []byte
	HeartFailure  []byte
	OtherSymptoms []byte
}

type Allergy struct {
	AllergyIndicator []byte
	AllergyTo        []byte
}

type BodyMeasurement struct {
	Height       []byte
	Weight       []byte
	BMI          []byte
	BloodGlucose []byte
}

type MealPlan struct {
	Breakfast []byte
	Lunch     []byte
	Dinner    []byte
	Snack     []byte
	Remark    []byte
}

type Assessment struct {
	AttendanceDate      []byte
	Triage              []byte
	BloodPressure       []byte
	BodyTemperature     []byte
	Pulse               []byte
	Allergy             Allergy
	Diagnosis           []byte
	AttendanceSpecialty []byte
	ClinicalNote        []byte
	FollowUpPlan        FollowUpPlan
}

type FollowUpPlan struct {
	DrugPrescription     []byte
	Examination          []byte
	DischargeDestination []byte
	AdmittedTo           []byte
}

type RiskFactor struct {
	PreviousPCI       []byte
	VascularDiseases  []byte
	DiabetesMelilites []byte
	Hypertension      []byte
	SmokingHistory    []byte
}

type PreIntervention struct {
	AnginaType       []byte
	HeartFailure     []byte
	EjectionFraction []byte
}

type Procedure struct {
	PreIntervention PreIntervention
	ProcedureDate   []byte
	Urgency         []byte
	Result          []byte
	Device          []byte
	Type            []byte
	Complication    Complication
}

type ChestPainReferral struct {
	ReferTo        []byte
	Reason         []byte
	Duration       []byte
	Severity       []byte
	AnginaSymptoms []byte
	HistoryMI      []byte
	HistoryPTCA    []byte
	ChestXRayImage []byte
}
