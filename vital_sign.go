package hdsfhir

import (
	"encoding/json"
	"gitlab.mitre.org/intervention-engine/fhir/models"
	"log"
	"strconv"
)

type VitalSign struct {
	Entry
	Description    string             `json:"description"`
	Interpretation string             `json:"interpretation"`
	Values         []PhysicalQuantity `json:"values"`
}

// TODO: need to handle CodedValue as well
type PhysicalQuantity struct {
	Unit   string `json:"units"`
	Scalar string `json:"scalar"`
}

func (self *VitalSign) ToFhirModel() models.Observation {
	fhirObservation := models.Observation{}
	fhirObservation.Name = self.ConvertCodingToFHIR()
	fhirObservation.Encounter = models.Reference{Reference: self.Patient.MatchingEncounter(self.Entry).ServerURL}
	fhirObservation.Comments = self.Description

	if len(self.Values) > 1 {
		panic("cannot handle more than one value")
	} else if len(self.Values) == 1 {
		value := self.Values[0]
		if val, err := strconv.ParseFloat(value.Scalar, 64); err == nil {
			fhirObservation.ValueQuantity = models.Quantity{Units: value.Unit, Value: val}
		} else {
			fhirObservation.ValueString = value.Scalar
		}
	}

	fhirObservation.AppliesPeriod = self.AsFHIRPeriod()

	fhirObservation.Subject = models.Reference{Reference: self.Patient.ServerURL}
	x, _ := json.Marshal(fhirObservation.Subject)
	log.Println(string(x))
	return fhirObservation
}

func (self *VitalSign) ToJSON() []byte {
	fhirObservation := self.ToFhirModel()
	json, _ := json.Marshal(fhirObservation)
	return json
}
