package hdsfhir

import fhir "github.com/intervention-engine/fhir/models"

type Entry struct {
	TemporallyIdentified
	Patient     *Patient `json:"-"`
	StartTime   UnixTime `json:"start_time"`
	EndTime     UnixTime `json:"end_time"`
	Time        UnixTime `json:"time"`
	Oid         string   `json:"oid"`
	Codes       CodeMap  `json:"codes"`
	MoodCode    string   `json:"mood_code"`
	NegationInd bool     `json:"negationInd"`
	StatusCode  CodeMap  `json:"status_code"`
	Description string   `json:"description"`
}

func (e *Entry) GetFHIRPeriod() *fhir.Period {
	return &fhir.Period{
		Start: e.StartTime.FHIRDateTime(),
		End:   e.EndTime.FHIRDateTime(),
	}
}
