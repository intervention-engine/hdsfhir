package hdsfhir

import (
	"encoding/json"
	"io/ioutil"

	fhir "github.com/intervention-engine/fhir/models"
	"github.com/pebbe/util"
	. "gopkg.in/check.v1"
)

type MedicationSuite struct {
	Patient *Patient
}

var _ = Suite(&MedicationSuite{})

func (s *MedicationSuite) SetUpSuite(c *C) {
	data, err := ioutil.ReadFile("./fixtures/john_peters.json")
	util.CheckErr(err)

	s.Patient = &Patient{}
	err = json.Unmarshal(data, s.Patient)
	util.CheckErr(err)
}

func (s *MedicationSuite) TestMedicationFHIRModels(c *C) {
	models := s.Patient.Medications[0].FHIRModels()
	c.Assert(models, HasLen, 1)

	c.Assert(models[0], FitsTypeOf, &fhir.MedicationStatement{})
	medication := models[0].(*fhir.MedicationStatement)
	c.Assert(medication.Patient, DeepEquals, s.Patient.FHIRReference())
	c.Assert(medication.MedicationCodeableConcept.Text, Equals, "Medication, Order: BH Antidepressant medication (Code List: 2.16.840.1.113883.3.1257.1.972)")
	c.Assert(medication.MedicationCodeableConcept.MatchesCode("http://www.nlm.nih.gov/research/umls/rxnorm/", "1000048"), Equals, true)
	c.Assert(medication.EffectivePeriod.Start, DeepEquals, NewUnixTime(1349092800).FHIRDateTime())
	c.Assert(medication.EffectivePeriod.End, DeepEquals, NewUnixTime(1349092800).FHIRDateTime())
}

func (s *MedicationSuite) TestImmunizationFHIRModels(c *C) {
	models := s.Patient.Medications[1].FHIRModels()
	c.Assert(models, HasLen, 1)
	c.Assert(models[0], FitsTypeOf, &fhir.Immunization{})
	immunization := models[0].(*fhir.Immunization)
	c.Assert(immunization.Patient, DeepEquals, s.Patient.FHIRReference())
	c.Assert(immunization.Date, DeepEquals, NewUnixTime(1313409600).FHIRDateTime())
	c.Assert(immunization.VaccineType.Text, Equals, "Medication, Administered: Pneumococcal Vaccine (Code List: 2.16.840.1.113883.3.464.1003.110.12.1027)")
	c.Assert(immunization.VaccineType.MatchesCode("http://www2a.cdc.gov/vaccines/iis/iisstandards/vaccines.asp?rpt=cvx", "33"), Equals, true)
}
