package hdsfhir

import fhir "github.com/intervention-engine/fhir/models"

type Entry struct {
	TemporallyIdentified
	Patient        *Patient    `json:"-"`
	StartTime      *UnixTime   `json:"start_time"`
	EndTime        *UnixTime   `json:"end_time"`
	Time           *UnixTime   `json:"time"`
	Oid            string      `json:"oid"`
	Codes          CodeMap     `json:"codes"`
	MoodCode       string      `json:"mood_code"`
	NegationInd    bool        `json:"negationInd"`
	NegationReason *CodeObject `json:"negationReason"`
	StatusCode     CodeMap     `json:"status_code"`
	Description    string      `json:"description"`
}

func (e *Entry) GetFHIRPeriod() *fhir.Period {
	if e.StartTime == nil && e.EndTime == nil {
		return nil
	}

	period := new(fhir.Period)
	if e.StartTime != nil {
		period.Start = e.StartTime.FHIRDateTime()
	}
	if e.EndTime != nil {
		period.End = e.EndTime.FHIRDateTime()
	}

	return period
}
