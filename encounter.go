package hdsfhir

import (
	"encoding/json"
	"github.com/intervention-engine/fhir/models"
)

type Encounter struct {
	Entry
	DischargeDisposition map[string][]string `json:"severity"`
}

func (e *Encounter) ToFhirModel() models.Encounter {
	fhirEncounter := models.Encounter{}
	cc := e.ConvertCodingToFHIR()
	cc.Text = e.Description
	fhirEncounter.Type = []models.CodeableConcept{cc}
	fhirEncounter.Period = e.AsFHIRPeriod()
	fhirEncounter.Subject = models.Reference{Reference: e.Patient.ServerURL}
	return fhirEncounter
}

func (e *Encounter) ToJSON() []byte {
	json, _ := json.Marshal(e.ToFhirModel())
	return json
}
