package hdsfhir

import (
	"time"

	fhir "github.com/intervention-engine/fhir/models"
	. "gopkg.in/check.v1"
)

type UnixTimeSuite struct {
}

var _ = Suite(&UnixTimeSuite{})

func (s *UnixTimeSuite) TestNewUnixTime(c *C) {
	t := NewUnixTime(1445459340)
	c.Assert(t.Time().UTC(), DeepEquals, time.Date(2015, time.October, 21, 20, 29, 0, 0, time.UTC))
}

func (s *UnixTimeSuite) TestFHIRDateTime(c *C) {
	t := NewUnixTime(1445459340)
	c.Assert(t.FHIRDateTime().Precision, Equals, fhir.Precision(fhir.Timestamp))
	c.Assert(t.FHIRDateTime().Time.UTC(), DeepEquals, time.Date(2015, time.October, 21, 20, 29, 0, 0, time.UTC))
}
