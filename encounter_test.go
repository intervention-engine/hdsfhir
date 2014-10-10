package hdsfhir

import (
	"encoding/json"
	"github.com/pebbe/util"
	. "gopkg.in/check.v1"
	"io/ioutil"
)

type EncounterSuite struct {
	Encounter *Entry
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
}

func (s *EncounterSuite) TestToJSON(c *C) {
	data := EncounterToJSON(s.Patient, s.Encounter)
	c.Assert(data, NotNil)
}
