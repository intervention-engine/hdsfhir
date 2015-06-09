package models

import (
	"time"

	fhir "github.com/intervention-engine/fhir/models"
	. "gopkg.in/check.v1"
)

type EntrySuite struct {
}

var _ = Suite(&EntrySuite{})

func (s *EntrySuite) TestGetFHIRPeriod(c *C) {
	entry := Entry{StartTime: UnixTime(1320148800), EndTime: UnixTime(1320152400)}
	period := entry.GetFHIRPeriod()

	c.Assert(period, FitsTypeOf, &fhir.Period{})
	c.Assert(period.Start.Precision, Equals, fhir.Precision(fhir.Timestamp))
	c.Assert(period.Start.Time.UTC(), Equals, time.Date(2011, time.November, 1, 12, 0, 0, 0, time.UTC))
	c.Assert(period.End.Precision, Equals, fhir.Precision(fhir.Timestamp))
	c.Assert(period.End.Time.UTC(), Equals, time.Date(2011, time.November, 1, 13, 0, 0, 0, time.UTC))
}
