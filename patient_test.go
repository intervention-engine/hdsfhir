package hdsfhir

import (
	"encoding/json"
	"io/ioutil"
	"reflect"
	"testing"

	fhir "github.com/intervention-engine/fhir/models"
	"github.com/pebbe/util"
	. "gopkg.in/check.v1"
)

type PatientSuite struct {
	Patient *Patient
}

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

func (s *PatientSuite) SetUpSuite(c *C) {
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
}

func (s *PatientSuite) TestFHIRModels(c *C) {
	models := s.Patient.FHIRModels()
	patient := s.Patient.FHIRModel()
	c.Assert(models, HasLen, 20)

	// Test all references to the patient to ensure they were populated
	for i := range models {
		patientVal := reflect.ValueOf(models[i]).Elem().FieldByName("Subject")
		if !patientVal.IsValid() {
			patientVal = reflect.ValueOf(models[i]).Elem().FieldByName("Patient")
		}
		if patientVal.IsValid() {
			c.Assert(patientVal.Interface().(*fhir.Reference).Reference, Equals, "cid:"+patient.Id)
		}
	}
}

func (s *PatientSuite) TestFHIRModelReferences(c *C) {
	models := s.Patient.FHIRModels()
	refs := getAllReferences(models)
	for i := range refs {
		c.Assert(isReferenceValid(refs[i], models), Equals, true)
	}
}

func getAllReferences(models []interface{}) []*fhir.Reference {
	refs := make([]*fhir.Reference, 0)
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
		if ref.Reference == "cid:"+id {
			return true
		}
	}
	return false
}
