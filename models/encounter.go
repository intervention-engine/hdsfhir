package models

import (
	"encoding/json"

	fhir "github.com/intervention-engine/fhir/models"
)

type Encounter struct {
	Entry
	DischargeDisposition map[string][]string `json:"severity"`
}

func (e *Encounter) FHIRModel() fhir.Encounter {
	fhirEncounter := fhir.Encounter{}
	cc := e.Codes.FHIRCodeableConcept(e.Description)
	fhirEncounter.Type = []fhir.CodeableConcept{*cc}
	fhirEncounter.Period = e.GetFHIRPeriod()
	fhirEncounter.Subject = &fhir.Reference{Reference: e.Patient.ServerURL}
	return fhirEncounter
}

func (e *Encounter) ToJSON() []byte {
	json, _ := json.Marshal(e.FHIRModel())
	return json
}
