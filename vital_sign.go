package hdsfhir

import (
	"encoding/json"

	"github.com/intervention-engine/fhir/models"
)

type VitalSign struct {
	Entry
	Description    string `json:"description"`
	Interpretation string `json:"interpretation"`
	ThingWithResults
}

func (self *VitalSign) ToFhirModel() models.Observation {
	fhirObservation := models.Observation{Reliability: "ok", Status: "final"}
	fhirObservation.Name = self.ConvertCodingToFHIR()
	fhirObservation.Name.Text = self.Description
	fhirObservation.Encounter = models.Reference{Reference: self.Patient.MatchingEncounter(self.Entry).ServerURL}

	fhirObservation.AppliesPeriod = self.AsFHIRPeriod()
	self.HandleValues(&fhirObservation)

	fhirObservation.Subject = models.Reference{Reference: self.Patient.ServerURL}
	return fhirObservation
}

func (self *VitalSign) ToJSON() []byte {
	fhirObservation := self.ToFhirModel()
	json, _ := json.Marshal(fhirObservation)
	return json
}
