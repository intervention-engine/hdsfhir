package hdsfhir

import (
	"encoding/json"
	"io/ioutil"

	fhir "github.com/intervention-engine/fhir/models"
	"github.com/pebbe/util"
	. "gopkg.in/check.v1"
)

type ProcedureSuite struct {
	Patient    *Patient
	Procedures map[string]*Procedure
	Encounter  *Encounter
}

var _ = Suite(&ProcedureSuite{})

func (s *ProcedureSuite) SetUpSuite(c *C) {
	data, err := ioutil.ReadFile("./fixtures/procedures.json")
	util.CheckErr(err)

	s.Procedures = make(map[string]*Procedure)
	err = json.Unmarshal(data, &s.Procedures)
	util.CheckErr(err)

	s.Patient = &Patient{}
	s.Encounter = &Encounter{Entry: Entry{StartTime: 1320148800, EndTime: 1320152400}}
	s.Patient.Encounters = []*Encounter{s.Encounter}
	for _, procedure := range s.Procedures {
		procedure.Patient = s.Patient
	}
}

func (s *ProcedureSuite) TestProcedureResults(c *C) {
	procedureModels := s.Procedures["procedureResults"].FHIRModels()

	c.Assert(procedureModels, HasLen, 5)

	c.Assert(procedureModels[0], FitsTypeOf, &fhir.Procedure{})
	procedure := procedureModels[0].(*fhir.Procedure)
	c.Assert(procedure.Status, Equals, "completed")
	c.Assert(procedure.Subject, DeepEquals, s.Patient.FHIRReference())
	c.Assert(procedure.Code.Text, Equals, "Procedure, Result: Clinical Staging Procedure")
	c.Assert(procedure.Code.Coding, HasLen, 1)
	c.Assert(procedure.Code.MatchesCode("http://snomed.info/sct", "116783008"), Equals, true)
	c.Assert(procedure.NotPerformed, IsNil)
	c.Assert(procedure.ReasonNotPerformed, IsNil)
	c.Assert(procedure.BodySite, IsNil)
	c.Assert(procedure.PerformedPeriod.Start, DeepEquals, NewUnixTime(1320149800).FHIRDateTime())
	c.Assert(procedure.PerformedPeriod.End, DeepEquals, NewUnixTime(1320159800).FHIRDateTime())
	c.Assert(procedure.Encounter, DeepEquals, s.Encounter.FHIRReference())
	c.Assert(procedure.Report, HasLen, 1)

	c.Assert(procedureModels[1], FitsTypeOf, &fhir.DiagnosticReport{})
	report := procedureModels[1].(*fhir.DiagnosticReport)
	c.Assert(report.Status, Equals, "final")
	c.Assert(report.Code.Text, Equals, "Procedure findings narrative")
	c.Assert(report.Code.Coding, HasLen, 1)
	c.Assert(report.Code.MatchesCode("http://loinc.org", "59776-5"), Equals, true)
	c.Assert(report.Subject, DeepEquals, s.Patient.FHIRReference())
	c.Assert(report.Encounter, DeepEquals, s.Encounter.FHIRReference())
	c.Assert(report.EffectivePeriod.Start, DeepEquals, NewUnixTime(1320149800).FHIRDateTime())
	c.Assert(report.EffectivePeriod.End, DeepEquals, NewUnixTime(1320159800).FHIRDateTime())
	c.Assert(report.Issued, DeepEquals, NewUnixTime(1320159800).FHIRDateTime())
	c.Assert(report.Result, HasLen, 3)
	c.Assert(procedure.Report[0].Reference, Equals, "urn:uuid:"+report.Id)

	for i := 2; i < 5; i++ {
		c.Assert(procedureModels[i], FitsTypeOf, &fhir.Observation{})
		observation := procedureModels[i].(*fhir.Observation)
		c.Assert(observation.Code.Text, Equals, "Procedure, Result: Clinical Staging Procedure")
		c.Assert(observation.Code.Coding, HasLen, 1)
		c.Assert(observation.Code.MatchesCode("http://snomed.info/sct", "116783008"), Equals, true)
		c.Assert(observation.Subject, DeepEquals, s.Patient.FHIRReference())
		c.Assert(observation.Encounter, DeepEquals, s.Encounter.FHIRReference())
		c.Assert(observation.EffectivePeriod.Start, DeepEquals, NewUnixTime(1320149800).FHIRDateTime())
		c.Assert(observation.EffectivePeriod.End, DeepEquals, NewUnixTime(1320159800).FHIRDateTime())
		c.Assert(observation.Status, Equals, "final")
		c.Assert(observation.ValueQuantity, IsNil)
		c.Assert(observation.ValueString, Equals, "")
		switch i {
		case 2:
			c.Assert(observation.ValueCodeableConcept.Text, Equals, "Colon Distant Metastasis Status M0")
			c.Assert(observation.ValueCodeableConcept.MatchesCode("http://snomed.info/sct", "433581000124101"), Equals, true)
		case 3:
			c.Assert(observation.ValueCodeableConcept.Text, Equals, "Colon Cancer Regional Lymph Node Status N2b")
			c.Assert(observation.ValueCodeableConcept.MatchesCode("http://snomed.info/sct", "433571000124104"), Equals, true)
		case 4:
			c.Assert(observation.ValueCodeableConcept.Text, Equals, "Colon Cancer Primary Tumor Size T4a")
			c.Assert(observation.ValueCodeableConcept.MatchesCode("http://snomed.info/sct", "433491000124102"), Equals, true)
		}
		c.Assert(report.Result[i-2].Reference, Equals, "urn:uuid:"+observation.Id)
	}
}

