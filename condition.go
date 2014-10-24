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
		"onsetDate": UnixToFHIRDate(c.StartTime),
		"subject": map[string]string{
			"reference": c.Patient.ServerURL,
		},
	}
	if c.EndTime != 0 {
		f["abatementDate"] = UnixToFHIRDate(c.EndTime)
	}
	json, _ := json.Marshal(f)
	return json
}
