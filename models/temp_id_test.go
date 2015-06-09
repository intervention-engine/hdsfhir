package models

import (
	. "gopkg.in/check.v1"
)

type TempIDSuite struct {
}

var _ = Suite(&TempIDSuite{})

func (s *TempIDSuite) TestGetTempID(c *C) {
	tempID := TemporallyIdentified{}
	id := tempID.GetTempID()
	c.Assert(id, Not(Equals), "")
	c.Assert(tempID.GetTempID(), Equals, id)

	tempID2 := TemporallyIdentified{}
	id2 := tempID2.GetTempID()
	c.Assert(id2, Not(Equals), "")
	c.Assert(id2, Not(Equals), id)
}

func (s *TempIDSuite) TestFHIRReference(c *C) {
	tempID := TemporallyIdentified{}
	ref := tempID.FHIRReference()
	c.Assert(ref.Reference, Equals, "cid:"+tempID.GetTempID())
	c.Assert(ref, DeepEquals, tempID.FHIRReference())
}
