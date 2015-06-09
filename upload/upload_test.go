package upload

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"regexp"
	"strings"
	"testing"

	fhir "github.com/intervention-engine/fhir/models"
	"github.com/intervention-engine/hdsfhir/models"
	"github.com/pebbe/util"
	. "gopkg.in/check.v1"
)

type UploadSuite struct {
	JSONBlob []byte
}

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

func (s *UploadSuite) SetUpSuite(c *C) {
	data, err := ioutil.ReadFile("../fixtures/john_peters.json")
	util.CheckErr(err)
	s.JSONBlob = data
}

var _ = Suite(&UploadSuite{})

func (s *UploadSuite) TestUpload(c *C) {
	patient := &models.Patient{}
	err := json.Unmarshal(s.JSONBlob, patient)
	util.CheckErr(err)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.Assert(r.Header.Get("Content-Type"), Equals, "application/json+fhir")
		w.Header().Add("Location", "http://localhost/Patient/1234")
		fmt.Fprintln(w, "Created")
	}))
	defer ts.Close()
	refMap, err := UploadResources(patient.FHIRModels(), ts.URL)
	util.CheckErr(err)
	c.Assert(refMap[patient.GetTempID()], Equals, "http://localhost/Patient/1234")
}

func (s *UploadSuite) TestPostToFHIRServer(c *C) {
	patient := &models.Patient{}
	err := json.Unmarshal(s.JSONBlob, patient)
	util.CheckErr(err)
	resourceCount, patientCount, encounterCount, conditionCount, immunizationCount, observationCount, procedureCount, diagnosticReportCount, medicationCount, medicationStatementCount := 0, 0, 0, 0, 0, 0, 0, 0, 0, 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		output := "Created"
		decoder := json.NewDecoder(r.Body)
		switch {
		case strings.Contains(r.RequestURI, "Patient"):
			if isValid(decoder, &fhir.Patient{}) {
				w.Header().Add("Location", fmt.Sprintf("http://localhost/Patient/%d", resourceCount))
				patientCount++
			}
		case strings.Contains(r.RequestURI, "Encounter"):
			if isValid(decoder, &fhir.Encounter{}) {
				w.Header().Add("Location", fmt.Sprintf("http://localhost/Encounter/%d", resourceCount))
				encounterCount++
			}
		case strings.Contains(r.RequestURI, "Condition"):
			if isValid(decoder, &fhir.Condition{}) {
				w.Header().Add("Location", fmt.Sprintf("http://localhost/Condition/%d", resourceCount))
				conditionCount++
			}
		case strings.Contains(r.RequestURI, "Immunization"):
			if isValid(decoder, &fhir.Immunization{}) {
				w.Header().Add("Location", fmt.Sprintf("http://localhost/Immunization/%d", resourceCount))
				immunizationCount++
			}
		case strings.Contains(r.RequestURI, "Observation"):
			if isValid(decoder, &fhir.Observation{}) {
				w.Header().Add("Location", fmt.Sprintf("http://localhost/Observation/%d", resourceCount))
				observationCount++
			}
		case strings.Contains(r.RequestURI, "Procedure"):
			if isValid(decoder, &fhir.Procedure{}) {
				w.Header().Add("Location", fmt.Sprintf("http://localhost/Procedure/%d", resourceCount))
				procedureCount++
			}
		case strings.Contains(r.RequestURI, "DiagnosticReport"):
			if isValid(decoder, &fhir.DiagnosticReport{}) {
				w.Header().Add("Location", fmt.Sprintf("http://localhost/DiagnosticReport/%d", resourceCount))
				diagnosticReportCount++
			}
		case strings.Contains(r.RequestURI, "MedicationStatement"):
			if isValid(decoder, &fhir.MedicationStatement{}) {
				w.Header().Add("Location", fmt.Sprintf("http://localhost/MedicationStatement/%d", resourceCount))
				medicationStatementCount++
			}
		case strings.Contains(r.RequestURI, "Medication"):
			if isValid(decoder, &fhir.Medication{}) {
				w.Header().Add("Location", fmt.Sprintf("http://localhost/Medication/%d", resourceCount))
				medicationCount++
			}
		}
		fmt.Fprintln(w, output)
		resourceCount++
	}))
	defer ts.Close()
	refMap, err := UploadResources(patient.FHIRModels(), ts.URL)
	c.Assert(patientCount, Equals, 1)
	c.Assert(encounterCount, Equals, 4)
	c.Assert(conditionCount, Equals, 5)
	c.Assert(immunizationCount, Equals, 1)
	c.Assert(observationCount, Equals, 4)
	c.Assert(procedureCount, Equals, 2)
	c.Assert(diagnosticReportCount, Equals, 1)
	c.Assert(resourceCount, Equals, 20)
	c.Assert(medicationStatementCount, Equals, 1)
	c.Assert(medicationCount, Equals, 1)

	c.Assert(len(refMap), Equals, 20)
	c.Assert(refMap[patient.GetTempID()], Equals, "http://localhost/Patient/0")
}

func isValid(decoder *json.Decoder, model interface{}) bool {
	err := decoder.Decode(model)
	if err != nil {
		return false
	}

	_, isPatient := model.(*fhir.Patient)
	_, isMedication := model.(*fhir.Medication)
	if !isPatient && !isMedication {
		refs := getAllReferences(model)
		for _, ref := range refs {
			match, _ := regexp.MatchString("\\Ahttp://localhost/[^/]+/[0-9a-f]+\\z", ref.Reference)
			if !match {
				fmt.Printf("Invalid reference: %s", ref.Reference)
				return false
			}
		}
		return len(refs) > 0
	}
	return true
}

func getAllReferences(model interface{}) []*fhir.Reference {
	refs := make([]*fhir.Reference, 0)
	s := reflect.ValueOf(model).Elem()
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		if f.Type() == reflect.TypeOf(&fhir.Reference{}) && !f.IsNil() {
			refs = append(refs, f.Interface().(*fhir.Reference))
		} else if f.Type() == reflect.TypeOf([]fhir.Reference{}) {
			for j := 0; j < f.Len(); j++ {
				refs = append(refs, f.Index(j).Addr().Interface().(*fhir.Reference))
			}
		}
	}
	return refs
}
