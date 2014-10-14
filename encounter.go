package hdsfhir

import (
	"encoding/json"
)

func EncounterToJSON(p *Patient, e *Entry) []byte {
	f := map[string]interface{}{
		"type": map[string][]FHIRCoding{
			"coding": e.ConvertCodingToFHIR(),
		},
		"period": map[string]string{
			"start": UnixToFHIRDate(e.StartTime),
			"end":   UnixToFHIRDate(e.EndTime),
		},
		"subject": map[string]string{
			"reference": p.ServerURL,
		},
	}
	json, _ := json.Marshal(f)
	return json
}
