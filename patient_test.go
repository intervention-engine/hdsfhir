package hdsfhir

import (
	"encoding/json"
	"io/ioutil"
	"reflect"
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
	s.doTestFHIRTransactionBundle(c, false)
}

func (s *PatientSuite) TestFHIRTransactionBundleConditionalUpdate(c *C) {
	s.doTestFHIRTransactionBundle(c, true)
}

func (s *PatientSuite) doTestFHIRTransactionBundle(c *C, conditionalUpdate bool) {
	bundle := s.Patient.FHIRTransactionBundle(conditionalUpdate)
	c.Assert(bundle.Entry, HasLen, 21)
	c.Assert(bundle.Entry[0].Resource, FitsTypeOf, &fhir.Patient{})
	patientID := bundle.Entry[0].Resource.(*fhir.Patient).Id
	patientRef := "urn:uuid:" + patientID
	for i := range bundle.Entry {
		if conditionalUpdate && i == 0 {
			c.Assert(bundle.Entry[i].Request.Method, Equals, "PUT")
		} else {
			c.Assert(bundle.Entry[i].Request.Method, Equals, "POST")
		}
		switch t := bundle.Entry[i].Resource.(type) {
		case *fhir.Patient:
			if conditionalUpdate {
				c.Assert(bundle.Entry[i].Request.Url, Equals, "Patient?identifier=bc8f60f4cbde3d6c28974971b6880793")
			} else {
				c.Assert(bundle.Entry[i].Request.Url, Equals, "Patient")
			}
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
