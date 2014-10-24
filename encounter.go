package hdsfhir

import (
	"encoding/json"
)

type Encounter struct {
	Entry
	DischargeDisposition map[string][]string `json:"severity"`
}

func (e *Encounter) ToJSON() []byte {
	f := map[string]interface{}{
		"type": []FHIRCodableConcept{
			FHIRCodableConcept{Codings: e.ConvertCodingToFHIR()},
		},
		"period": map[string]string{
			"start": UnixToFHIRDate(e.StartTime),
			"end":   UnixToFHIRDate(e.EndTime),
		},
		"subject": map[string]string{
			"reference": e.Patient.ServerURL,
		},
	}
	json, _ := json.Marshal(f)
	return json
}
