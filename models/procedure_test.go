package models

import (
	"encoding/json"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/pebbe/util"
	. "gopkg.in/check.v1"
)

type ProcedureSuite struct {
	Procedures []*Procedure
	Patient    *Patient
}

var _ = Suite(&ProcedureSuite{})

func (suite *ProcedureSuite) SetUpSuite(c *C) {
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
	suite.Procedures = patient.Procedures
	for _, procedure := range patient.Procedures {
		procedure.Patient = patient
	}
}

func (suite *ProcedureSuite) TestToJSON(c *C) {
	for _, procedure := range suite.Procedures {
		data := procedure.ToJSON()
		c.Assert(data, NotNil)
	}
}

func (suite *ProcedureSuite) TestToToFhirModel(c *C) {
	procedure1 := suite.Procedures[0].ToFhirModel()

	c.Assert(procedure1.Subject.Reference, Equals, suite.Patient.ServerURL)
	c.Assert(procedure1.Type.Coding[0].System, Equals, "http://snomed.info/sct")
	c.Assert(procedure1.Type.Coding[0].Code, Equals, "116783008")
	c.Assert(procedure1.Date.Start.Time, Equals, time.Unix(1320149800, 0))
	c.Assert(procedure1.Date.End.Time, Equals, time.Unix(1320159800, 0))
	c.Assert(procedure1.Encounter.Reference, Equals, "http://www.example.com/Encounter/0")
	c.Assert(procedure1.Notes, Equals, "Procedure, Result: Clinical Staging Procedure")

	procedure2 := suite.Procedures[1].ToFhirModel()
	c.Assert(procedure2.Subject.Reference, Equals, suite.Patient.ServerURL)

	icd9, icd10, snomed := false, false, false
	for _, code := range procedure2.Type.Coding {
		switch code.System {
		case "http://hl7.org/fhir/sid/icd-10":
			c.Assert(code.Code, Equals, "0210093")
			icd10 = true
		case "http://snomed.info/sct":
			c.Assert(code.Code, Equals, "10190003")
			snomed = true
		case "http://hl7.org/fhir/sid/icd-9":
			c.Assert(code.Code, Equals, "36.10")
			icd9 = true
		}
	}
	c.Assert(icd9 && icd10 && snomed, Equals, true)
	c.Assert(procedure2.Date.Start.Time, Equals, time.Unix(1362239100, 0))
	c.Assert(procedure2.Date.End.Time, Equals, time.Unix(1362242700, 0))
	c.Assert(procedure2.Encounter.Reference, Equals, "")
	c.Assert(procedure2.Notes, Equals, "Procedure, Performed: Hospital measures-CABG")

}

func (suite *ProcedureSuite) TestResultHandling(c *C) {
	procedure1 := suite.Procedures[0]
	procedure1.ProcessResultObservations()
	c.Assert(len(procedure1.ResultObservations), Equals, 3)
	procedure1.ProcessResultReport()
	c.Assert(procedure1.Report, NotNil)

}
