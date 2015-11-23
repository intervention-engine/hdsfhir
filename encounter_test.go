package hdsfhir

import (
	"encoding/json"
	"io/ioutil"

	fhir "github.com/intervention-engine/fhir/models"
	"github.com/pebbe/util"
	. "gopkg.in/check.v1"
)

type EncounterSuite struct {
	Patient    *Patient
	Encounters map[string]*Encounter
}

var _ = Suite(&EncounterSuite{})

func (s *EncounterSuite) SetUpSuite(c *C) {
	data, err := ioutil.ReadFile("./fixtures/encounters.json")
	util.CheckErr(err)

	s.Encounters = make(map[string]*Encounter)
	err = json.Unmarshal(data, &s.Encounters)
	util.CheckErr(err)

	s.Patient = &Patient{}
	for _, encounter := range s.Encounters {
		encounter.Patient = s.Patient
	}
}

func (s *EncounterSuite) TestOfficeVisit(c *C) {
	encounter := s.Encounters["officeVisit"].FHIRModels()[0].(*fhir.Encounter)
	c.Assert(encounter.Status, Equals, "finished")
	c.Assert(encounter.Type, HasLen, 1)
	c.Assert(encounter.Type[0].Text, Equals, "Encounter, Performed: Office Visit (Code List: 2.16.840.1.113883.3.464.1003.101.12.1001)")
	c.Assert(encounter.Type[0].Coding, HasLen, 1)
	c.Assert(encounter.Type[0].MatchesCode("http://www.ama-assn.org/go/cpt", "99201"), Equals, true)
	c.Assert(encounter.Patient, DeepEquals, s.Patient.FHIRReference())
	c.Assert(encounter.Period.Start, DeepEquals, NewUnixTime(1320148800).FHIRDateTime())
	c.Assert(encounter.Period.End, DeepEquals, NewUnixTime(1320152400).FHIRDateTime())
	c.Assert(encounter.Reason, IsNil)
	c.Assert(encounter.Hospitalization, IsNil)
}

func (s *EncounterSuite) TestOrderedOfficeVisit(c *C) {
	encounter := s.Encounters["orderedOfficeVisit"].FHIRModels()[0].(*fhir.Encounter)
	c.Assert(encounter.Status, Equals, "planned")
	c.Assert(encounter.Type, HasLen, 1)
	c.Assert(encounter.Type[0].Text, Equals, "Encounter, Order: Office Visit (Code List: 2.16.840.1.113883.3.464.1003.101.12.1001)")
	c.Assert(encounter.Type[0].MatchesCode("http://www.ama-assn.org/go/cpt", "99201"), Equals, true)
	c.Assert(encounter.Patient, DeepEquals, s.Patient.FHIRReference())
	c.Assert(encounter.Period.Start, DeepEquals, NewUnixTime(1320148800).FHIRDateTime())
	c.Assert(encounter.Period.End, DeepEquals, NewUnixTime(1320148800).FHIRDateTime())
	c.Assert(encounter.Reason, HasLen, 1)
	c.Assert(encounter.Reason[0].Text, Equals, "")
	c.Assert(encounter.Reason[0].Coding, HasLen, 1)
	c.Assert(encounter.Reason[0].MatchesCode("http://snomed.info/sct", "2070002"), Equals, true)
	c.Assert(encounter.Hospitalization, IsNil)
}

func (s *EncounterSuite) TestInpatientEncounter(c *C) {
	encounter := s.Encounters["inpatientEncounter"].FHIRModels()[0].(*fhir.Encounter)
	c.Assert(encounter.Status, Equals, "finished")
	c.Assert(encounter.Type, HasLen, 1)
	c.Assert(encounter.Type[0].Text, Equals, "Encounter, Performed: Inpatient Encounter")
	c.Assert(encounter.Type[0].Coding, HasLen, 1)
	c.Assert(encounter.Type[0].MatchesCode("http://www.ama-assn.org/go/cpt", "99235"), Equals, true)
	c.Assert(encounter.Patient, DeepEquals, s.Patient.FHIRReference())
	c.Assert(encounter.Period.Start, DeepEquals, NewUnixTime(1320148800).FHIRDateTime())
	c.Assert(encounter.Period.End, DeepEquals, NewUnixTime(1320152400).FHIRDateTime())
	c.Assert(encounter.Reason, IsNil)
	c.Assert(encounter.Hospitalization.DischargeDisposition.Coding, HasLen, 1)
	c.Assert(encounter.Hospitalization.DischargeDisposition.MatchesCode("urn:oid:2.16.840.1.113883.3.88.12.80.33", "1"), Equals, true)
}
