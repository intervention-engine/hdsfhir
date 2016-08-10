package hdsfhir

import (
	"time"

	fhir "github.com/intervention-engine/fhir/models"
	. "gopkg.in/check.v1"
)

type EntrySuite struct {
}

var _ = Suite(&EntrySuite{})

func (s *EntrySuite) TestGetFHIRPeriod(c *C) {
	entry := Entry{StartTime: NewUnixTime(1320148800), EndTime: NewUnixTime(1320152400)}
	period := entry.GetFHIRPeriod()

	c.Assert(period, FitsTypeOf, &fhir.Period{})
	c.Assert(period.Start.Precision, Equals, fhir.Precision(fhir.Timestamp))
	c.Assert(period.Start.Time.UTC(), Equals, time.Date(2011, time.November, 1, 12, 0, 0, 0, time.UTC))
	c.Assert(period.End.Precision, Equals, fhir.Precision(fhir.Timestamp))
	c.Assert(period.End.Time.UTC(), Equals, time.Date(2011, time.November, 1, 13, 0, 0, 0, time.UTC))
}

func (s *EntrySuite) TestGetFHIRPeriodWithNoStart(c *C) {
	entry := Entry{EndTime: NewUnixTime(1320152400)}
	period := entry.GetFHIRPeriod()

	c.Assert(period, FitsTypeOf, &fhir.Period{})
	c.Assert(period.Start, IsNil)
	c.Assert(period.End.Precision, Equals, fhir.Precision(fhir.Timestamp))
	c.Assert(period.End.Time.UTC(), Equals, time.Date(2011, time.November, 1, 13, 0, 0, 0, time.UTC))
}

func (s *EntrySuite) TestGetFHIRPeriodWithNoEnd(c *C) {
	entry := Entry{StartTime: NewUnixTime(1320148800)}
	period := entry.GetFHIRPeriod()

	c.Assert(period, FitsTypeOf, &fhir.Period{})
	c.Assert(period.Start.Precision, Equals, fhir.Precision(fhir.Timestamp))
	c.Assert(period.Start.Time.UTC(), Equals, time.Date(2011, time.November, 1, 12, 0, 0, 0, time.UTC))
	c.Assert(period.End, IsNil)
}

func (s *EntrySuite) TestGetFHIRPeriodWithNoStartOrEnd(c *C) {
	entry := Entry{}
	period := entry.GetFHIRPeriod()

	c.Assert(period, IsNil)
}
