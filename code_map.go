package hdsfhir

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

type CodeObject struct {
	Code       string `json:"code"`
	CodeSystem string `json:"codeSystem"`
}

func (c *CodeObject) FHIRCodeableConcept(text string) *fhir.CodeableConcept {
	concept := &fhir.CodeableConcept{}
	concept.Coding = []fhir.Coding{
		{System: CodeSystemMap[c.CodeSystem], Code: c.Code},
	}
	concept.Text = text
	return concept
}

var CodeSystemMap = map[string]string{
	"CPT":                             "http://www.ama-assn.org/go/cpt",
	"LOINC":                           "http://loinc.org",
	"SNOMED-CT":                       "http://snomed.info/sct",
	"RxNorm":                          "http://www.nlm.nih.gov/research/umls/rxnorm",
	"ICD-9-CM":                        "http://hl7.org/fhir/sid/icd-9",
	"ICD-10-CM":                       "http://hl7.org/fhir/sid/icd-10",
	"ICD-9-PCS":                       "http://hl7.org/fhir/sid/icd-9",
	"ICD-10-PCS":                      "http://hl7.org/fhir/sid/icd-10",
	"NDC":                             "http://www.fda.gov/Drugs/InformationOnDrugs",
	"CVX":                             "http://www2a.cdc.gov/vaccines/iis/iisstandards/vaccines.asp?rpt=cvx",
	"HCP":                             "urn:oid:2.16.840.1.113883.6.14",
	"HCPCS":                           "urn:oid:2.16.840.1.113883.6.285",
	"HL7 Marital Status":              "http://hl7.org/fhir/ValueSet/v3-MaritalStatus",
	"HITSP C80 Observation Status":    "http://hl7.org/fhir/ValueSet/v3-ObservationInterpretation",
	"NCI Thesaurus":                   "urn:oid:2.16.840.1.113883.3.26.1.1",
	"FDA SPL":                         "urn:oid:2.16.840.1.113883.3.26.1.1",
	"FDA":                             "urn:oid:2.16.840.1.113883.3.88.12.80.20",
	"UNII":                            "http://fdasis.nlm.nih.gov",
	"HL7 ActStatus":                   "http://hl7.org/fhir/ValueSet/v3-ActStatus",
	"HL7 Healthcare Service Location": "urn:oid:2.16.840.1.113883.6.259",
	"HSLOC":                       "urn:oid:2.16.840.1.113883.6.259",
	"DischargeDisposition":        "urn:oid:2.16.840.1.113883.12.112",
	"HL7 Act Code":                "http://hl7.org/fhir/ValueSet/v3-ActCode",
	"HL7 Relationship Code":       "urn:oid:2.16.840.1.113883.1.11.18877",
	"CDC Race":                    "urn:oid:2.16.840.1.113883.6.238",
	"NLM Mesh":                    "urn:oid:2.16.840.1.113883.6.177",
	"Religious Affiliation":       "http://hl7.org/fhir/ValueSet/v3-ReligiousAffiliation",
	"HL7 ActNoImmunicationReason": "urn:oid:2.16.840.1.113883.1.11.19717",
	"NUBC": "urn:oid:2.16.840.1.113883.3.88.12.80.33",
	"HL7 Observation Interpretation": "urn:oid:2.16.840.1.113883.1.11.78",
	"Source of Payment Typology":     "urn:oid:2.16.840.1.113883.3.221.5",
	"SOP":               "urn:oid:2.16.840.1.113883.3.221.5",
	"CDT":               "urn:oid:2.16.840.1.113883.6.13",
	"AdministrativeSex": "urn:oid:2.16.840.1.113883.18.2",
}
