package hdsfhir

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"

	fhir "github.com/intervention-engine/fhir/models"
	"github.com/pebbe/util"
	"github.com/satori/go.uuid"
	. "gopkg.in/check.v1"
)

type PatientSuite struct {
	Patient *Patient
}

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

func (s *PatientSuite) SetUpTest(c *C) {
	data, err := ioutil.ReadFile("./fixtures/john_peters.json")
	util.CheckErr(err)

	s.Patient = &Patient{}
	err = json.Unmarshal(data, s.Patient)
	util.CheckErr(err)
}

var _ = Suite(&PatientSuite{})

func (s *PatientSuite) TestPatientFHIRModel(c *C) {
	model := s.Patient.FHIRModel()
	c.Assert(model, FitsTypeOf, &fhir.Patient{})
	c.Assert(model.Name, HasLen, 1)
	c.Assert(model.Name[0].Given, HasLen, 1)
	c.Assert(model.Name[0].Given[0], Equals, "John")
	c.Assert(model.Name[0].Family, HasLen, 1)
	c.Assert(model.Name[0].Family[0], Equals, "Peters")
	c.Assert(model.Gender, Equals, "male")
	c.Assert(model.BirthDate, DeepEquals, NewUnixTime(665420400).FHIRDate())
	c.Assert(model.Identifier, HasLen, 1)
	c.Assert(model.Identifier[0].Type.MatchesCode("http://hl7.org/fhir/v2/0203", "MR"), Equals, true)
	c.Assert(model.Identifier[0].Value, Equals, "bc8f60f4cbde3d6c28974971b6880793")
}

func (s *PatientSuite) TestPatientFHIRModelWithNoMRN(c *C) {
	s.Patient.MedicalRecordNumber = ""
	model := s.Patient.FHIRModel()
	c.Assert(model, FitsTypeOf, &fhir.Patient{})
	c.Assert(model.Name, HasLen, 1)
	c.Assert(model.Gender, Equals, "male")
	c.Assert(model.BirthDate, DeepEquals, NewUnixTime(665420400).FHIRDate())
	c.Assert(model.Identifier, HasLen, 0)
}

func (s *PatientSuite) TestFHIRModels(c *C) {
	models := s.Patient.FHIRModels()
	patient := s.Patient.FHIRModel()
	c.Assert(models, HasLen, 21)

	typeMap := make(map[string]int)
	for i := range models {
		// Build a map to count the returned models by resource type
		class := reflect.TypeOf(models[i]).Elem().Name()
		count, _ := typeMap[class]
		typeMap[class] = count + 1

		// Test all references to the patient to ensure they were populated
		patientVal := reflect.ValueOf(models[i]).Elem().FieldByName("Subject")
		if !patientVal.IsValid() {
			patientVal = reflect.ValueOf(models[i]).Elem().FieldByName("Patient")
		}
		if patientVal.IsValid() {
			c.Assert(patientVal.Interface().(*fhir.Reference).Reference, Equals, "urn:uuid:"+patient.Id)
		}
	}

	// Test the resource type counts
	c.Assert(typeMap, HasLen, 9)
	c.Assert(typeMap["Condition"], Equals, 5)
	c.Assert(typeMap["DiagnosticReport"], Equals, 1)
	c.Assert(typeMap["Encounter"], Equals, 4)
	c.Assert(typeMap["Immunization"], Equals, 2)
	c.Assert(typeMap["MedicationStatement"], Equals, 1)
	c.Assert(typeMap["Observation"], Equals, 4)
	c.Assert(typeMap["Patient"], Equals, 1)
	c.Assert(typeMap["Procedure"], Equals, 2)
	c.Assert(typeMap["AllergyIntolerance"], Equals, 1)
}

func (s *PatientSuite) TestFHIRModelReferences(c *C) {
	models := s.Patient.FHIRModels()
	refs := getAllReferences(models)
	for i := range refs {
		c.Assert(isReferenceValid(refs[i], models), Equals, true)
	}
}

