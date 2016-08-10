package hdsfhir

import fhir "github.com/intervention-engine/fhir/models"

type Allergy struct {
	Entry
	Reaction *CodeObject `json:"reaction"`
	Severity *CodeObject `json:"severity"`
}

func (a *Allergy) FHIRModels() []interface{} {
	fhirAllergy := &fhir.AllergyIntolerance{}
	fhirAllergy.Id = a.GetTempID()
	if a.StartTime != nil {
		fhirAllergy.Onset = a.StartTime.FHIRDateTime()
	} else if a.Time != nil {
		fhirAllergy.Onset = a.Time.FHIRDateTime()
	}
	fhirAllergy.Patient = a.Patient.FHIRReference()
	fhirAllergy.Substance = a.Codes.FHIRCodeableConcept(a.Description)
	fhirAllergy.Status = a.convertStatus()
	fhirAllergy.Criticality = a.convertCriticality()
	if a.Reaction != nil {
		cc := a.Reaction.FHIRCodeableConcept("")
		fhirAllergy.Reaction = []fhir.AllergyIntoleranceReactionComponent{
			{
				Substance:     fhirAllergy.Substance,
				Manifestation: []fhir.CodeableConcept{*cc},
				Onset:         fhirAllergy.Onset,
				Severity:      a.convertSeverity(),
			},
		}
	}

	return []interface{}{fhirAllergy}
}

// convertStatus maps the status to a code in the "required" FHIR value set:
//   http://hl7.org/fhir/DSTU2/valueset-allergy-intolerance-status.html
// If the status cannot be reliably mapped, active is assumed.
func (a *Allergy) convertStatus() string {
	var status string

	if a.NegationInd {
		return "refuted"
	}

	statusConcept := a.StatusCode.FHIRCodeableConcept("")
	switch {
	case statusConcept.MatchesCode("http://snomed.info/sct", "55561003"):
		status = "active"
	case statusConcept.MatchesCode("http://snomed.info/sct", "73425007"):
		status = "inactive"
	case statusConcept.MatchesCode("http://snomed.info/sct", "413322009"):
		status = "resolved"
	default:
		status = "active"
	}

	// In order to remain consistent, fix the status if there is an end date
	// and it doesn't match the start date (in which case we can't be sure it's really an end)
	if status == "" && a.EndTime != nil && a.StartTime != nil && *a.EndTime != *a.StartTime {
		status = "resolved"
	}

	return status
}

// convertCriticality maps the severity to a CodeableConcept. FHIR has a "required" value set for
// criticality:
//   http://hl7.org/fhir/DSTU2/valueset-allergy-intolerance-criticality.html
// If the severity can't be mapped, criticality will be left blank
func (a *Allergy) convertCriticality() string {
	if a.Severity == nil {
		return ""
	}

	criticality := a.Severity.FHIRCodeableConcept("")
	switch {
	case criticality.MatchesCode("http://snomed.info/sct", "399166001"):
		return "CRITH"
	case criticality.MatchesCode("http://snomed.info/sct", "255604002"):
		return "CRITL"
	case criticality.MatchesCode("http://snomed.info/sct", "371923003"):
		// Mild to moderate: translate to L
		return "CRITL"
	case criticality.MatchesCode("http://snomed.info/sct", "6736007"):
		// Moderate: tough to call L or H, translate to CRITU (unable to determine)
		return "CRITU"
	case criticality.MatchesCode("http://snomed.info/sct", "371924009"):
		// Moderate to severe: err on the side of safety, translate to CRITH
		return "CRITH"
	case criticality.MatchesCode("http://snomed.info/sct", "24484000"):
		return "CRITH"
	}

	return ""
}

// convertSeverity maps the severity to a CodeableConcept. FHIR has a "required" value set for
// severity:
//   http://hl7.org/fhir/DSTU2/valueset-reaction-event-severity.html
// If the severity can't be mapped, severity will be left blank
func (a *Allergy) convertSeverity() string {
	if a.Severity == nil {
		return ""
	}

	severity := a.Severity.FHIRCodeableConcept("")
	switch {
	case severity.MatchesCode("http://snomed.info/sct", "399166001"):
		return "severe"
	case severity.MatchesCode("http://snomed.info/sct", "255604002"):
		return "mild"
	case severity.MatchesCode("http://snomed.info/sct", "371923003"):
		// Mild to moderate: translate to moderate
		return "moderate"
	case severity.MatchesCode("http://snomed.info/sct", "6736007"):
		return "moderate"
	case severity.MatchesCode("http://snomed.info/sct", "371924009"):
		// Moderate to severe: err on the side of safety, translate to severe
		return "severe"
	case severity.MatchesCode("http://snomed.info/sct", "24484000"):
		return "severe"
	}

	return ""
}
