package hdsfhir

import fhir "github.com/intervention-engine/fhir/models"

type Medication struct {
	Entry
}

func (m *Medication) FHIRModels() []interface{} {
	// Sometimes immunizations come across as medications, so we need to support that
	_, isImmunization := m.Codes["CVX"]
	if isImmunization {
		return m.convertImmunization()
	}

	return m.convertMedication()
}

func (m *Medication) convertMedication() []interface{} {
	// TODO: Use more specific medication resources (MedicationOrder, MedicationAdministration, etc.)
	fhirMedicationStatement := &fhir.MedicationStatement{Id: m.GetTempID()}
	fhirMedicationStatement.Patient = m.Patient.FHIRReference()
	fhirMedicationStatement.Status = m.convertMedicationStatus()
	if m.NegationInd {
		t := true
		fhirMedicationStatement.WasNotTaken = &t
	}
	if len(m.NegationReason) > 0 {
		cc := m.NegationReason.FHIRCodeableConcept("")
		fhirMedicationStatement.ReasonNotTaken = []fhir.CodeableConcept{*cc}
	}
	fhirMedicationStatement.EffectivePeriod = m.GetFHIRPeriod()
	fhirMedicationStatement.MedicationCodeableConcept = m.Codes.FHIRCodeableConcept(m.Description)

	// Ignoring dosage, route, etc.

	return []interface{}{fhirMedicationStatement}
}

// convertMedicationStatus maps the status to a code in the required FHIR value set:
//   http://hl7.org/fhir/DSTU2/valueset-medication-statement-status.html
func (m *Medication) convertMedicationStatus() string {
	var status string
	statusConcept := m.StatusCode.FHIRCodeableConcept("")
	switch {
	case statusConcept.MatchesCode("http://hl7.org/fhir/ValueSet/v3-ActStatus", "active"):
		status = "active"
	case statusConcept.MatchesCode("http://hl7.org/fhir/ValueSet/v3-ActStatus", "cancelled"):
		status = "entered-in-error"
	case statusConcept.MatchesCode("http://hl7.org/fhir/ValueSet/v3-ActStatus", "held"):
		status = "intended"
	case statusConcept.MatchesCode("http://hl7.org/fhir/ValueSet/v3-ActStatus", "new"):
		status = "intended"
	case statusConcept.MatchesCode("http://hl7.org/fhir/ValueSet/v3-ActStatus", "suspended"):
		status = "completed"
	case statusConcept.MatchesCode("http://hl7.org/fhir/ValueSet/v3-ActStatus", "nullified"):
		status = "entered-in-error"
	case statusConcept.MatchesCode("http://hl7.org/fhir/ValueSet/v3-ActStatus", "obsolete"):
		status = "cancelled"
	// NOTE: this is not a real ActStatus, but HDS seems to use it
	case statusConcept.MatchesCode("http://hl7.org/fhir/ValueSet/v3-ActStatus", "ordered"):
		status = "intended"
	// NOTE: this is not a real ActStatus, but HDS seems to use it
	case statusConcept.MatchesCode("http://hl7.org/fhir/ValueSet/v3-ActStatus", "discharge"):
		status = "intended"
	// NOTE: this is not a real ActStatus, but HDS seems to use it
	case statusConcept.MatchesCode("http://hl7.org/fhir/ValueSet/v3-ActStatus", "dispensed"):
		status = "intended"
	case m.MoodCode == "RQO":
		status = "intended"
	case len(m.StatusCode) == 0 && m.EndTime == 0:
		status = "active"
	default:
		status = "completed"
	}

	return status
}

func (m *Medication) convertImmunization() []interface{} {
	fhirImmunization := &fhir.Immunization{Id: m.GetTempID()}
	fhirImmunization.Status = m.convertImmunizationStatus()
	fhirImmunization.Date = m.StartTime.FHIRDateTime()
	fhirImmunization.VaccineCode = m.Codes.FHIRCodeableConcept(m.Description)
	fhirImmunization.Patient = m.Patient.FHIRReference()
	if m.NegationInd {
		t := true
		fhirImmunization.WasNotGiven = &t
	}
	if len(m.NegationReason) > 0 {
		cc := m.NegationReason.FHIRCodeableConcept("")
		fhirImmunization.Explanation = &fhir.ImmunizationExplanationComponent{
			ReasonNotGiven: []fhir.CodeableConcept{*cc},
		}
	}

	// Ignoring dosage, route, etc.

	return []interface{}{fhirImmunization}
}

// convertImmunizationStatus maps the status to a code in the required FHIR value set:
//   http://hl7.org/fhir/DSTU2/valueset-medication-admin-status.html
func (m *Medication) convertImmunizationStatus() string {
	var status string
	statusConcept := m.StatusCode.FHIRCodeableConcept("")
	switch {
	case statusConcept.MatchesCode("http://hl7.org/fhir/ValueSet/v3-ActStatus", "aborted"):
		status = "stopped"
	case statusConcept.MatchesCode("http://hl7.org/fhir/ValueSet/v3-ActStatus", "active"):
		status = "in-progress"
	case statusConcept.MatchesCode("http://hl7.org/fhir/ValueSet/v3-ActStatus", "cancelled"):
		status = "entered-in-error"
	case statusConcept.MatchesCode("http://hl7.org/fhir/ValueSet/v3-ActStatus", "held"):
		status = "on-hold"
	case statusConcept.MatchesCode("http://hl7.org/fhir/ValueSet/v3-ActStatus", "new"):
		status = "intended"
	case statusConcept.MatchesCode("http://hl7.org/fhir/ValueSet/v3-ActStatus", "suspended"):
		status = "on-hold"
	case statusConcept.MatchesCode("http://hl7.org/fhir/ValueSet/v3-ActStatus", "nullified"):
		status = "entered-in-error"
	case statusConcept.MatchesCode("http://hl7.org/fhir/ValueSet/v3-ActStatus", "obsolete"):
		status = "entered-in-error"
	// NOTE: this is not a real ActStatus, but HDS seems to use it
	case statusConcept.MatchesCode("http://hl7.org/fhir/ValueSet/v3-ActStatus", "ordered"):
		status = "intended"
	case m.MoodCode == "RQO":
		status = "intended"
	case len(m.StatusCode) == 0 && m.EndTime == 0:
		status = "in-progress"
	default:
		status = "completed"
	}

	return status
}
