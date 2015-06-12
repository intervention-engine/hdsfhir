package hdsfhir

import (
	"encoding/json"
	"io/ioutil"

	fhir "github.com/intervention-engine/fhir/models"
	"github.com/pebbe/util"
	. "gopkg.in/check.v1"
)

type ConditionSuite struct {
	Patient *Patient
}

var _ = Suite(&ConditionSuite{})

func (s *ConditionSuite) SetUpSuite(c *C) {
	data, err := ioutil.ReadFile("./fixtures/john_peters.json")
	util.CheckErr(err)

	s.Patient = &Patient{}
	err = json.Unmarshal(data, s.Patient)
	util.CheckErr(err)
}

func (s *ConditionSuite) TestFHIRModels(c *C) {
	models := s.Patient.Conditions[0].FHIRModels()
	c.Assert(models, HasLen, 1)
	c.Assert(models[0], FitsTypeOf, &fhir.Condition{})
	condition := models[0].(*fhir.Condition)
	c.Assert(condition.Subject, DeepEquals, s.Patient.FHIRReference())
	c.Assert(condition.OnsetDate, DeepEquals, NewUnixTime(1330603200).FHIRDateTime())
	c.Assert(condition.AbatementDate, IsNil)
	c.Assert(condition.Code.Text, Equals, "Diagnosis, Active: Heart Failure (Code List: 2.16.840.1.113883.3.526.3.376)")
	c.Assert(condition.Code.MatchesCode("http://snomed.info/sct", "10091002"), Equals, true)
	c.Assert(condition.Code.MatchesCode("http://hl7.org/fhir/sid/icd-9", "428.0"), Equals, true)
	c.Assert(condition.Code.MatchesCode("http://hl7.org/fhir/sid/icd-10", "I50.1"), Equals, true)
}