func (s *ProcedureSuite) TestProcedurePerformed(c *C) {
	procedureModels := s.Procedures["procedurePerformed"].FHIRModels()
	c.Assert(procedureModels, HasLen, 1)
	c.Assert(procedureModels[0], FitsTypeOf, &fhir.Procedure{})

	procedure := procedureModels[0].(*fhir.Procedure)
	c.Assert(procedure.Status, Equals, "completed")
	c.Assert(procedure.Subject, DeepEquals, s.Patient.FHIRReference())
	c.Assert(procedure.Code.Text, Equals, "Procedure, Performed: Hospital measures-CABG")
	c.Assert(procedure.Code.Coding, HasLen, 3)
	c.Assert(procedure.Code.MatchesCode("http://hl7.org/fhir/sid/icd-9", "36.10"), Equals, true)
	c.Assert(procedure.Code.MatchesCode("http://hl7.org/fhir/sid/icd-10", "0210093"), Equals, true)
	c.Assert(procedure.Code.MatchesCode("http://snomed.info/sct", "10190003"), Equals, true)
	c.Assert(procedure.NotPerformed, IsNil)
	c.Assert(procedure.ReasonNotPerformed, IsNil)
	c.Assert(procedure.BodySite, HasLen, 1)
	c.Assert(procedure.BodySite[0].Text, Equals, "")
	c.Assert(procedure.BodySite[0].Coding, HasLen, 1)
	c.Assert(procedure.BodySite[0].MatchesCode("http://snomed.info/sct", "50018008"), Equals, true)
	c.Assert(procedure.PerformedPeriod.Start, DeepEquals, NewUnixTime(1362239100).FHIRDateTime())
	c.Assert(procedure.PerformedPeriod.End, DeepEquals, NewUnixTime(1362242700).FHIRDateTime())
	c.Assert(procedure.Encounter, IsNil)
	c.Assert(procedure.Report, IsNil)
}

func (s *ProcedureSuite) TestProcedureOrdered(c *C) {
	procedureModels := s.Procedures["procedureOrdered"].FHIRModels()
	c.Assert(procedureModels, HasLen, 1)
	c.Assert(procedureModels[0], FitsTypeOf, &fhir.ProcedureRequest{})

	procedureRequest := procedureModels[0].(*fhir.ProcedureRequest)
	c.Assert(procedureRequest.Subject, DeepEquals, s.Patient.FHIRReference())
	c.Assert(procedureRequest.Status, Equals, "accepted")
	c.Assert(procedureRequest.Code.Text, Equals, "Procedure, Ordered: Hospital measures-CABG")
	c.Assert(procedureRequest.Code.Coding, HasLen, 3)
	c.Assert(procedureRequest.Code.MatchesCode("http://hl7.org/fhir/sid/icd-9", "36.10"), Equals, true)
	c.Assert(procedureRequest.Code.MatchesCode("http://hl7.org/fhir/sid/icd-10", "0210093"), Equals, true)
	c.Assert(procedureRequest.Code.MatchesCode("http://snomed.info/sct", "10190003"), Equals, true)
	c.Assert(procedureRequest.OrderedOn, DeepEquals, NewUnixTime(1362239100).FHIRDateTime())
	c.Assert(procedureRequest.BodySite, HasLen, 1)
	c.Assert(procedureRequest.BodySite[0].Text, Equals, "")
	c.Assert(procedureRequest.BodySite[0].Coding, HasLen, 1)
	c.Assert(procedureRequest.BodySite[0].MatchesCode("http://snomed.info/sct", "50018008"), Equals, true)
	c.Assert(procedureRequest.Encounter, IsNil)
}

