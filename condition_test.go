package hdsfhir

import (
	"encoding/json"
	"io/ioutil"

	fhir "github.com/intervention-engine/fhir/models"
	"github.com/pebbe/util"
	. "gopkg.in/check.v1"
)

type ConditionSuite struct {
	Patient    *Patient
	Conditions map[string]*Condition
}

var _ = Suite(&ConditionSuite{})

func (s *ConditionSuite) SetUpSuite(c *C) {
	data, err := ioutil.ReadFile("./fixtures/conditions.json")
	util.CheckErr(err)

	s.Conditions = make(map[string]*Condition)
	err = json.Unmarshal(data, &s.Conditions)
	util.CheckErr(err)

	s.Patient = &Patient{}
	for _, condition := range s.Conditions {
		condition.Patient = s.Patient
	}
}

func (s *ConditionSuite) TestActiveCondition(c *C) {
	condition := s.Conditions["active"].FHIRModels()[0].(*fhir.Condition)
	c.Assert(condition.Patient, DeepEquals, s.Patient.FHIRReference())
	c.Assert(condition.Code.Text, Equals, "Diagnosis, Active: Heart Failure (Code List: 2.16.840.1.113883.3.526.3.376)")
	c.Assert(condition.Code.Coding, HasLen, 3)
	c.Assert(condition.Code.MatchesCode("http://snomed.info/sct", "10091002"), Equals, true)
	c.Assert(condition.Code.MatchesCode("http://hl7.org/fhir/sid/icd-9", "428.0"), Equals, true)
	c.Assert(condition.Code.MatchesCode("http://hl7.org/fhir/sid/icd-10", "I50.1"), Equals, true)
	c.Assert(condition.ClinicalStatus, Equals, "active")
	c.Assert(condition.VerificationStatus, Equals, "confirmed")
	c.Assert(condition.Severity.MatchesCode("http://snomed.info/sct", "24484000"), Equals, true)
	c.Assert(condition.Severity.Text, Equals, "Severe")
	c.Assert(condition.OnsetDateTime, DeepEquals, NewUnixTime(1330603200).FHIRDateTime())
	c.Assert(condition.AbatementDateTime, IsNil)
}

func (s *ConditionSuite) TestInactiveCondition(c *C) {
	condition := s.Conditions["inactive"].FHIRModels()[0].(*fhir.Condition)
	c.Assert(condition.Patient, DeepEquals, s.Patient.FHIRReference())
	c.Assert(condition.Code.Text, Equals, "Diagnosis, Inactive: Heart Failure (Code List: 2.16.840.1.113883.3.526.3.376)")
	c.Assert(condition.Code.MatchesCode("http://snomed.info/sct", "10091002"), Equals, true)
	c.Assert(condition.Code.MatchesCode("http://hl7.org/fhir/sid/icd-9", "428.0"), Equals, true)
	c.Assert(condition.Code.MatchesCode("http://hl7.org/fhir/sid/icd-10", "I50.1"), Equals, true)
	c.Assert(condition.ClinicalStatus, Equals, "remission")
	c.Assert(condition.VerificationStatus, Equals, "confirmed")
	c.Assert(condition.Severity, IsNil)
	c.Assert(condition.OnsetDateTime, DeepEquals, NewUnixTime(1330603200).FHIRDateTime())
	c.Assert(condition.AbatementDateTime, DeepEquals, NewUnixTime(1330624800).FHIRDateTime())
}

func (s *ConditionSuite) TestResolvedCondition(c *C) {
	condition := s.Conditions["resolved"].FHIRModels()[0].(*fhir.Condition)
	c.Assert(condition.Patient, DeepEquals, s.Patient.FHIRReference())
	c.Assert(condition.Code.Text, Equals, "Diagnosis, Resolved: Heart Failure (Code List: 2.16.840.1.113883.3.526.3.376)")
	c.Assert(condition.Code.MatchesCode("http://snomed.info/sct", "10091002"), Equals, true)
	c.Assert(condition.Code.MatchesCode("http://hl7.org/fhir/sid/icd-9", "428.0"), Equals, true)
	c.Assert(condition.Code.MatchesCode("http://hl7.org/fhir/sid/icd-10", "I50.1"), Equals, true)
	c.Assert(condition.ClinicalStatus, Equals, "resolved")
	c.Assert(condition.VerificationStatus, Equals, "confirmed")
	c.Assert(condition.Severity, IsNil)
	c.Assert(condition.OnsetDateTime, DeepEquals, NewUnixTime(1330603200).FHIRDateTime())
	c.Assert(condition.AbatementDateTime, DeepEquals, NewUnixTime(1330624800).FHIRDateTime())
}

