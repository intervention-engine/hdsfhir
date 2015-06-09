package models

import fhir "github.com/intervention-engine/fhir/models"

type CodeMap map[string][]string

func (c *CodeMap) FHIRCodeableConcept(text string) *fhir.CodeableConcept {
	concept := &fhir.CodeableConcept{}
	codings := make([]fhir.Coding, 0)
	for codeSystem, codes := range *c {
		codeSystemURL := CodeSystemMap[codeSystem]
		for _, code := range codes {
			coding := fhir.Coding{System: codeSystemURL, Code: code}
			codings = append(codings, coding)
		}
	}
	concept.Coding = codings
	concept.Text = text
	return concept
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
	"NDC":        "http://www.fda.gov/Drugs/InformationOnDrugs",
	"CVX":        "http://www2a.cdc.gov/vaccines/iis/iisstandards/vaccines.asp?rpt=cvx",
}