func (s *ProcedureSuite) TestProcedureNotPerformed(c *C) {
	procedureModels := s.Procedures["procedureNotPerformed"].FHIRModels()
	c.Assert(procedureModels, HasLen, 1)
	c.Assert(procedureModels[0], FitsTypeOf, &fhir.Procedure{})

	procedure := procedureModels[0].(*fhir.Procedure)
	c.Assert(procedure.Status, Equals, "completed")
	c.Assert(procedure.Subject, DeepEquals, s.Patient.FHIRReference())
	c.Assert(procedure.Code.Text, Equals, "Procedure, Performed: Hospital measures-CABG")
	c.Assert(procedure.Code.Coding, HasLen, 3)
	c.Assert(procedure.Code.MatchesCode("http://hl7.org/fhir/sid/icd-9", "36.10"), Equals, true)
	c.Assert(procedure.Code.MatchesCode("http://hl7.org/fhir/sid/icd-10", "0210093"), Equals, true)
	c.Assert(procedure.Code.MatchesCode("http://snomed.info/sct", "10190003"), Equals, true)
	c.Assert(*procedure.NotPerformed, Equals, true)
	c.Assert(procedure.ReasonNotPerformed, HasLen, 1)
	c.Assert(procedure.ReasonNotPerformed[0].Text, Equals, "")
	c.Assert(procedure.ReasonNotPerformed[0].Coding, HasLen, 1)
	c.Assert(procedure.ReasonNotPerformed[0].MatchesCode("http://snomed.info/sct", "397807004"), Equals, true)
	c.Assert(procedure.BodySite, IsNil)
	c.Assert(procedure.PerformedPeriod.Start, DeepEquals, NewUnixTime(1362239100).FHIRDateTime())
	c.Assert(procedure.PerformedPeriod.End, DeepEquals, NewUnixTime(1362239100).FHIRDateTime())
	c.Assert(procedure.Encounter, IsNil)
	c.Assert(procedure.Report, IsNil)
}

func (s *ProcedureSuite) TestProcedureRejected(c *C) {
	procedureModels := s.Procedures["procedureRejected"].FHIRModels()
	c.Assert(procedureModels, HasLen, 1)
	c.Assert(procedureModels[0], FitsTypeOf, &fhir.ProcedureRequest{})

	procedureRequest := procedureModels[0].(*fhir.ProcedureRequest)
	c.Assert(procedureRequest.Subject, DeepEquals, s.Patient.FHIRReference())
	c.Assert(procedureRequest.Status, Equals, "rejected")
	c.Assert(procedureRequest.Code.Text, Equals, "Procedure, Ordered: Hospital measures-CABG")
	c.Assert(procedureRequest.Code.Coding, HasLen, 3)
	c.Assert(procedureRequest.Code.MatchesCode("http://hl7.org/fhir/sid/icd-9", "36.10"), Equals, true)
	c.Assert(procedureRequest.Code.MatchesCode("http://hl7.org/fhir/sid/icd-10", "0210093"), Equals, true)
	c.Assert(procedureRequest.Code.MatchesCode("http://snomed.info/sct", "10190003"), Equals, true)
	c.Assert(procedureRequest.OrderedOn, DeepEquals, NewUnixTime(1362239100).FHIRDateTime())
	c.Assert(procedureRequest.BodySite, HasLen, 1)
	c.Assert(procedureRequest.BodySite[0].Text, Equals, "")
	c.Assert(procedureRequest.BodySite[0].Coding, HasLen, 1)
	c.Assert(procedureRequest.BodySite[0].MatchesCode("http://snomed.info/sct", "50018008"), Equals, true)
	c.Assert(procedureRequest.Encounter, IsNil)
}
