package models

import (
	"time"

	fhir "github.com/intervention-engine/fhir/models"
)

type UnixTime int64

func NewUnixTime(offset int64) *UnixTime {
	t := UnixTime(offset)
	return &t
}

func (t *UnixTime) Time() time.Time {
	return time.Unix(int64(*t), 0)
}

func (t *UnixTime) FHIRDateTime() *fhir.FHIRDateTime {
	return &fhir.FHIRDateTime{Time: t.Time(), Precision: fhir.Timestamp}
}
