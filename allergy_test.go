package hdsfhir

import (
	"encoding/json"
	"io/ioutil"

	fhir "github.com/intervention-engine/fhir/models"
	"github.com/pebbe/util"
	. "gopkg.in/check.v1"
)

type AllergySuite struct {
	Patient   *Patient
	Allergies map[string]*Allergy
}

var _ = Suite(&AllergySuite{})

func (s *AllergySuite) SetUpSuite(c *C) {
	data, err := ioutil.ReadFile("./fixtures/allergies.json")
	util.CheckErr(err)

	s.Allergies = make(map[string]*Allergy)
	err = json.Unmarshal(data, &s.Allergies)
	util.CheckErr(err)

	s.Patient = &Patient{}
	for _, allergy := range s.Allergies {
		allergy.Patient = s.Patient
	}
}

func (s *AllergySuite) TestActiveAllergy(c *C) {
	allergy := s.Allergies["active"].FHIRModels()[0].(*fhir.AllergyIntolerance)
	c.Assert(allergy.Onset, DeepEquals, NewUnixTime(1325396520).FHIRDateTime())
	c.Assert(allergy.Patient, DeepEquals, s.Patient.FHIRReference())
	c.Assert(allergy.Substance.Text, Equals, "Medication, Allergy: Influenza Vaccine (Code List: 2.16.840.1.113883.3.526.3.1254)")
	c.Assert(allergy.Substance.Coding, HasLen, 1)
	c.Assert(allergy.Substance.MatchesCode("http://www2a.cdc.gov/vaccines/iis/iisstandards/vaccines.asp?rpt=cvx", "111"), Equals, true)
	c.Assert(allergy.Status, Equals, "active")
	c.Assert(allergy.Criticality, Equals, "CRITH")
	c.Assert(allergy.Reaction, HasLen, 1)
	c.Assert(allergy.Reaction[0].Manifestation, HasLen, 1)
	c.Assert(allergy.Reaction[0].Manifestation[0].Text, Equals, "")
	c.Assert(allergy.Reaction[0].Manifestation[0].Coding, HasLen, 1)
	c.Assert(allergy.Reaction[0].Manifestation[0].MatchesCode("http://snomed.info/sct", "421581006"), Equals, true)
}

func (s *AllergySuite) TestDefaultStatusAllergy(c *C) {
	allergy := s.Allergies["defaultStatus"].FHIRModels()[0].(*fhir.AllergyIntolerance)
	c.Assert(allergy.Onset, DeepEquals, NewUnixTime(1325396520).FHIRDateTime())
	c.Assert(allergy.Patient, DeepEquals, s.Patient.FHIRReference())
	c.Assert(allergy.Substance.Text, Equals, "Medication, Allergy: Influenza Vaccine (Code List: 2.16.840.1.113883.3.526.3.1254)")
	c.Assert(allergy.Substance.Coding, HasLen, 1)
	c.Assert(allergy.Substance.MatchesCode("http://www2a.cdc.gov/vaccines/iis/iisstandards/vaccines.asp?rpt=cvx", "111"), Equals, true)
	c.Assert(allergy.Status, Equals, "active")
	c.Assert(allergy.Criticality, Equals, "")
	c.Assert(allergy.Reaction, IsNil)
}

func (s *AllergySuite) TestInactiveAllergy(c *C) {
	allergy := s.Allergies["inactive"].FHIRModels()[0].(*fhir.AllergyIntolerance)
	c.Assert(allergy.Onset, DeepEquals, NewUnixTime(1325396520).FHIRDateTime())
	c.Assert(allergy.Patient, DeepEquals, s.Patient.FHIRReference())
	c.Assert(allergy.Substance.Text, Equals, "Medication, Allergy: Influenza Vaccine (Code List: 2.16.840.1.113883.3.526.3.1254)")
	c.Assert(allergy.Substance.Coding, HasLen, 1)
	c.Assert(allergy.Substance.MatchesCode("http://www2a.cdc.gov/vaccines/iis/iisstandards/vaccines.asp?rpt=cvx", "111"), Equals, true)
	c.Assert(allergy.Status, Equals, "inactive")
	c.Assert(allergy.Criticality, Equals, "")
	c.Assert(allergy.Reaction, IsNil)
}

