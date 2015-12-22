package hdsfhir

import (
	"encoding/json"
	"strconv"

	fhir "github.com/intervention-engine/fhir/models"
)

type ResultValue struct {
	TemporallyIdentified
	Physical *PhysicalQuantityResult
	Coded    *CodedResult
}

func (v *ResultValue) FHIRModels() []interface{} {
	observation := &fhir.Observation{}
	observation.Id = v.GetTempID()
	observation.Status = "final"
	if v.Physical != nil {
		if val, err := strconv.ParseFloat(v.Physical.Scalar, 64); err == nil {
			observation.ValueQuantity = &fhir.Quantity{Unit: v.Physical.Unit, Value: &val}
		} else {
			observation.ValueString = v.Physical.Scalar
		}
	} else {
		observation.ValueCodeableConcept = v.Coded.Codes.FHIRCodeableConcept(v.Coded.Description)
	}

	return []interface{}{observation}
}

func (v *ResultValue) UnmarshalJSON(data []byte) (err error) {
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
		v.Coded = local
	case "PhysicalQuantityResultValue":
		local := &PhysicalQuantityResult{}
		json.Unmarshal(data, local)
		v.Physical = local
	default:
		local := &PhysicalQuantityResult{}
		json.Unmarshal(data, local)
		v.Physical = local
	}

	return nil

}

// Result Types
type PhysicalQuantityResult struct {
	Unit   string `json:"unit"`
	Scalar string `json:"scalar"`
}

type CodedResult struct {
	Codes       CodeMap `json:"codes"`
	Description string  `json:"description"`
}
