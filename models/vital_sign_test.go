package models

import (
	"encoding/json"
	"io/ioutil"

	fhir "github.com/intervention-engine/fhir/models"
	"github.com/pebbe/util"
	. "gopkg.in/check.v1"
)

type VitalSignSuite struct {
	Patient *Patient
}

var _ = Suite(&VitalSignSuite{})

func (suite *VitalSignSuite) SetUpSuite(c *C) {
	data, err := ioutil.ReadFile("../fixtures/john_peters.json")
	util.CheckErr(err)

	suite.Patient = &Patient{}
	err = json.Unmarshal(data, suite.Patient)
	util.CheckErr(err)
}

func (suite *VitalSignSuite) TestFHIRModels(c *C) {
	models := suite.Patient.VitalSigns[0].FHIRModels()
	c.Assert(models, HasLen, 1)
	c.Assert(models[0], FitsTypeOf, &fhir.Observation{})

	data := models[0].(*fhir.Observation)
	c.Assert(data.Subject, DeepEquals, suite.Patient.FHIRReference())
	c.Assert(data.Name.Text, Equals, "Laboratory Test, Result: HbA1c Laboratory Test")
	c.Assert(data.Encounter, DeepEquals, suite.Patient.Encounters[0].FHIRReference())
	c.Assert(data.ValueQuantity.Value, Equals, float64(8))
	c.Assert(data.ValueQuantity.Units, Equals, "%")
	c.Assert(data.AppliesPeriod.Start, DeepEquals, NewUnixTime(1320149800).FHIRDateTime())
	c.Assert(data.AppliesPeriod.End, DeepEquals, NewUnixTime(1320149800).FHIRDateTime())
	c.Assert(data.Name.Text, Equals, "Laboratory Test, Result: HbA1c Laboratory Test")
	c.Assert(data.Name.MatchesCode("http://loinc.org", "17856-6"), Equals, true)
}