func (s *AllergySuite) TestResolvedAllergy(c *C) {
	allergy := s.Allergies["resolved"].FHIRModels()[0].(*fhir.AllergyIntolerance)
	c.Assert(allergy.Onset, DeepEquals, NewUnixTime(1325396520).FHIRDateTime())
	c.Assert(allergy.Patient, DeepEquals, s.Patient.FHIRReference())
	c.Assert(allergy.Substance.Text, Equals, "Medication, Allergy: Influenza Vaccine (Code List: 2.16.840.1.113883.3.526.3.1254)")
	c.Assert(allergy.Substance.Coding, HasLen, 1)
	c.Assert(allergy.Substance.MatchesCode("http://www2a.cdc.gov/vaccines/iis/iisstandards/vaccines.asp?rpt=cvx", "111"), Equals, true)
	c.Assert(allergy.Status, Equals, "resolved")
	c.Assert(allergy.Criticality, Equals, "")
	c.Assert(allergy.Reaction, IsNil)
}

func (s *AllergySuite) TestMildAllergy(c *C) {
	allergy := s.Allergies["mild"].FHIRModels()[0].(*fhir.AllergyIntolerance)
	c.Assert(allergy.Onset, DeepEquals, NewUnixTime(1325396520).FHIRDateTime())
	c.Assert(allergy.Patient, DeepEquals, s.Patient.FHIRReference())
	c.Assert(allergy.Substance.Text, Equals, "Medication, Allergy: Influenza Vaccine (Code List: 2.16.840.1.113883.3.526.3.1254)")
	c.Assert(allergy.Substance.Coding, HasLen, 1)
	c.Assert(allergy.Substance.MatchesCode("http://www2a.cdc.gov/vaccines/iis/iisstandards/vaccines.asp?rpt=cvx", "111"), Equals, true)
	c.Assert(allergy.Status, Equals, "active")
	c.Assert(allergy.Criticality, Equals, "CRITL")
	c.Assert(allergy.Reaction, IsNil)
}

func (s *AllergySuite) TestModerateAllergy(c *C) {
	allergy := s.Allergies["moderate"].FHIRModels()[0].(*fhir.AllergyIntolerance)
	c.Assert(allergy.Onset, DeepEquals, NewUnixTime(1325396520).FHIRDateTime())
	c.Assert(allergy.Patient, DeepEquals, s.Patient.FHIRReference())
	c.Assert(allergy.Substance.Text, Equals, "Medication, Allergy: Influenza Vaccine (Code List: 2.16.840.1.113883.3.526.3.1254)")
	c.Assert(allergy.Substance.Coding, HasLen, 1)
	c.Assert(allergy.Substance.MatchesCode("http://www2a.cdc.gov/vaccines/iis/iisstandards/vaccines.asp?rpt=cvx", "111"), Equals, true)
	c.Assert(allergy.Status, Equals, "active")
	c.Assert(allergy.Criticality, Equals, "CRITU")
	c.Assert(allergy.Reaction, IsNil)
}

func (s *AllergySuite) TestNegatedAllergy(c *C) {
	allergy := s.Allergies["negated"].FHIRModels()[0].(*fhir.AllergyIntolerance)
	c.Assert(allergy.Onset, DeepEquals, NewUnixTime(1325396520).FHIRDateTime())
	c.Assert(allergy.Patient, DeepEquals, s.Patient.FHIRReference())
	c.Assert(allergy.Substance.Text, Equals, "Medication, Allergy: Influenza Vaccine (Code List: 2.16.840.1.113883.3.526.3.1254)")
	c.Assert(allergy.Substance.Coding, HasLen, 1)
	c.Assert(allergy.Substance.MatchesCode("http://www2a.cdc.gov/vaccines/iis/iisstandards/vaccines.asp?rpt=cvx", "111"), Equals, true)
	c.Assert(allergy.Status, Equals, "refuted")
	c.Assert(allergy.Criticality, Equals, "")
	c.Assert(allergy.Reaction, IsNil)
}
