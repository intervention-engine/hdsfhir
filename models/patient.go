package models

import (
	"encoding/json"
	"time"

	fhir "github.com/intervention-engine/fhir/models"
)

type Patient struct {
	FirstName     string        `json:"first"`
	LastName      string        `json:"last"`
	UnixBirthTime int64         `json:"birthdate"`
	Gender        string        `json:"gender"`
	Encounters    []*Encounter  `json:"encounters"`
	Conditions    []*Condition  `json:"conditions"`
	VitalSigns    []*VitalSign  `json:"vital_signs"`
	Procedures    []*Procedure  `json:"procedures"`
	Medications   []*Medication `json:"medications"`
	ServerURL     string        `json:"-"`
}

// TODO: :allergies, :care_goals, :immunizations, :medical_equipment, :results, :social_history, :support, :advance_directives, :insurance_providers, :functional_statuses

func (p *Patient) SetServerURL(url string) {
	p.ServerURL = url
}

func (p *Patient) BirthTime() time.Time {
	return time.Unix(p.UnixBirthTime, 0)
}

func (p *Patient) PostToFHIRServer(baseURL string) {
	Upload(p, baseURL+"/Patient")
	for _, encounter := range p.Encounters {
		encounter.Patient = p
		Upload(encounter, baseURL+"/Encounter")
	}
	for _, condition := range p.Conditions {
		condition.Patient = p
		Upload(condition, baseURL+"/Condition")
	}
	for _, observation := range p.VitalSigns {
		// find matching encounter
		observation.Patient = p
		Upload(observation, baseURL+"/Observation")
	}

	for _, procedure := range p.Procedures {
		procedure.Patient = p
		procedure.UploadResults(baseURL)
		Upload(procedure, baseURL+"/Procedure")

	}

	for _, med := range p.Medications {
		_, cvxExists := med.Codes["CVX"] // Ignores medications that are coded with CVX
		if cvxExists != true {
			med.Patient = p
			med.BaseUrl = baseURL
			Upload(med, baseURL+"/MedicationStatement")
		}
	}
}

func (self *Patient) MatchingEncounter(entry Entry) Encounter {
	for _, encounter := range self.Encounters {
		// TODO: Overlaps may not be the right thing here... maybe closest?
		if encounter.StartTime <= entry.EndTime && encounter.EndTime >= entry.StartTime {
			return *encounter
		}
	}
	return Encounter{}
}

func (p *Patient) ToJSON() []byte {
	fhirPatient := &fhir.Patient{}
	fhirPatient.Name = []fhir.HumanName{fhir.HumanName{Given: []string{p.FirstName}, Family: []string{p.LastName}}}
	fhirPatient.Gender = &fhir.CodeableConcept{Coding: []fhir.Coding{fhir.Coding{System: "http://hl7.org/fhir/v3/AdministrativeGender", Code: p.Gender}}}
	fhirPatient.BirthDate = &fhir.FHIRDateTime{Time: p.BirthTime(), Precision: fhir.Precision("date")}
	json, _ := json.Marshal(fhirPatient)
	return json
}
