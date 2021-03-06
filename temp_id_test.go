package hdsfhir

import (
	"github.com/pebbe/util"
	"github.com/satori/go.uuid"
	. "gopkg.in/check.v1"
)

type TempIDSuite struct {
}

var _ = Suite(&TempIDSuite{})

func (s *TempIDSuite) TestGetTempID(c *C) {
	tempID := TemporallyIdentified{}
	id := tempID.GetTempID()
	_, err := uuid.FromString(id)
	util.CheckErr(err)
	c.Assert(tempID.GetTempID(), Equals, id)

	tempID2 := TemporallyIdentified{}
	id2 := tempID2.GetTempID()
	_, err = uuid.FromString(id2)
	util.CheckErr(err)
	c.Assert(id2, Not(Equals), id)
}

func (s *TempIDSuite) TestFHIRReference(c *C) {
	tempID := TemporallyIdentified{}
	ref := tempID.FHIRReference()
	c.Assert(ref.Reference, Equals, "urn:uuid:"+tempID.GetTempID())
	c.Assert(ref, DeepEquals, tempID.FHIRReference())
}
