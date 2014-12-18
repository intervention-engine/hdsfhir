package hdsfhir

import (
	"encoding/json"
	"github.com/intervention-engine/fhir/models"
)

type Condition struct {
	Entry
	Severity map[string][]string `json:"severity"`
}

func (c *Condition) ToFhirModel() models.Condition {
	fhirCondition := models.Condition{}
	fhirCondition.Code = c.ConvertCodingToFHIR()
	fhirCondition.OnsetDate = models.FHIRDateTime{Time: c.StartTimestamp(), Precision: models.Timestamp}
	fhirCondition.Subject = models.Reference{Reference: c.Patient.ServerURL}

	if c.EndTime != 0 {
		fhirCondition.AbatementDate = models.FHIRDateTime{Time: c.EndTimestamp(), Precision: models.Timestamp}
	}
	return fhirCondition
}

func (c *Condition) ToJSON() []byte {
	json, _ := json.Marshal(c.ToFhirModel())
	return json
}