func (s *PatientSuite) TestFHIRTransactionBundle(c *C) {
	bundle := s.Patient.FHIRTransactionBundle(false)
	c.Assert(bundle.Entry, HasLen, 21)
	c.Assert(bundle.Entry[0].Resource, FitsTypeOf, &fhir.Patient{})
	patientID := bundle.Entry[0].Resource.(*fhir.Patient).Id
	patientRef := "urn:uuid:" + patientID
	for i := range bundle.Entry {
		c.Assert(bundle.Entry[i].Request.Method, Equals, "POST")
		switch t := bundle.Entry[i].Resource.(type) {
		case *fhir.Patient:
			c.Assert(bundle.Entry[i].Request.Url, Equals, "Patient")
		case *fhir.Encounter:
			c.Assert(t.Patient.Reference, Equals, patientRef)
			c.Assert(bundle.Entry[i].Request.Url, Equals, "Encounter")
		case *fhir.Condition:
			c.Assert(t.Patient.Reference, Equals, patientRef)
			c.Assert(bundle.Entry[i].Request.Url, Equals, "Condition")
		case *fhir.Observation:
			c.Assert(t.Subject.Reference, Equals, patientRef)
			c.Assert(bundle.Entry[i].Request.Url, Equals, "Observation")
		case *fhir.Procedure:
			c.Assert(t.Subject.Reference, Equals, patientRef)
			c.Assert(bundle.Entry[i].Request.Url, Equals, "Procedure")
		case *fhir.DiagnosticReport:
			c.Assert(t.Subject.Reference, Equals, patientRef)
			c.Assert(bundle.Entry[i].Request.Url, Equals, "DiagnosticReport")
		case *fhir.MedicationStatement:
			c.Assert(t.Patient.Reference, Equals, patientRef)
			c.Assert(bundle.Entry[i].Request.Url, Equals, "MedicationStatement")
		case *fhir.Immunization:
			c.Assert(t.Patient.Reference, Equals, patientRef)
			c.Assert(bundle.Entry[i].Request.Url, Equals, "Immunization")
		case *fhir.AllergyIntolerance:
			c.Assert(t.Patient.Reference, Equals, patientRef)
			c.Assert(bundle.Entry[i].Request.Url, Equals, "AllergyIntolerance")
		default:
			c.Fail()
		}
	}
}

func (s *PatientSuite) TestFHIRTransactionBundleConditionalUpdate(c *C) {
	bundle := s.Patient.FHIRTransactionBundle(true)
	c.Assert(bundle.Entry, HasLen, 21)
	c.Assert(bundle.Entry[0].Resource, FitsTypeOf, &fhir.Patient{})
	patientID := bundle.Entry[0].Resource.(*fhir.Patient).Id
	patientRef := url.QueryEscape("urn:uuid:" + patientID)
	for i := range bundle.Entry {
		c.Assert(bundle.Entry[i].Request.Method, Equals, "PUT")
	}
	assertURL(c, bundle, 0, "Patient?identifier=bc8f60f4cbde3d6c28974971b6880793")
	assertURL(c, bundle, 1, "Encounter?date=sa%s&date=lt%s&patient=%s&type=http://www.ama-assn.org/go/cpt|99201", ld("2011-11-01T11:59:59"), ld("2011-11-01T12:00:01"), patientRef)
	assertURL(c, bundle, 2, "Encounter?date=sa%s&date=lt%s&patient=%s&type=http://www.ama-assn.org/go/cpt|99201", ld("2012-10-01T11:59:59"), ld("2012-10-01T12:00:01"), patientRef)
	assertURL(c, bundle, 3, "Encounter?date=sa%s&date=lt%s&patient=%s&type=http://www.ama-assn.org/go/cpt|99201", ld("2012-11-01T11:59:59"), ld("2012-11-01T12:00:01"), patientRef)
	assertURL(c, bundle, 4, "Encounter?date=sa%s&date=lt%s&patient=%s&type=http://snomed.info/sct|171047005", ld("2012-11-01T12:44:59"), ld("2012-11-01T12:45:01"), patientRef)
	assertURL(c, bundle, 5, "Condition?code=http://hl7.org/fhir/sid/icd-10|I50.1,http://hl7.org/fhir/sid/icd-9|428.0,http://snomed.info/sct|10091002&onset=%s&patient=%s", ld("2012-03-01T12:00:00"), patientRef)
	assertURL(c, bundle, 6, "Condition?code=http://snomed.info/sct|981000124106&onset=%s&patient=%s", ld("2012-03-01T12:05:00"), patientRef)
	assertURL(c, bundle, 7, "Condition?code=http://hl7.org/fhir/sid/icd-10|I20.0,http://hl7.org/fhir/sid/icd-9|411.0,http://snomed.info/sct|123641001&onset=%s&patient=%s", ld("2012-03-01T12:05:00"), patientRef)
	assertURL(c, bundle, 8, "Condition?code=http://hl7.org/fhir/sid/icd-10|I20.0,http://hl7.org/fhir/sid/icd-9|411.0,http://snomed.info/sct|123641001&onset=%s&patient=%s", ld("2012-03-01T12:05:00"), patientRef)
	assertURL(c, bundle, 9, "Condition?code=http://hl7.org/fhir/sid/icd-9|401.1,http://snomed.info/sct|10725009&onset=%s&patient=%s", ld("2012-03-01T12:30:00"), patientRef)
	assertURL(c, bundle, 10, "Observation?code=http://loinc.org|17856-6&date=sa%s&date=lt%s&patient=%s&value-quantity=8||%%25", ld("2011-11-01T12:16:39"), ld("2011-11-01T12:16:41"), patientRef)
	assertURL(c, bundle, 11, "Procedure?code=http://snomed.info/sct|116783008&date=sa%s&date=lt%s&patient=%s", ld("2011-11-01T12:16:39"), ld("2011-11-01T12:16:41"), patientRef)
	assertURL(c, bundle, 12, "DiagnosticReport?code=http://loinc.org|59776-5&date=sa%s&date=lt%s&patient=%s", ld("2011-11-01T12:16:39"), ld("2011-11-01T12:16:41"), patientRef)
	assertURL(c, bundle, 13, "Observation?code=http://snomed.info/sct|116783008&date=sa%s&date=lt%s&patient=%s&value-concept=http://snomed.info/sct|433581000124101", ld("2011-11-01T12:16:39"), ld("2011-11-01T12:16:41"), patientRef)
	assertURL(c, bundle, 14, "Observation?code=http://snomed.info/sct|116783008&date=sa%s&date=lt%s&patient=%s&value-concept=http://snomed.info/sct|433571000124104", ld("2011-11-01T12:16:39"), ld("2011-11-01T12:16:41"), patientRef)
	assertURL(c, bundle, 15, "Observation?code=http://snomed.info/sct|116783008&date=sa%s&date=lt%s&patient=%s&value-concept=http://snomed.info/sct|433491000124102", ld("2011-11-01T12:16:39"), ld("2011-11-01T12:16:41"), patientRef)
	assertURL(c, bundle, 16, "Procedure?code=http://hl7.org/fhir/sid/icd-10|0210093,http://hl7.org/fhir/sid/icd-9|36.10,http://snomed.info/sct|10190003&date=sa%s&date=lt%s&patient=%s", ld("2013-03-02T15:44:59"), ld("2013-03-02T15:45:01"), patientRef)
	assertURL(c, bundle, 17, "MedicationStatement?code=http://www.nlm.nih.gov/research/umls/rxnorm/|1000048&effectivedate=sa%s&effectivedate=lt%s&patient=%s", ld("2012-10-01T11:59:59"), ld("2012-10-01T12:00:01"), patientRef)
	assertURL(c, bundle, 18, "Immunization?date=%s&patient=%s&vaccine-code=http://www2a.cdc.gov/vaccines/iis/iisstandards/vaccines.asp%%3Frpt%%3Dcvx|33", ld("2011-08-15T12:00:00"), patientRef)
	assertURL(c, bundle, 19, "Immunization?date=%s&patient=%s&vaccine-code=http://www2a.cdc.gov/vaccines/iis/iisstandards/vaccines.asp%%3Frpt%%3Dcvx|03", ld("2010-01-11T00:08:28"), patientRef)
	assertURL(c, bundle, 20, "AllergyIntolerance?substance=http://www2a.cdc.gov/vaccines/iis/iisstandards/vaccines.asp%%3Frpt%%3Dcvx|111&onset=%s&patient=%s", ld("2012-01-01T05:42:00"), patientRef)
}

