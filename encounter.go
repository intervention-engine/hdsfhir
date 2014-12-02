package hdsfhir

import (
	"encoding/json"
	"gitlab.mitre.org/intervention-engine/fhir/models"
)

type Encounter struct {
	Entry
	DischargeDisposition map[string][]string `json:"severity"`
}

func (e *Encounter) Period() models.Period {
	period := models.Period{}
	period.Start = models.FHIRDateTime{Time: e.StartTimestamp(), Precision: models.Timestamp}
	period.End = models.FHIRDateTime{Time: e.EndTimestamp(), Precision: models.Timestamp}
	return period
}

func (e *Encounter) ToJSON() []byte {
	fhirEncounter := models.Encounter{}
	fhirEncounter.Type = []models.CodeableConcept{e.ConvertCodingToFHIR()}
	fhirEncounter.Period = e.Period()
	fhirEncounter.Subject = models.Reference{Reference: e.Patient.ServerURL}
	json, _ := json.Marshal(fhirEncounter)
	return json
}
