package hdsfhir

import (
	"bytes"
	"gitlab.mitre.org/intervention-engine/fhir/models"
	"net/http"
	"time"
)

type Uploadable interface {
	ToJSON() []byte
	SetServerURL(string)
}

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
	concept := models.CodeableConcept{}
	codings := make([]models.Coding, 0)
	for codeSystem, codes := range e.Codes {
		codeSystemURL := CodeSystemMap[codeSystem]
		for _, code := range codes {
			coding := models.Coding{System: codeSystemURL, Code: code}
			codings = append(codings, coding)
		}
	}
	concept.Coding = codings
	return concept
}

func (e *Entry) StartTimestamp() time.Time {
	return time.Unix(e.StartTime, 0)
}

func (e *Entry) EndTimestamp() time.Time {
	return time.Unix(e.EndTime, 0)
}

func UnixToFHIRDate(unixTime int64) string {
	return time.Unix(unixTime, 0).Format("2006-01-02")
}

func Upload(thing Uploadable, url string) {
	body := bytes.NewReader(thing.ToJSON())
	response, err := http.Post(url, "application/json+fhir", body)
	if err != nil {
		panic("HTTP request failed")
	}
	thing.SetServerURL(response.Header.Get("Location"))
}

var CodeSystemMap = map[string]string{
	"CPT":       "http://www.ama-assn.org/go/cpt",
	"LOINC":     "http://loinc.org",
	"SNOMED-CT": "http://snomed.info/sct",
	"RxNorm":    "http://www.nlm.nih.gov/research/umls/rxnorm/",
	"ICD-9-CM":  "http://hl7.org/fhir/sid/icd-9",
	"ICD-10-CM": "http://hl7.org/fhir/sid/icd-10",
}
