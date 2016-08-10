package hdsfhir

import (
	"encoding/json"
	"io/ioutil"

	fhir "github.com/intervention-engine/fhir/models"
	"github.com/pebbe/util"
	. "gopkg.in/check.v1"
)

type VitalSignSuite struct {
	Patient    *Patient
	VitalSigns map[string]*VitalSign
	Encounter  *Encounter
}

var _ = Suite(&VitalSignSuite{})

func (s *VitalSignSuite) SetUpSuite(c *C) {
	data, err := ioutil.ReadFile("./fixtures/vital_signs.json")
	util.CheckErr(err)

	s.VitalSigns = make(map[string]*VitalSign)
	err = json.Unmarshal(data, &s.VitalSigns)
	util.CheckErr(err)

	s.Patient = &Patient{}
	s.Encounter = &Encounter{Entry: Entry{StartTime: NewUnixTime(1320148800), EndTime: NewUnixTime(1320152400)}}
	s.Patient.Encounters = []*Encounter{s.Encounter}
	for _, vital := range s.VitalSigns {
		vital.Patient = s.Patient
	}
}

func (s *VitalSignSuite) TestFHIRModels(c *C) {
	models := s.VitalSigns["hba1c"].FHIRModels()
	c.Assert(models, HasLen, 1)
	c.Assert(models[0], FitsTypeOf, &fhir.Observation{})

	data := models[0].(*fhir.Observation)
	c.Assert(data.Subject, DeepEquals, s.Patient.FHIRReference())
	c.Assert(data.Code.Text, Equals, "Laboratory Test, Result: HbA1c Laboratory Test")
	c.Assert(data.Code.Coding, HasLen, 1)
	c.Assert(data.Code.MatchesCode("http://loinc.org", "17856-6"), Equals, true)
	c.Assert(data.Encounter, DeepEquals, s.Encounter.FHIRReference())
	c.Assert(*data.ValueQuantity.Value, Equals, float64(8))
	c.Assert(data.ValueQuantity.Unit, Equals, "%")
	c.Assert(data.Interpretation.Text, Equals, "")
	c.Assert(data.Interpretation.Coding, HasLen, 1)
	c.Assert(data.Interpretation.MatchesCode("urn:oid:2.16.840.1.113883.1.11.78", "A"), Equals, true)
	c.Assert(data.EffectivePeriod.Start, DeepEquals, NewUnixTime(1320149800).FHIRDateTime())
	c.Assert(data.EffectivePeriod.End, DeepEquals, NewUnixTime(1320149800).FHIRDateTime())
}
