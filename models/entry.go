package models

import fhir "github.com/intervention-engine/fhir/models"

type Entry struct {
	Patient     *Patient            `json:"-"`
	StartTime   UnixTime            `json:"start_time"`
	EndTime     UnixTime            `json:"end_time"`
	Time        UnixTime            `json:"time"`
	Oid         string              `json:"oid"`
	Codes       CodeMap             `json:"codes"`
	NegationInd bool                `json:"negationInd"`
	StatusCode  map[string][]string `json:"status_code"`
	Description string              `json:"description"`
	ServerURL   string              `json:"-"`
}

func (e *Entry) SetServerURL(url string) {
	e.ServerURL = url
}

func (e *Entry) GetFHIRPeriod() *fhir.Period {
	period := &fhir.Period{}
	period.Start = &fhir.FHIRDateTime{Time: e.StartTime.Time(), Precision: fhir.Timestamp}
	period.End = &fhir.FHIRDateTime{Time: e.EndTime.Time(), Precision: fhir.Timestamp}
	return period
}
