package hdsfhir

import (
	"encoding/json"
	"gitlab.mitre.org/intervention-engine/fhir/models"
	"time"
)

type Patient struct {
	FirstName     string       `json:"first"`
	LastName      string       `json:"last"`
	UnixBirthTime int64        `json:"birthdate"`
	Gender        string       `json:"gender"`
	Encounters    []*Encounter `json:"encounters"`
	Conditions    []*Condition `json:"conditions"`
	ServerURL     string       `json:"-"`
}

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
}

func (p *Patient) ToJSON() []byte {
	fhirPatient := models.Patient{}
	fhirPatient.Name = []models.HumanName{models.HumanName{Given: []string{p.FirstName}, Family: []string{p.LastName}}}
	fhirPatient.Gender = models.CodeableConcept{Coding: []models.Coding{models.Coding{System: "http://hl7.org/fhir/v3/AdministrativeGender", Code: p.Gender}}}
	fhirPatient.BirthDate = models.FHIRDateTime{Time: p.BirthTime(), Precision: models.Precision("date")}
	json, _ := json.Marshal(fhirPatient)
	return json
}
