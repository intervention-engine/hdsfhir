package hdsfhir

import (
	"encoding/json"
	"io/ioutil"
	"net/url"
	"reflect"
	"strings"
	"testing"

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
	patientRef := "urn%3Auuid%3A" + patientID
	_ = patientRef
	for i := range bundle.Entry {
		c.Assert(bundle.Entry[i].Request.Method, Equals, "PUT")
	}
	assertURL(c, bundle, 0, "Patient?identifier=bc8f60f4cbde3d6c28974971b6880793")
	assertURL(c, bundle, 1, "Encounter?date=sa2011-11-01T07%3A59%3A59\u0026date=lt2011-11-01T08%3A00%3A01\u0026patient="+patientRef+"\u0026type=http%3A%2F%2Fwww.ama-assn.org%2Fgo%2Fcpt%7C99201")
	assertURL(c, bundle, 2, "Encounter?date=sa2012-10-01T07%3A59%3A59\u0026date=lt2012-10-01T08%3A00%3A01\u0026patient="+patientRef+"\u0026type=http%3A%2F%2Fwww.ama-assn.org%2Fgo%2Fcpt%7C99201")
	assertURL(c, bundle, 3, "Encounter?date=sa2012-11-01T07%3A59%3A59\u0026date=lt2012-11-01T08%3A00%3A01\u0026patient="+patientRef+"\u0026type=http%3A%2F%2Fwww.ama-assn.org%2Fgo%2Fcpt%7C99201")
	assertURL(c, bundle, 4, "Encounter?date=sa2012-11-01T08%3A44%3A59\u0026date=lt2012-11-01T08%3A45%3A01\u0026patient="+patientRef+"\u0026type=http%3A%2F%2Fsnomed.info%2Fsct%7C171047005")
	assertURL(c, bundle, 5, "Condition?code=http%3A%2F%2Fhl7.org%2Ffhir%2Fsid%2Ficd-10%7CI50.1%2Chttp%3A%2F%2Fhl7.org%2Ffhir%2Fsid%2Ficd-9%7C428.0%2Chttp%3A%2F%2Fsnomed.info%2Fsct%7C10091002\u0026onset=2012-03-01T07%3A00%3A00\u0026patient="+patientRef)
	assertURL(c, bundle, 6, "Condition?code=http%3A%2F%2Fsnomed.info%2Fsct%7C981000124106\u0026onset=2012-03-01T07%3A05%3A00\u0026patient="+patientRef)
	assertURL(c, bundle, 7, "Condition?code=http%3A%2F%2Fhl7.org%2Ffhir%2Fsid%2Ficd-10%7CI20.0%2Chttp%3A%2F%2Fhl7.org%2Ffhir%2Fsid%2Ficd-9%7C411.0%2Chttp%3A%2F%2Fsnomed.info%2Fsct%7C123641001\u0026onset=2012-03-01T07%3A05%3A00\u0026patient="+patientRef)
	assertURL(c, bundle, 8, "Condition?code=http%3A%2F%2Fhl7.org%2Ffhir%2Fsid%2Ficd-10%7CI20.0%2Chttp%3A%2F%2Fhl7.org%2Ffhir%2Fsid%2Ficd-9%7C411.0%2Chttp%3A%2F%2Fsnomed.info%2Fsct%7C123641001\u0026onset=2012-03-01T07%3A05%3A00\u0026patient="+patientRef)
	assertURL(c, bundle, 9, "Condition?code=http%3A%2F%2Fhl7.org%2Ffhir%2Fsid%2Ficd-9%7C401.1%2Chttp%3A%2F%2Fsnomed.info%2Fsct%7C10725009\u0026onset=2012-03-01T07%3A30%3A00\u0026patient="+patientRef)
	assertURL(c, bundle, 10, "Observation?code=http%3A%2F%2Floinc.org%7C17856-6\u0026date=sa2011-11-01T08%3A16%3A39\u0026date=lt2011-11-01T08%3A16%3A41\u0026patient="+patientRef)
	assertURL(c, bundle, 11, "Procedure?code=http%3A%2F%2Fsnomed.info%2Fsct%7C116783008\u0026date=sa2011-11-01T08%3A16%3A39\u0026date=lt2011-11-01T08%3A16%3A41\u0026patient="+patientRef)
	assertURL(c, bundle, 12, "DiagnosticReport?code=http%3A%2F%2Floinc.org%7C59776-5\u0026date=sa2011-11-01T08%3A16%3A39\u0026date=lt2011-11-01T08%3A16%3A41\u0026patient="+patientRef)
	assertURL(c, bundle, 13, "Observation?code=http%3A%2F%2Fsnomed.info%2Fsct%7C116783008\u0026date=sa2011-11-01T08%3A16%3A39\u0026date=lt2011-11-01T08%3A16%3A41\u0026patient="+patientRef)
	assertURL(c, bundle, 14, "Observation?code=http%3A%2F%2Fsnomed.info%2Fsct%7C116783008\u0026date=sa2011-11-01T08%3A16%3A39\u0026date=lt2011-11-01T08%3A16%3A41\u0026patient="+patientRef)
	assertURL(c, bundle, 15, "Observation?code=http%3A%2F%2Fsnomed.info%2Fsct%7C116783008\u0026date=sa2011-11-01T08%3A16%3A39\u0026date=lt2011-11-01T08%3A16%3A41\u0026patient="+patientRef)
	assertURL(c, bundle, 16, "Procedure?code=http%3A%2F%2Fhl7.org%2Ffhir%2Fsid%2Ficd-10%7C0210093%2Chttp%3A%2F%2Fhl7.org%2Ffhir%2Fsid%2Ficd-9%7C36.10%2Chttp%3A%2F%2Fsnomed.info%2Fsct%7C10190003\u0026date=sa2013-03-02T10%3A44%3A59\u0026date=lt2013-03-02T10%3A45%3A01\u0026patient="+patientRef)
	assertURL(c, bundle, 17, "MedicationStatement?code=http%3A%2F%2Fwww.nlm.nih.gov%2Fresearch%2Fumls%2Frxnorm%2F%7C1000048\u0026effectivedate=sa2012-10-01T07%3A59%3A59\u0026effectivedate=lt2012-10-01T08%3A00%3A01\u0026patient="+patientRef)
	assertURL(c, bundle, 18, "Immunization?date=2011-08-15T08%3A00%3A00\u0026patient="+patientRef+"\u0026vaccine-code=http%3A%2F%2Fwww2a.cdc.gov%2Fvaccines%2Fiis%2Fiisstandards%2Fvaccines.asp%3Frpt%3Dcvx%7C33")
	assertURL(c, bundle, 19, "Immunization?date=2010-01-10T19%3A08%3A28\u0026patient="+patientRef+"\u0026vaccine-code=http%3A%2F%2Fwww2a.cdc.gov%2Fvaccines%2Fiis%2Fiisstandards%2Fvaccines.asp%3Frpt%3Dcvx%7C03")
	assertURL(c, bundle, 20, "AllergyIntolerance?substance=http%3A%2F%2Fwww2a.cdc.gov%2Fvaccines%2Fiis%2Fiisstandards%2Fvaccines.asp%3Frpt%3Dcvx%7C111\u0026onset=2012-01-01T00%3A42%3A00\u0026patient="+patientRef)
}

func assertURL(c *C, b *fhir.Bundle, i int, u string) {
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
