package hdsfhir

import (
	"net/url"
	"sort"
	"strings"

	"github.com/intervention-engine/fhir/models"
)

// ConvertToConditionalUpdates converts a bundle containing POST requests to a bundle with PUT requests using
// conditional updates.  For patient resources, the update is based on the Medical Record Number.  For all other
// resources it is based on reasonable indicators of sameness (such as equal dates and codes).
func ConvertToConditionalUpdates(bundle *models.Bundle) error {
	for _, entry := range bundle.Entry {
		values := url.Values{}
		switch t := entry.Resource.(type) {
		case *models.AllergyIntolerance:
			if check(t.Patient, t.Substance, t.Onset) {
				values.Set("patient", refToParamValue(t.Patient, "Patient"))
				values.Set("code", ccToParamValue(t.Substance))
				values.Set("onset", dateToParamValue(t.Onset))
			}
		case *models.Condition:
			if check(t.Patient, t.Code, t.OnsetDateTime) {
				values.Set("patient", refToParamValue(t.Patient, "Patient"))
				values.Set("code", ccToParamValue(t.Code))
				values.Set("onset", dateToParamValue(t.OnsetDateTime))
			}
		case *models.DiagnosticReport:
			// TODO: Consider if this query is precise enough, consider searching on results too
			if check(t.Subject, t.Code, t.EffectivePeriod) {
				values.Set("patient", refToParamValue(t.Subject, "Patient"))
				values.Set("code", ccToParamValue(t.Code))
				values.Set("date", dateToParamValue(t.EffectivePeriod.Start))
			}
		case *models.Encounter:
			if check(t.Patient, t.Type, t.Period) {
				values.Set("patient", refToParamValue(t.Patient, "Patient"))
				for _, cc := range t.Type {
					values.Add("type", ccToParamValue(&cc))
				}
				// TODO: the date param references "a date within the period the encounter lasted."  Is this OK?
				values.Set("date", dateToParamValue(t.Period.Start))
			}
		case *models.Immunization:
			if check(t.Patient, t.VaccineCode, t.Date) {
				values.Set("patient", refToParamValue(t.Patient, "Patient"))
				values.Set("vaccine-code", ccToParamValue(t.VaccineCode))
				values.Set("date", dateToParamValue(t.Date))
			}
		case *models.MedicationStatement:
			if check(t.Patient, t.MedicationCodeableConcept, t.EffectivePeriod) {
				values.Set("patient", refToParamValue(t.Patient, "Patient"))
				values.Set("code", ccToParamValue(t.MedicationCodeableConcept))
				values.Set("effectivedate", dateToParamValue(t.EffectivePeriod.Start))
			}
		case *models.Observation:
			if check(t.Subject, t.Code, t.EffectivePeriod) {
				values.Set("patient", refToParamValue(t.Subject, "Patient"))
				values.Set("code", ccToParamValue(t.Code))
				values.Set("date", dateToParamValue(t.EffectivePeriod.Start))
			}
		case *models.Procedure:
			if check(t.Subject, t.Code, t.PerformedPeriod) {
				values.Set("patient", refToParamValue(t.Subject, "Patient"))
				values.Set("code", ccToParamValue(t.Code))
				values.Set("date", dateToParamValue(t.PerformedPeriod.Start))
			}
		case *models.ProcedureRequest:
			// We can't do anything meaningful because ProcedureRequest does not have search params
			// for code or orderedOn.  We simply can't get precise enough for a conditional update.
		case *models.Patient:
			if len(t.Identifier) > 0 && t.Identifier[0].Value != "" {
				values.Set("identifier", t.Identifier[0].Value)
			}
		}
		if entry.Request.Method == "POST" && len(values) > 0 {
			entry.Request.Method = "PUT"
			entry.Request.Url += "?" + values.Encode()
		}
	}
	return nil
}

func check(things ...interface{}) bool {
	for _, t := range things {
		switch t := t.(type) {
		case *models.CodeableConcept:
			if t == nil || len(t.Coding) == 0 {
				return false
			}
		case []models.CodeableConcept:
			if len(t) == 0 || !check(&t[0]) {
				return false
			}
		case *models.FHIRDateTime:
			if t == nil || t.Time.IsZero() {
				return false
			}
		case *models.Period:
			if t == nil || !check(t.Start) {
				return false
			}
		case *models.Reference:
			if t == nil || t.Reference == "" {
				return false
			}
		}
	}
	return true
}

func ccToParamValue(cc *models.CodeableConcept) string {
	codes := make([]string, len(cc.Coding))
	for i := range cc.Coding {
		codes[i] = cc.Coding[i].System + "|" + cc.Coding[i].Code
	}
	// sort for predictability (a.k.a., easier testing)
	sort.Strings(codes)
	return strings.Join(codes, ",")
}

func dateToParamValue(date *models.FHIRDateTime) string {
	return date.Time.Format("2006-01-02T15:04:05")
}

func refToParamValue(ref *models.Reference, resourceType string) string {
	return resourceType + "/" + strings.TrimPrefix(ref.Reference, "urn:uuid:")
}
