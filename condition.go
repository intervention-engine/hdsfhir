package hdsfhir

import (
	"encoding/json"
	"gitlab.mitre.org/intervention-engine/fhir/models"
)

type Condition struct {
	Entry
	Severity map[string][]string `json:"severity"`
}

func (c *Condition) ToJSON() []byte {
	fhirCondition := models.Condition{}
	fhirCondition.Code = c.ConvertCodingToFHIR()
	fhirCondition.OnsetDate = models.FHIRDateTime{Time: c.StartTimestamp(), Precision: models.Timestamp}
	fhirCondition.Subject = models.Reference{Reference: c.Patient.ServerURL}

	if c.EndTime != 0 {
		fhirCondition.AbatementDate = models.FHIRDateTime{Time: c.EndTimestamp(), Precision: models.Timestamp}
	}
	json, _ := json.Marshal(fhirCondition)
	return json
}
