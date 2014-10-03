package hdsfhir

import (
	"encoding/json"
	"github.com/pebbe/util"
	. "gopkg.in/check.v1"
	"io/ioutil"
)

type HDSSuite struct {
	JSONBlob *[]byte
}

func (s *HDSSuite) SetUpSuite(c *C) {
	data, err := ioutil.ReadFile("./fixtures/john_peters.json")
	util.CheckErr(err)
	s.JSONBlob = data
}

var _ = Suite(&HDSSuite{})

func (s *HDSSuite) TestExtractPATIENT(c *C) {
	patient := Patient{}
	err := json.Unmarshal(s.JSONBlob, &patient)
	util.CheckErr(err)
	c.Assert("John", Equals, patient.FirstName)
}