// ld takes in a string utc date and converts to a string date in the local timezone
func ld(utcDate string) string {
	if utcDate == "" {
		return ""
	}
	d, err := time.ParseInLocation("2006-01-02T15:04:05", utcDate, time.UTC)
	util.CheckErr(err)
	local := d.In(time.Local).Format("2006-01-02T15:04:05-07:00")
	return url.QueryEscape(local)
}

func assertURL(c *C, b *fhir.Bundle, i int, u string, args ...interface{}) {
	u = fmt.Sprintf(u, args...)
	obtSplit := strings.SplitN(b.Entry[i].Request.Url, "?", 2)
	obtVals, err := url.ParseQuery(obtSplit[1])
	util.CheckErr(err)
	expSplit := strings.SplitN(u, "?", 2)
	expVals, err := url.ParseQuery(expSplit[1])
	util.CheckErr(err)
	c.Assert(obtSplit[0], Equals, expSplit[0])
	c.Assert(obtVals, DeepEquals, expVals)
}

func (s *PatientSuite) TestFHIRBundleReferences(c *C) {
	s.doTestFHIRBundleReferences(c, false)
	s.doTestFHIRBundleReferences(c, true)
}

func (s *PatientSuite) doTestFHIRBundleReferences(c *C, conditionalUpdate bool) {
	bundle := s.Patient.FHIRTransactionBundle(conditionalUpdate)
	models := make([]interface{}, len(bundle.Entry))
	for i := range bundle.Entry {
		models[i] = bundle.Entry[i].Resource
	}
	refs := getAllReferences(models)
	for i := range refs {
		c.Assert(isReferenceValid(refs[i], models), Equals, true)
	}
}

func getAllReferences(models []interface{}) []*fhir.Reference {
	var refs []*fhir.Reference
	for i := range models {
		s := reflect.ValueOf(models[i]).Elem()
		for j := 0; j < s.NumField(); j++ {
			f := s.Field(j)
			if f.Type() == reflect.TypeOf(&fhir.Reference{}) && !f.IsNil() {
				refs = append(refs, f.Interface().(*fhir.Reference))
			} else if f.Type() == reflect.TypeOf([]fhir.Reference{}) {
				for k := 0; k < f.Len(); k++ {
					refs = append(refs, f.Index(k).Addr().Interface().(*fhir.Reference))
				}
			}
		}
	}
	return refs
}

func isReferenceValid(ref *fhir.Reference, models []interface{}) bool {
	for i := range models {
		id := reflect.ValueOf(models[i]).Elem().FieldByName("Id").String()
		_, err := uuid.FromString(id)
		util.CheckErr(err)
		if ref.Reference == "urn:uuid:"+id {
			return true
		}
	}
	return false
}
