package hdsfhir

import (
	"encoding/json"
	"strconv"

	"github.com/intervention-engine/fhir/models"
)

type ThingWithResults struct {
	Values []ResultValue `json:"values"`
}

type ResultValue struct {
	Physical *PhysicalQuantityResult
	Coded    *CodedResult
}

func (self *ThingWithResults) HandleValues(fhirObservation *models.Observation) {
	if len(self.Values) > 1 {
		// observation can have only one value
		panic("cannot handle more than one value... FHIR does not support more than one result")
	} else if len(self.Values) == 1 {
		self.HandleValue(fhirObservation, self.Values[0])
	}
}

func (self *ThingWithResults) HandleValue(fhirObservation *models.Observation, value ResultValue) {
	if value.Physical != nil {
		if val, err := strconv.ParseFloat(value.Physical.Scalar, 64); err == nil {
			fhirObservation.ValueQuantity = models.Quantity{Units: value.Physical.Unit, Value: val}
		} else {
			fhirObservation.ValueString = value.Physical.Scalar
		}
	} else {
		fhirObservation.ValueCodeableConcept = value.Coded.ConvertCodingToFHIR()
	}
}

func (self *ResultValue) UnmarshalJSON(data []byte) (err error) {
	// check if we have a coded or physical result value
	type ValueType struct {
		Type string `json:"_type"`
	}
	t := &ValueType{}
	json.Unmarshal(data, t)

	switch t.Type {
	case "CodedResultValue":
		local := &CodedResult{}
		json.Unmarshal(data, local)
		self.Coded = local
	case "PhysicalQuantityResultValue":
		local := &PhysicalQuantityResult{}
		json.Unmarshal(data, local)
		self.Physical = local
	default:
		local := &PhysicalQuantityResult{}
		json.Unmarshal(data, local)
		self.Physical = local
	}

	return nil

}

// Result Types
type PhysicalQuantityResult struct {
	Unit   string `json:"units"`
	Scalar string `json:"scalar"`
}

type CodedResult struct {
	Codes       map[string][]string `json:"codes"`
	Description string              `json:"description"`
}

func (self *CodedResult) ConvertCodingToFHIR() models.CodeableConcept {
	c := ConvertCodeMapToFHIR(self.Codes)
	c.Text = self.Description
	return c
}
