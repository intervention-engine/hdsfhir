package models

import (
	"encoding/json"

	fhir "github.com/intervention-engine/fhir/models"
)

type VitalSign struct {
	Entry
	Description    string        `json:"description"`
	Interpretation string        `json:"interpretation"`
	Values         []ResultValue `json:"values"`
}

func (self *VitalSign) FHIRModel() fhir.Observation {
	if len(self.Values) != 1 {
		panic("FHIR Observations must gave one and only one value")
	}
	observation := self.Values[0].FHIRModel()
	observation.Name = self.Codes.FHIRCodeableConcept(self.Description)
	observation.Encounter = &fhir.Reference{Reference: self.Patient.MatchingEncounter(self.Entry).ServerURL}
	observation.AppliesPeriod = self.GetFHIRPeriod()
	observation.Subject = &fhir.Reference{Reference: self.Patient.ServerURL}
	return observation
}

func (self *VitalSign) ToJSON() []byte {
	fhirObservation := self.FHIRModel()
	json, _ := json.Marshal(fhirObservation)
	return json
}
