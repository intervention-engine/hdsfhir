package models

import (
	"encoding/json"
	"github.com/pebbe/util"
	. "gopkg.in/check.v1"
	"io/ioutil"
	"strconv"
	"time"
)

type VitalSignSuite struct {
	VitalSign *VitalSign
	Patient   *Patient
}

var _ = Suite(&VitalSignSuite{})

func (suite *VitalSignSuite) SetUpSuite(c *C) {
	data, err := ioutil.ReadFile("../fixtures/john_peters.json")
	util.CheckErr(err)
	patient := &Patient{}
	err = json.Unmarshal(data, patient)
	patient.ServerURL = "http://www.example.com/Patient/1"
	suite.Patient = patient
	for i, encounter := range patient.Encounters {
		encounter.ServerURL = "http://www.example.com/Encounter/" + strconv.Itoa(i)
	}
	util.CheckErr(err)
	suite.VitalSign = patient.VitalSigns[0]
	suite.VitalSign.Patient = patient
}

func (suite *VitalSignSuite) TestToJSON(c *C) {
	data := suite.VitalSign.ToJSON()
	c.Assert(data, NotNil)
}

func (suite *VitalSignSuite) TestFHIRModel(c *C) {
	data := suite.VitalSign.FHIRModel()
	c.Assert(data.Subject.Reference, Equals, suite.Patient.ServerURL)
	c.Assert(data.Name.Text, Equals, suite.VitalSign.Description)
	c.Assert(data.Encounter.Reference, Equals, "http://www.example.com/Encounter/0")
	c.Assert(data.ValueQuantity.Value, Equals, float64(8))
	c.Assert(data.ValueQuantity.Units, Equals, "%")
	c.Assert(data.AppliesPeriod.Start.Time, Equals, time.Unix(1320149800, 0))
	c.Assert(data.AppliesPeriod.End.Time, Equals, time.Unix(1320149800, 0))
	c.Assert(data.Name.Coding[0].Code, Equals, "17856-6")
	c.Assert(data.Name.Coding[0].System, Equals, "http://loinc.org")
}
