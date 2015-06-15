package hdsfhir

import . "gopkg.in/check.v1"

type CodeMapSuite struct {
}

var _ = Suite(&CodeMapSuite{})

func (s *CodeMapSuite) TestCodeMapToCodeableConcept(c *C) {
	codeMap := CodeMap{
		"SNOMED-CT": []string{"1234", "5678"},
		"CPT":       []string{"abcd"}}

	concept := codeMap.FHIRCodeableConcept("test")
	c.Assert(concept.Text, Equals, "test")
	c.Assert(concept.Coding, HasLen, 3)
	c.Assert(concept.MatchesCode("http://snomed.info/sct", "1234"), Equals, true)
	c.Assert(concept.MatchesCode("http://snomed.info/sct", "5678"), Equals, true)
	c.Assert(concept.MatchesCode("http://www.ama-assn.org/go/cpt", "abcd"), Equals, true)
}
