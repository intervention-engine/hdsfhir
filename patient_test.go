package hdsfhir

import (
	"encoding/json"
	"github.com/pebbe/util"
	. "gopkg.in/check.v1"
	"io/ioutil"
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
	patient := Patient{}
	err := json.Unmarshal(s.JSONBlob, &patient)
	util.CheckErr(err)
	c.Assert("John", Equals, patient.FirstName)
}

func (s *HDSSuite) TestBirthTime(c *C) {
	patient := Patient{}
	err := json.Unmarshal(s.JSONBlob, &patient)
	util.CheckErr(err)
	c.Assert(time.February, Equals, patient.BirthTime().Month())
}

func (s *HDSSuite) TestToJSON(c *C) {
	patient := Patient{}
	err := json.Unmarshal(s.JSONBlob, &patient)
	util.CheckErr(err)
	data := patient.ToJSON()
	c.Assert(data, NotNil)
}
