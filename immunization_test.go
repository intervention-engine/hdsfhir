package hdsfhir

import (
	"encoding/json"
	"io/ioutil"

	fhir "github.com/intervention-engine/fhir/models"
	"github.com/pebbe/util"
	. "gopkg.in/check.v1"
)

type ImmunizationSuite struct {
	Patient       *Patient
	Immunizations map[string]*Immunization
}

var _ = Suite(&ImmunizationSuite{})

func (s *ImmunizationSuite) SetUpSuite(c *C) {
	data, err := ioutil.ReadFile("./fixtures/immunizations.json")
	util.CheckErr(err)

	s.Immunizations = make(map[string]*Immunization)
	err = json.Unmarshal(data, &s.Immunizations)
	util.CheckErr(err)

	s.Patient = &Patient{}
	for _, immunization := range s.Immunizations {
		immunization.Patient = s.Patient
	}
}

func (s *ImmunizationSuite) TestImmunizationAdministered(c *C) {
	immunization := s.Immunizations["immunizationAdministered"].FHIRModels()[0].(*fhir.Immunization)
	c.Assert(immunization.Status, Equals, "completed")
	c.Assert(immunization.Date, DeepEquals, NewUnixTime(1263168508).FHIRDateTime())
	c.Assert(immunization.VaccineCode.Text, Equals, "MMR")
	c.Assert(immunization.VaccineCode.Coding, HasLen, 1)
	c.Assert(immunization.VaccineCode.MatchesCode("http://www2a.cdc.gov/vaccines/iis/iisstandards/vaccines.asp?rpt=cvx", "03"), Equals, true)
	c.Assert(immunization.Patient, DeepEquals, s.Patient.FHIRReference())
	c.Assert(immunization.VaccinationProtocol, HasLen, 1)
	c.Assert(*immunization.VaccinationProtocol[0].DoseSequence, Equals, uint32(2))
	c.Assert(immunization.WasNotGiven, IsNil)
	c.Assert(immunization.Explanation, IsNil)
}

func (s *ImmunizationSuite) TestImmunizationNotAdministered(c *C) {
	immunization := s.Immunizations["immunizationNotAdministered"].FHIRModels()[0].(*fhir.Immunization)
	c.Assert(immunization.Status, Equals, "completed")
	c.Assert(immunization.Date, DeepEquals, NewUnixTime(1263168508).FHIRDateTime())
	c.Assert(immunization.VaccineCode.Text, Equals, "MMR")
	c.Assert(immunization.VaccineCode.Coding, HasLen, 1)
	c.Assert(immunization.VaccineCode.MatchesCode("http://www2a.cdc.gov/vaccines/iis/iisstandards/vaccines.asp?rpt=cvx", "03"), Equals, true)
	c.Assert(immunization.Patient, DeepEquals, s.Patient.FHIRReference())
	c.Assert(*immunization.WasNotGiven, Equals, true)
	c.Assert(immunization.Explanation.Reason, IsNil)
	c.Assert(immunization.Explanation.ReasonNotGiven, HasLen, 1)
	c.Assert(immunization.Explanation.ReasonNotGiven[0].Text, Equals, "")
	c.Assert(immunization.Explanation.ReasonNotGiven[0].Coding, HasLen, 1)
	c.Assert(immunization.Explanation.ReasonNotGiven[0].MatchesCode("http://snomed.info/sct", "591000119102"), Equals, true)
	c.Assert(immunization.VaccinationProtocol, IsNil)
}
