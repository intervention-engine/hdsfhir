package models

import (
	. "gopkg.in/check.v1"
)

type CoreSuite struct {
}

var _ = Suite(&CoreSuite{})

func (s *CoreSuite) TestExtractPatient(c *C) {
	entry := Entry{Codes: map[string][]string{
		"SNOMED-CT": []string{"1234", "5678"},
		"CPT":       []string{"abcd"},
	}}

	codings := entry.ConvertCodingToFHIR()
	found := false
	for _, coding := range codings.Coding {
		if coding.Code == "abcd" && coding.System == "http://www.ama-assn.org/go/cpt" {
			found = true
		}
	}
	c.Assert(found, Equals, true)
}
