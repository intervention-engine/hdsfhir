package hdsfhir

import (
	"encoding/json"
	"fmt"
	"github.com/pebbe/util"
	. "gopkg.in/check.v1"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
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
	resourceCount := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.Assert(r.Header.Get("Content-Type"), Equals, "application/json+fhir")
		switch {
		case strings.Contains(r.RequestURI, "Patient"):
			w.Header().Add("Location", fmt.Sprintf("http://localhost/Patients/%d", resourceCount))
		case strings.Contains(r.RequestURI, "Encounter"):
			w.Header().Add("Location", fmt.Sprintf("http://localhost/Encounters/%d", resourceCount))
		case strings.Contains(r.RequestURI, "Condition"):
			w.Header().Add("Location", fmt.Sprintf("http://localhost/Conditions/%d", resourceCount))
		}
		fmt.Fprintln(w, "Created")
		resourceCount++
	}))
	defer ts.Close()
	patient.PostToFHIRServer(ts.URL)
	c.Assert(patient.ServerURL, Equals, "http://localhost/Patients/0")
	c.Assert(resourceCount, Equals, 11)
	c.Assert(patient.Encounters[0].ServerURL, Equals, "http://localhost/Encounters/1")
}
