package hdsfhir

import (
	"encoding/json"
)

type Condition struct {
	Entry
	Severity map[string][]string `json:"severity"`
}

func (c *Condition) ToJSON() []byte {
	f := map[string]interface{}{
		"code": map[string][]FHIRCoding{
			"coding": c.ConvertCodingToFHIR(),
		},
		"onsetDate":     UnixToFHIRDate(c.StartTime),
		"abatementDate": UnixToFHIRDate(c.EndTime),
		"subject": map[string]string{
			"reference": c.Patient.ServerURL,
		},
	}
	json, _ := json.Marshal(f)
	return json
}
