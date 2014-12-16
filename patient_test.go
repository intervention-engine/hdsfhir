package hdsfhir

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/pebbe/util"
	. "gopkg.in/check.v1"
)

type HDSSuite struct {
	JSONBlob []byte
}

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

func (s *HDSSuite) SetUpSuite(c *C) {
	data, err := ioutil.ReadFile("./fixtures/john_peters.json")
	util.CheckErr(err)
	s.JSONBlob = data
}

var _ = Suite(&HDSSuite{})

func (s *HDSSuite) TestExtractPatient(c *C) {
	patient := &Patient{}
	err := json.Unmarshal(s.JSONBlob, patient)
	util.CheckErr(err)
	c.Assert("John", Equals, patient.FirstName)
}

func (s *HDSSuite) TestBirthTime(c *C) {
	patient := &Patient{}
	err := json.Unmarshal(s.JSONBlob, patient)
	util.CheckErr(err)
	c.Assert(time.February, Equals, patient.BirthTime().Month())
}

func (s *HDSSuite) TestToJSON(c *C) {
	patient := &Patient{}
	err := json.Unmarshal(s.JSONBlob, patient)
	util.CheckErr(err)
	data := patient.ToJSON()
	c.Assert(data, NotNil)
}

func (s *HDSSuite) TestUpload(c *C) {
	patient := &Patient{}
	err := json.Unmarshal(s.JSONBlob, patient)
	util.CheckErr(err)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.Assert(r.Header.Get("Content-Type"), Equals, "application/json+fhir")
		w.Header().Add("Location", "http://localhost/patients/1234")
		fmt.Fprintln(w, "Created")
	}))
	defer ts.Close()
	Upload(patient, ts.URL)
	c.Assert(patient.ServerURL, Equals, "http://localhost/patients/1234")
}

func (s *HDSSuite) TestPostToFHIRServer(c *C) {
	patient := &Patient{}
	err := json.Unmarshal(s.JSONBlob, patient)
	util.CheckErr(err)
	resourceCount, patientCount, encounterCount, conditionCount, observationCount, procedureCount, diagnosticReportCount := 0, 0, 0, 0, 0, 0, 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.Assert(r.Header.Get("Content-Type"), Equals, "application/json+fhir")
		switch {
		case strings.Contains(r.RequestURI, "Patient"):
			w.Header().Add("Location", fmt.Sprintf("http://localhost/Patients/%d", resourceCount))
			patientCount++
		case strings.Contains(r.RequestURI, "Encounter"):
			w.Header().Add("Location", fmt.Sprintf("http://localhost/Encounters/%d", resourceCount))
			encounterCount++
		case strings.Contains(r.RequestURI, "Condition"):
			w.Header().Add("Location", fmt.Sprintf("http://localhost/Conditions/%d", resourceCount))
			conditionCount++
		case strings.Contains(r.RequestURI, "Observation"):
			w.Header().Add("Location", fmt.Sprintf("http://localhost/Observations/%d", resourceCount))
			observationCount++
		case strings.Contains(r.RequestURI, "Procedure"):
			w.Header().Add("Location", fmt.Sprintf("http://localhost/Procedures/%d", resourceCount))
			procedureCount++
		case strings.Contains(r.RequestURI, "DiagnosticReport"):
			w.Header().Add("Location", fmt.Sprintf("http://localhost/DiagnosticReports/%d", resourceCount))
			diagnosticReportCount++
		}
		fmt.Fprintln(w, "Created")
		resourceCount++
	}))
	defer ts.Close()
	patient.PostToFHIRServer(ts.URL)
	c.Assert(patientCount, Equals, 1)
	c.Assert(encounterCount, Equals, 4)
	c.Assert(conditionCount, Equals, 5)
	c.Assert(observationCount, Equals, 4)
	c.Assert(procedureCount, Equals, 2)
	c.Assert(diagnosticReportCount, Equals, 1)
	c.Assert(patient.ServerURL, Equals, "http://localhost/Patients/0")
	c.Assert(resourceCount, Equals, 17)
	c.Assert(patient.Encounters[0].ServerURL, Equals, "http://localhost/Encounters/1")
}
