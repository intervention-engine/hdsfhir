package hdsfhir

import (
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"

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
				addRefParam(values, "patient", t.Patient)
				addCCParam(values, "substance", t.Substance)
				addDateParam(values, "onset", t.Onset)
			}
		case *models.Condition:
			if check(t.Patient, t.Code, t.OnsetDateTime) {
				addRefParam(values, "patient", t.Patient)
				addCCParam(values, "code", t.Code)
				addDateParam(values, "onset", t.OnsetDateTime)
			}
		case *models.DiagnosticReport:
			// TODO: Consider if this query is precise enough, consider searching on results too
			if check(t.Subject, t.Code, t.EffectivePeriod) {
				addRefParam(values, "patient", t.Subject)
				addCCParam(values, "code", t.Code)
				addPeriodParam(values, "date", t.EffectivePeriod)
			}
		case *models.Encounter:
			if check(t.Patient, t.Type, t.Period) {
				addRefParam(values, "patient", t.Patient)
				for _, cc := range t.Type {
					addCCParam(values, "type", &cc)
				}
				// TODO: the date param references "a date within the period the encounter lasted."  Is this OK?
				addPeriodParam(values, "date", t.Period)
			}
		case *models.Immunization:
			if check(t.Patient, t.VaccineCode, t.Date) {
				addRefParam(values, "patient", t.Patient)
				addCCParam(values, "vaccine-code", t.VaccineCode)
				addDateParam(values, "date", t.Date)
			}
		case *models.MedicationStatement:
			if check(t.Patient, t.MedicationCodeableConcept, t.EffectivePeriod) {
				addRefParam(values, "patient", t.Patient)
				addCCParam(values, "code", t.MedicationCodeableConcept)
				addPeriodParam(values, "effectivedate", t.EffectivePeriod)
			}
		case *models.Observation:
			if check(t.Subject, t.Code, t.EffectivePeriod) {
				addRefParam(values, "patient", t.Subject)
				addCCParam(values, "code", t.Code)
				addPeriodParam(values, "date", t.EffectivePeriod)
				if check(t.ValueCodeableConcept) {
					addCCParam(values, "value-concept", t.ValueCodeableConcept)
				} else if check(t.ValueQuantity) {
					q := t.ValueQuantity
					if q.Code != "" {
						values.Add("value-quantity", fmt.Sprintf("%g|%s|%s", *q.Value, q.System, q.Code))
					} else {
						values.Add("value-quantity", fmt.Sprintf("%g|%s|%s", *q.Value, q.System, q.Unit))
					}
				} else if check(t.ValueString) {
					values.Add("value-string", t.ValueString)
				}
			}
		case *models.Procedure:
			if check(t.Subject, t.Code, t.PerformedPeriod) {
				addRefParam(values, "patient", t.Subject)
				addCCParam(values, "code", t.Code)
				addPeriodParam(values, "date", t.PerformedPeriod)
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
		case *models.Quantity:
			if t == nil || t.Value == nil {
				return false
			}
		case *models.Reference:
			if t == nil || t.Reference == "" {
				return false
			}
		case string:
			if t == "" {
				return false
			}
		}
	}
	return true
}

func addCCParam(values url.Values, name string, cc *models.CodeableConcept) {
	codes := make([]string, len(cc.Coding))
	for i := range cc.Coding {
		codes[i] = cc.Coding[i].System + "|" + cc.Coding[i].Code
	}
	// sort for predictability (a.k.a., easier testing)
	sort.Strings(codes)
	values.Add(name, strings.Join(codes, ","))
}

func addDateParam(values url.Values, name string, date *models.FHIRDateTime) {
	values.Add(name, date.Time.Format("2006-01-02T15:04:05"))
}

func addPeriodParam(values url.Values, name string, period *models.Period) {
	// Due to the way searching on periods is defined, the only way we can try to match on the start date is by using
	// a combination of sa (starts after) and lt (less than) query parameters.
	l := period.Start.Time.Add(-1 * time.Second)
	h := period.Start.Time.Add(1 * time.Second)
	values.Add(name, "sa"+l.Format("2006-01-02T15:04:05"))
	values.Add(name, "lt"+h.Format("2006-01-02T15:04:05"))
}

func addRefParam(values url.Values, name string, ref *models.Reference) {
	values.Add(name, ref.Reference)
}
