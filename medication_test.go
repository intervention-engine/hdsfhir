package hdsfhir

import (
	"encoding/json"
	"io/ioutil"

	fhir "github.com/intervention-engine/fhir/models"
	"github.com/pebbe/util"
	. "gopkg.in/check.v1"
)

type MedicationSuite struct {
	Patient     *Patient
	Medications map[string]*Medication
}

var _ = Suite(&MedicationSuite{})

func (s *MedicationSuite) SetUpSuite(c *C) {
	data, err := ioutil.ReadFile("./fixtures/medications.json")
	util.CheckErr(err)

	s.Medications = make(map[string]*Medication)
	err = json.Unmarshal(data, &s.Medications)
	util.CheckErr(err)

	s.Patient = &Patient{}
	for _, medication := range s.Medications {
		medication.Patient = s.Patient
	}
}

func (s *MedicationSuite) TestMedicationOrdered(c *C) {
	medication := s.Medications["medicationOrdered"].FHIRModels()[0].(*fhir.MedicationStatement)
	c.Assert(medication.Patient, DeepEquals, s.Patient.FHIRReference())
	c.Assert(medication.Status, Equals, "intended")
	c.Assert(medication.WasNotTaken, IsNil)
	c.Assert(medication.ReasonNotTaken, IsNil)
	c.Assert(medication.EffectivePeriod.Start, DeepEquals, NewUnixTime(1349092800).FHIRDateTime())
	c.Assert(medication.EffectivePeriod.End, DeepEquals, NewUnixTime(1349092800).FHIRDateTime())
	c.Assert(medication.MedicationCodeableConcept.Text, Equals, "Medication, Order: BH Antidepressant medication (Code List: 2.16.840.1.113883.3.1257.1.972)")
	c.Assert(medication.MedicationCodeableConcept.Coding, HasLen, 1)
	c.Assert(medication.MedicationCodeableConcept.MatchesCode("http://www.nlm.nih.gov/research/umls/rxnorm/", "1000048"), Equals, true)
}

func (s *MedicationSuite) TestMedicationDispensed(c *C) {
	medication := s.Medications["medicationDispensed"].FHIRModels()[0].(*fhir.MedicationStatement)
	c.Assert(medication.Patient, DeepEquals, s.Patient.FHIRReference())
	c.Assert(medication.Status, Equals, "intended")
	c.Assert(medication.WasNotTaken, IsNil)
	c.Assert(medication.ReasonNotTaken, IsNil)
	c.Assert(medication.EffectivePeriod.Start, DeepEquals, NewUnixTime(1349092800).FHIRDateTime())
	c.Assert(medication.EffectivePeriod.End, DeepEquals, NewUnixTime(1349092800).FHIRDateTime())
	c.Assert(medication.MedicationCodeableConcept.Text, Equals, "Medication, Dispensed: BH Antidepressant medication (Code List: 2.16.840.1.113883.3.1257.1.972)")
	c.Assert(medication.MedicationCodeableConcept.Coding, HasLen, 1)
	c.Assert(medication.MedicationCodeableConcept.MatchesCode("http://www.nlm.nih.gov/research/umls/rxnorm/", "1000048"), Equals, true)
}

func (s *MedicationSuite) TestMedicationNotOrdered(c *C) {
	medication := s.Medications["medicationNotOrdered"].FHIRModels()[0].(*fhir.MedicationStatement)
	c.Assert(medication.Patient, DeepEquals, s.Patient.FHIRReference())
	c.Assert(medication.Status, Equals, "intended")
	c.Assert(*medication.WasNotTaken, Equals, true)
	c.Assert(medication.ReasonNotTaken, HasLen, 1)
	c.Assert(medication.ReasonNotTaken[0].Text, Equals, "")
	c.Assert(medication.ReasonNotTaken[0].Coding, HasLen, 1)
	c.Assert(medication.ReasonNotTaken[0].MatchesCode("http://snomed.info/sct", "416098002"), Equals, true)
	c.Assert(medication.EffectivePeriod.Start, DeepEquals, NewUnixTime(1349092800).FHIRDateTime())
	c.Assert(medication.EffectivePeriod.End, DeepEquals, NewUnixTime(1349092800).FHIRDateTime())
	c.Assert(medication.MedicationCodeableConcept.Text, Equals, "Medication, Order: BH Antidepressant medication (Code List: 2.16.840.1.113883.3.1257.1.972)")
	c.Assert(medication.MedicationCodeableConcept.Coding, HasLen, 1)
	c.Assert(medication.MedicationCodeableConcept.MatchesCode("http://www.nlm.nih.gov/research/umls/rxnorm/", "1000048"), Equals, true)
}

func (s *MedicationSuite) TestImmunizationAdministered(c *C) {
	immunization := s.Medications["immunizationAdministered"].FHIRModels()[0].(*fhir.Immunization)
	c.Assert(immunization.Status, Equals, "completed")
	c.Assert(immunization.Date, DeepEquals, NewUnixTime(1313409600).FHIRDateTime())
	c.Assert(immunization.VaccineCode.Text, Equals, "Immunization, Administered: Pneumococcal Vaccine (Code List: 2.16.840.1.113883.3.464.1003.110.12.1027)")
	c.Assert(immunization.VaccineCode.Coding, HasLen, 1)
	c.Assert(immunization.VaccineCode.MatchesCode("http://www2a.cdc.gov/vaccines/iis/iisstandards/vaccines.asp?rpt=cvx", "33"), Equals, true)
	c.Assert(immunization.Patient, DeepEquals, s.Patient.FHIRReference())
	c.Assert(immunization.WasNotGiven, IsNil)
	c.Assert(immunization.Explanation, IsNil)
}

func (s *MedicationSuite) TestImmunizationNotAdministered(c *C) {
	immunization := s.Medications["immunizationNotAdministered"].FHIRModels()[0].(*fhir.Immunization)
	c.Assert(immunization.Status, Equals, "completed")
	c.Assert(immunization.Date, DeepEquals, NewUnixTime(1313409600).FHIRDateTime())
	c.Assert(immunization.VaccineCode.Text, Equals, "Immunization, Administered: Pneumococcal Vaccine (Code List: 2.16.840.1.113883.3.464.1003.110.12.1027)")
	c.Assert(immunization.VaccineCode.Coding, HasLen, 1)
	c.Assert(immunization.VaccineCode.MatchesCode("http://www2a.cdc.gov/vaccines/iis/iisstandards/vaccines.asp?rpt=cvx", "33"), Equals, true)
	c.Assert(immunization.Patient, DeepEquals, s.Patient.FHIRReference())
	c.Assert(*immunization.WasNotGiven, Equals, true)
	c.Assert(immunization.Explanation.Reason, IsNil)
	c.Assert(immunization.Explanation.ReasonNotGiven, HasLen, 1)
	c.Assert(immunization.Explanation.ReasonNotGiven[0].Text, Equals, "")
	c.Assert(immunization.Explanation.ReasonNotGiven[0].Coding, HasLen, 1)
	c.Assert(immunization.Explanation.ReasonNotGiven[0].MatchesCode("http://snomed.info/sct", "591000119102"), Equals, true)
}
