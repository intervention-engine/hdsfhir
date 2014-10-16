package hdsfhir

import (
	"encoding/json"
	"github.com/pebbe/util"
	. "gopkg.in/check.v1"
	"io/ioutil"
)

type EncounterSuite struct {
	Encounter *Encounter
	Patient   *Patient
}

var _ = Suite(&EncounterSuite{})

func (s *EncounterSuite) SetUpSuite(c *C) {
	data, err := ioutil.ReadFile("./fixtures/john_peters.json")
	util.CheckErr(err)
	patient := &Patient{}
	err = json.Unmarshal(data, patient)
	s.Patient = patient
	util.CheckErr(err)
	s.Encounter = &patient.Encounters[0]
	s.Encounter.Patient = patient
}

func (s *EncounterSuite) TestToJSON(c *C) {
	data := s.Encounter.ToJSON()
	c.Assert(data, NotNil)
}
