package hdsfhir

import (
	"time"
)

type FHIRCoding struct {
	System string `json:"system"`
	Code   string `json:"code"`
}

type Entry struct {
	StartTime   int64               `json:"start_time"`
	EndTime     int64               `json:"end_time"`
	Time        int64               `json:"time"`
	Oid         string              `json:"oid"`
	Codes       map[string][]string `json:"codes"`
	NegationInd bool                `json:"negationInd"`
	StatusCode  map[string][]string `json:"status_code"`
}

func (e *Entry) ConvertCodingToFHIR() []FHIRCoding {
	codings := make([]FHIRCoding, 3)
	for codeSystem, codes := range e.Codes {
		codeSystemURL := CodeSystemMap[codeSystem]
		for _, code := range codes {
			coding := FHIRCoding{System: codeSystemURL, Code: code}
			codings = append(codings, coding)
		}
	}
	return codings
}

func UnixToFHIRDate(unixTime int64) string {
	return time.Unix(unixTime, 0).Format("2006-01-02")
}

var CodeSystemMap = map[string]string{
	"CPT":       "http://www.ama-assn.org/go/cpt",
	"LOINC":     "http://loinc.org",
	"SNOMED-CT": "http://snomed.info/sct",
	"RxNorm":    "http://www.nlm.nih.gov/research/umls/rxnorm/",
	"ICD-9-CM":  "http://hl7.org/fhir/sid/icd-9",
	"ICD-10-CM": "http://hl7.org/fhir/sid/icd-10",
}
