package hdsfhir

import (
	"encoding/json"
	"github.com/pebbe/util"
	. "gopkg.in/check.v1"
	"io/ioutil"
)

type MedicationSuite struct {
	Medication *Medication
	Patient    *Patient
}

var _ = Suite(&EncounterSuite{})

func (s *MedicationSuite) SetUpSuite(c *C) {
	data, err := ioutil.ReadFile("./fixtures/john_peters.json")
	util.CheckErr(err)
	patient := &Patient{}
	err = json.Unmarshal(data, patient)
	s.Patient = patient
	util.CheckErr(err)
	s.Medication = patient.Medications[0]
	s.Medication.Patient = patient
}

func (s *MedicationSuite) TestNewMedicationWrapper(c *C) {
	wrapper := NewMedicationWrapper(*s.Medication)
	c.Assert(wrapper.Medication.Code.Coding[0].Code, Equals, "1000048")
	c.Assert(wrapper.Medication.Code.Coding[0].System, Equals, "http://www.nlm.nih.gov/research/umls/rxnorm/")
}
