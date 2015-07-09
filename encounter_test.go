package hdsfhir

import (
	"encoding/json"
	"io/ioutil"

	fhir "github.com/intervention-engine/fhir/models"
	"github.com/pebbe/util"
	. "gopkg.in/check.v1"
)

type EncounterSuite struct {
	Patient *Patient
}

var _ = Suite(&EncounterSuite{})

func (s *EncounterSuite) SetUpSuite(c *C) {
	data, err := ioutil.ReadFile("./fixtures/john_peters.json")
	util.CheckErr(err)

	s.Patient = &Patient{}
	err = json.Unmarshal(data, s.Patient)
	util.CheckErr(err)
}

func (s *EncounterSuite) TestFHIRModels(c *C) {
	models := s.Patient.Encounters[0].FHIRModels()
	c.Assert(models, HasLen, 1)
	c.Assert(models[0], FitsTypeOf, &fhir.Encounter{})
	encounter := models[0].(*fhir.Encounter)
	c.Assert(encounter.Patient, DeepEquals, s.Patient.FHIRReference())
	c.Assert(encounter.Period.Start, DeepEquals, NewUnixTime(1320148800).FHIRDateTime())
	c.Assert(encounter.Period.End, DeepEquals, NewUnixTime(1320152400).FHIRDateTime())
	c.Assert(encounter.Type, HasLen, 1)
	c.Assert(encounter.Type[0].Text, Equals, "Encounter, Performed: Office Visit (Code List: 2.16.840.1.113883.3.464.1003.101.12.1001)")
	c.Assert(encounter.Type[0].MatchesCode("http://www.ama-assn.org/go/cpt", "99201"), Equals, true)
}
