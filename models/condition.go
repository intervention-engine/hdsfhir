package models

import (
	"encoding/json"

	fhir "github.com/intervention-engine/fhir/models"
)

type Condition struct {
	Entry
	Severity map[string][]string `json:"severity"`
}

func (c *Condition) FHIRModel() fhir.Condition {
	fhirCondition := fhir.Condition{}
	fhirCondition.Code = c.Codes.FHIRCodeableConcept(c.Description)
	fhirCondition.OnsetDate = &fhir.FHIRDateTime{Time: c.StartTime.Time(), Precision: fhir.Timestamp}
	fhirCondition.Subject = &fhir.Reference{Reference: c.Patient.ServerURL}

	if c.EndTime != 0 {
		fhirCondition.AbatementDate = &fhir.FHIRDateTime{Time: c.EndTime.Time(), Precision: fhir.Timestamp}
	}
	return fhirCondition
}

func (c *Condition) ToJSON() []byte {
	json, _ := json.Marshal(c.FHIRModel())
	return json
}