func (s *ConditionSuite) TestUnsetStatusWithAbatementCondition(c *C) {
	condition := s.Conditions["unsetStatusAbated"].FHIRModels()[0].(*fhir.Condition)
	c.Assert(condition.Patient, DeepEquals, s.Patient.FHIRReference())
	c.Assert(condition.Code.Text, Equals, "Diagnosis, Resolved: Heart Failure (Code List: 2.16.840.1.113883.3.526.3.376)")
	c.Assert(condition.Code.MatchesCode("http://snomed.info/sct", "10091002"), Equals, true)
	c.Assert(condition.Code.MatchesCode("http://hl7.org/fhir/sid/icd-9", "428.0"), Equals, true)
	c.Assert(condition.Code.MatchesCode("http://hl7.org/fhir/sid/icd-10", "I50.1"), Equals, true)
	c.Assert(condition.ClinicalStatus, Equals, "resolved")
	c.Assert(condition.VerificationStatus, Equals, "confirmed")
	c.Assert(condition.Severity, IsNil)
	c.Assert(condition.OnsetDateTime, DeepEquals, NewUnixTime(1330603200).FHIRDateTime())
	c.Assert(condition.AbatementDateTime, DeepEquals, NewUnixTime(1330624800).FHIRDateTime())
}

func (s *ConditionSuite) TestActiveStatusWithAbatementCondition(c *C) {
	condition := s.Conditions["activeStatusAbated"].FHIRModels()[0].(*fhir.Condition)
	c.Assert(condition.Patient, DeepEquals, s.Patient.FHIRReference())
	c.Assert(condition.Code.Text, Equals, "Diagnosis, Inactive: Heart Failure (Code List: 2.16.840.1.113883.3.526.3.376)")
	c.Assert(condition.Code.MatchesCode("http://snomed.info/sct", "10091002"), Equals, true)
	c.Assert(condition.Code.MatchesCode("http://hl7.org/fhir/sid/icd-9", "428.0"), Equals, true)
	c.Assert(condition.Code.MatchesCode("http://hl7.org/fhir/sid/icd-10", "I50.1"), Equals, true)
	c.Assert(condition.ClinicalStatus, Equals, "remission")
	c.Assert(condition.VerificationStatus, Equals, "confirmed")
	c.Assert(condition.Severity, IsNil)
	c.Assert(condition.OnsetDateTime, DeepEquals, NewUnixTime(1330603200).FHIRDateTime())
	c.Assert(condition.AbatementDateTime, DeepEquals, NewUnixTime(1330624800).FHIRDateTime())
}

func (s *ConditionSuite) TestNegatedCondition(c *C) {
	condition := s.Conditions["negated"].FHIRModels()[0].(*fhir.Condition)
	c.Assert(condition.Patient, DeepEquals, s.Patient.FHIRReference())
	c.Assert(condition.Code.Text, Equals, "Diagnosis, Active: Heart Failure (Code List: 2.16.840.1.113883.3.526.3.376)")
	c.Assert(condition.Code.MatchesCode("http://snomed.info/sct", "10091002"), Equals, true)
	c.Assert(condition.Code.MatchesCode("http://hl7.org/fhir/sid/icd-9", "428.0"), Equals, true)
	c.Assert(condition.Code.MatchesCode("http://hl7.org/fhir/sid/icd-10", "I50.1"), Equals, true)
	c.Assert(condition.ClinicalStatus, Equals, "")
	c.Assert(condition.VerificationStatus, Equals, "refuted")
	c.Assert(condition.Severity, IsNil)
	c.Assert(condition.OnsetDateTime, IsNil)
	c.Assert(condition.AbatementDateTime, IsNil)
}
