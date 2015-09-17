package hdsfhir

import (
	fhir "github.com/intervention-engine/fhir/models"
	. "gopkg.in/check.v1"
)

type ResultValueSuite struct {
}

var _ = Suite(&ResultValueSuite{})

func (s *ResultValueSuite) TestPhysicalQuantityResult(c *C) {
	result := ResultValue{Physical: &PhysicalQuantityResult{Unit: "mg/dL", Scalar: "130"}}
	models := result.FHIRModels()
	c.Assert(models, HasLen, 1)
	c.Assert(models[0], FitsTypeOf, &fhir.Observation{})
	model := models[0].(*fhir.Observation)
	c.Assert(model.Status, Equals, "final")
	c.Assert(model.ValueQuantity.Unit, Equals, "mg/dL")
	c.Assert(*model.ValueQuantity.Value, Equals, float64(130))
	c.Assert(model.ValueString, Equals, "")
	c.Assert(model.ValueCodeableConcept, IsNil)
}

func (s *ResultValueSuite) TestStringResult(c *C) {
	result := ResultValue{Physical: &PhysicalQuantityResult{Scalar: "positive"}}
	models := result.FHIRModels()
	c.Assert(models, HasLen, 1)
	c.Assert(models[0], FitsTypeOf, &fhir.Observation{})
	model := models[0].(*fhir.Observation)
	c.Assert(model.Status, Equals, "final")
	c.Assert(model.ValueString, Equals, "positive")
	c.Assert(model.ValueQuantity, IsNil)
	c.Assert(model.ValueCodeableConcept, IsNil)
}

func (s *ResultValueSuite) TestCodedResult(c *C) {
	result := ResultValue{
		Coded: &CodedResult{
			Description: "PHQ-9 Tool",
			Codes: map[string][]string{
				"LOINC": []string{"44249-1"},
			},
		},
	}
	models := result.FHIRModels()
	c.Assert(models, HasLen, 1)
	c.Assert(models[0], FitsTypeOf, &fhir.Observation{})
	model := models[0].(*fhir.Observation)
	c.Assert(model.Status, Equals, "final")
	c.Assert(model.ValueCodeableConcept.Text, Equals, "PHQ-9 Tool")
	c.Assert(model.ValueCodeableConcept.Coding, HasLen, 1)
	c.Assert(model.ValueCodeableConcept.MatchesCode("http://loinc.org", "44249-1"), Equals, true)
	c.Assert(model.ValueString, Equals, "")
	c.Assert(model.ValueQuantity, IsNil)
}
