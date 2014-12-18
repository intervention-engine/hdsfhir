package hdsfhir

import (
	"time"

	"github.com/intervention-engine/fhir/models"
)

type FHIRCoding struct {
	System string `json:"system"`
	Code   string `json:"code"`
}

type FHIRCodableConcept struct {
	Codings []FHIRCoding `json:"coding"`
}

type FHIRName struct {
	FirstName []string `json:"given"`
	LastName  []string `json:"family"`
}
type Entry struct {
	Patient     *Patient            `json:"-"`
	StartTime   int64               `json:"start_time"`
	EndTime     int64               `json:"end_time"`
	Time        int64               `json:"time"`
	Oid         string              `json:"oid"`
	Codes       map[string][]string `json:"codes"`
	NegationInd bool                `json:"negationInd"`
	StatusCode  map[string][]string `json:"status_code"`
	Description string              `json:"description"`
	ServerURL   string              `json:"-"`
}

func (e *Entry) SetServerURL(url string) {
	e.ServerURL = url
}

func (e *Entry) ConvertCodingToFHIR() models.CodeableConcept {
	return ConvertCodeMapToFHIR(e.Codes)
}

func (e *Entry) StartTimestamp() time.Time {
	return time.Unix(e.StartTime, 0)
}

func (e *Entry) EndTimestamp() time.Time {
	return time.Unix(e.EndTime, 0)
}

func (e *Entry) AsFHIRPeriod() models.Period {
	period := models.Period{}
	period.Start = models.FHIRDateTime{Time: e.StartTimestamp(), Precision: models.Timestamp}
	period.End = models.FHIRDateTime{Time: e.EndTimestamp(), Precision: models.Timestamp}
	return period
}

func (self *Entry) Overlaps(other Entry) bool {
	return (self.StartTime <= other.EndTime && self.EndTime >= other.StartTime)
}

func ConvertCodeMapToFHIR(codeMap map[string][]string) models.CodeableConcept {
	concept := models.CodeableConcept{}
	codings := make([]models.Coding, 0)
	for codeSystem, codes := range codeMap {
		codeSystemURL := CodeSystemMap[codeSystem]
		for _, code := range codes {
			coding := models.Coding{System: codeSystemURL, Code: code}
			codings = append(codings, coding)
		}
	}
	concept.Coding = codings
	return concept
}

func UnixToFHIRDate(unixTime int64) string {
	return time.Unix(unixTime, 0).Format("2006-01-02")
}

var CodeSystemMap = map[string]string{
	"CPT":        "http://www.ama-assn.org/go/cpt",
	"LOINC":      "http://loinc.org",
	"SNOMED-CT":  "http://snomed.info/sct",
	"RxNorm":     "http://www.nlm.nih.gov/research/umls/rxnorm/",
	"ICD-9-CM":   "http://hl7.org/fhir/sid/icd-9",
	"ICD-10-CM":  "http://hl7.org/fhir/sid/icd-10",
	"ICD-9-PCS":  "http://hl7.org/fhir/sid/icd-9",
	"ICD-10-PCS": "http://hl7.org/fhir/sid/icd-10",
}
