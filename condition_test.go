package hdsfhir

import (
	"encoding/json"
	"github.com/pebbe/util"
	. "gopkg.in/check.v1"
	"io/ioutil"
)

type ConditionSuite struct {
	Condition *Condition
	Patient   *Patient
}

var _ = Suite(&EncounterSuite{})

func (s *ConditionSuite) SetUpSuite(c *C) {
	data, err := ioutil.ReadFile("./fixtures/john_peters.json")
	util.CheckErr(err)
	patient := &Patient{}
	err = json.Unmarshal(data, patient)
	s.Patient = patient
	util.CheckErr(err)
	s.Condition = &patient.Conditions[0]
	s.Condition.Patient = patient
}

func (s *ConditionSuite) TestToJSON(c *C) {
	data := s.Condition.ToJSON()
	c.Assert(data, NotNil)
}
