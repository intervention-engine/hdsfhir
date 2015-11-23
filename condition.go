package hdsfhir

import fhir "github.com/intervention-engine/fhir/models"

type Condition struct {
	Entry
	Severity CodeMap `json:"severity"`
}

func (c *Condition) FHIRModels() []interface{} {
	fhirCondition := &fhir.Condition{Id: c.GetTempID()}
	fhirCondition.Patient = c.Patient.FHIRReference()
	fhirCondition.Code = c.Codes.FHIRCodeableConcept(c.Description)
	// TODO: consider setting category to "diagnosis"
	fhirCondition.ClinicalStatus = c.convertClinicalStatus()
	if c.NegationInd {
		fhirCondition.VerificationStatus = "refuted"
	} else {
		fhirCondition.VerificationStatus = "confirmed"
	}
	fhirCondition.Severity = c.convertSeverity()
	if c.StartTime != 0 {
		fhirCondition.OnsetDateTime = c.StartTime.FHIRDateTime()
	}
	if c.EndTime != 0 {
		fhirCondition.AbatementDateTime = c.EndTime.FHIRDateTime()
	}

	return []interface{}{fhirCondition}
}

// convertClinicalStatus maps the clinical status to a code in the "preferred" FHIR value set:
//   http://hl7.org/fhir/DSTU2/valueset-condition-clinical.html
// If the status cannot be reliably mapped, an empty code will be returned.
func (c *Condition) convertClinicalStatus() string {
	var status string
	statusConcept := c.StatusCode.FHIRCodeableConcept("")
	switch {
	case statusConcept.MatchesCode("http://snomed.info/sct", "55561003"):
		status = "active"
	case statusConcept.MatchesCode("http://snomed.info/sct", "73425007"):
		status = "remission"
	case statusConcept.MatchesCode("http://snomed.info/sct", "413322009"):
		status = "resolved"
	case statusConcept.MatchesCode("http://hl7.org/fhir/ValueSet/v3-ActStatus", "active"):
		status = "active"
	}

	// In order to remain consistent, fix the status if there is an end date (abatement)
	if c.EndTime != 0 {
		if status == "" {
			status = "resolved"
		} else if status == "active" {
			status = "remission" // closest concept to active with an abatement
		}
	}

	return status
}

// convertSeverity maps the severity to a CodeableConcept. If possible, it will add a display name.
// FHIR has a "preferred" value set for severity:
//   http://hl7.org/fhir/DSTU2/valueset-condition-severity.html
// Note that this is a subset of the value set that is often used:
//   https://phinvads.cdc.gov/vads/ViewValueSet.action?id=72FDBFB5-A277-DE11-9B52-0015173D1785
// Rather than make a subjective mapping decision, we keep the codes as-is, which is likely
// more in line with the larger value set (containing "mild to moderate" and "moderate to severe")
func (c *Condition) convertSeverity() *fhir.CodeableConcept {
	if len(c.Severity) == 0 {
		return nil
	}

	severity := c.Severity.FHIRCodeableConcept("")
	switch {
	case severity.MatchesCode("http://snomed.info/sct", "399166001"):
		severity.Text = "Fatal"
	case severity.MatchesCode("http://snomed.info/sct", "255604002"):
		severity.Text = "Mild"
	case severity.MatchesCode("http://snomed.info/sct", "371923003"):
		severity.Text = "Mild to moderate"
	case severity.MatchesCode("http://snomed.info/sct", "6736007"):
		severity.Text = "Moderate"
	case severity.MatchesCode("http://snomed.info/sct", "371924009"):
		severity.Text = "Moderate to severe"
	case severity.MatchesCode("http://snomed.info/sct", "24484000"):
		severity.Text = "Severe"
	}

	return severity
}
