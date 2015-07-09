package hdsfhir

import (
	"encoding/json"
	"io/ioutil"

	fhir "github.com/intervention-engine/fhir/models"
	"github.com/pebbe/util"
	. "gopkg.in/check.v1"
)

type ProcedureSuite struct {
	Patient *Patient
}

var _ = Suite(&ProcedureSuite{})

func (suite *ProcedureSuite) SetUpSuite(c *C) {
	data, err := ioutil.ReadFile("./fixtures/john_peters.json")
	util.CheckErr(err)

	suite.Patient = &Patient{}
	err = json.Unmarshal(data, suite.Patient)
	util.CheckErr(err)
}

func (suite *ProcedureSuite) TestFHIRModels(c *C) {
	procedureModels := suite.Patient.Procedures[0].FHIRModels()

	c.Assert(procedureModels, HasLen, 5)

	c.Assert(procedureModels[0], FitsTypeOf, &fhir.Procedure{})
	procedure := procedureModels[0].(*fhir.Procedure)
	c.Assert(procedure.Patient, DeepEquals, suite.Patient.FHIRReference())
	c.Assert(procedure.Type.Text, Equals, "Procedure, Result: Clinical Staging Procedure")
	c.Assert(procedure.Type.MatchesCode("http://snomed.info/sct", "116783008"), Equals, true)
	c.Assert(procedure.PerformedPeriod.Start, DeepEquals, NewUnixTime(1320149800).FHIRDateTime())
	c.Assert(procedure.PerformedPeriod.End, DeepEquals, NewUnixTime(1320159800).FHIRDateTime())
	c.Assert(procedure.Encounter, DeepEquals, suite.Patient.Encounters[0].FHIRReference())
	c.Assert(procedure.Report, HasLen, 1)
	c.Assert(procedure.Notes, Equals, "Procedure, Result: Clinical Staging Procedure")

	c.Assert(procedureModels[1], FitsTypeOf, &fhir.DiagnosticReport{})
	report := procedureModels[1].(*fhir.DiagnosticReport)
	c.Assert(report.Subject, DeepEquals, suite.Patient.FHIRReference())
	c.Assert(report.Result, HasLen, 3)
	c.Assert(procedure.Report[0].Reference, Equals, "cid:"+report.Id)

	for i := 2; i < 5; i++ {
		c.Assert(procedureModels[i], FitsTypeOf, &fhir.Observation{})
		observation := procedureModels[i].(*fhir.Observation)
		c.Assert(observation.Subject, DeepEquals, suite.Patient.FHIRReference())
		c.Assert(observation.EffectivePeriod.Start, DeepEquals, NewUnixTime(1320149800).FHIRDateTime())
		c.Assert(observation.EffectivePeriod.End, DeepEquals, NewUnixTime(1320159800).FHIRDateTime())
		c.Assert(observation.Code.Text, Equals, "Procedure, Result: Clinical Staging Procedure")
		c.Assert(observation.Code.MatchesCode("http://snomed.info/sct", "116783008"), Equals, true)
		c.Assert(observation.Reliability, Equals, "ok")
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
		c.Assert(report.Result[i-2].Reference, Equals, "cid:"+observation.Id)
	}

	procedureModels = suite.Patient.Procedures[1].FHIRModels()
	c.Assert(procedureModels, HasLen, 1)
	c.Assert(procedureModels[0], FitsTypeOf, &fhir.Procedure{})

	procedure = procedureModels[0].(*fhir.Procedure)
	c.Assert(procedure.Patient, DeepEquals, suite.Patient.FHIRReference())
	c.Assert(procedure.Type.Text, Equals, "Procedure, Performed: Hospital measures-CABG")
	c.Assert(procedure.Type.MatchesCode("http://hl7.org/fhir/sid/icd-9", "36.10"), Equals, true)
	c.Assert(procedure.Type.MatchesCode("http://hl7.org/fhir/sid/icd-10", "0210093"), Equals, true)
	c.Assert(procedure.Type.MatchesCode("http://snomed.info/sct", "10190003"), Equals, true)
	c.Assert(procedure.PerformedPeriod.Start, DeepEquals, NewUnixTime(1362239100).FHIRDateTime())
	c.Assert(procedure.PerformedPeriod.End, DeepEquals, NewUnixTime(1362242700).FHIRDateTime())
	c.Assert(procedure.Encounter, IsNil)
	c.Assert(procedure.Report, HasLen, 0)
	c.Assert(procedure.Notes, Equals, "Procedure, Performed: Hospital measures-CABG")
}
