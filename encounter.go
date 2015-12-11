package hdsfhir

import fhir "github.com/intervention-engine/fhir/models"

type Encounter struct {
	Entry
	Reason               *Entry      `json:"reason"`
	DischargeDisposition *CodeObject `json:"dischargeDisposition"`
}

func (e *Encounter) FHIRModels() []interface{} {
	fhirEncounter := &fhir.Encounter{Id: e.GetTempID()}
	fhirEncounter.Status = e.convertStatus()
	typeConcept := e.Codes.FHIRCodeableConcept(e.Description)
	fhirEncounter.Type = []fhir.CodeableConcept{*typeConcept}
	fhirEncounter.Patient = e.Patient.FHIRReference()
	fhirEncounter.Period = e.GetFHIRPeriod()
	if e.Reason != nil && len(e.Reason.Codes) > 0 {
		reasonConcept := e.Reason.Codes.FHIRCodeableConcept("")
		fhirEncounter.Reason = []fhir.CodeableConcept{*reasonConcept}
	}
	if e.DischargeDisposition != nil {
		fhirEncounter.Hospitalization = &fhir.EncounterHospitalizationComponent{
			DischargeDisposition: e.DischargeDisposition.FHIRCodeableConcept(""),
		}
	}

	return []interface{}{fhirEncounter}
}

// convertStatus maps the status to a code in the required FHIR value set:
//   http://hl7.org/fhir/DSTU2/valueset-encounter-state.html
// If the status cannot be reliably mapped, "finished" will be assumed.  Note that this code is
// built to handle even some statuses that HDS does not currently return (active, cancelled, etc.)
func (e *Encounter) convertStatus() string {
	var status string
	statusConcept := e.StatusCode.FHIRCodeableConcept("")
	switch {
	// Negated encounters are rare, but if we run into one, call it cancelled
	case e.NegationInd:
		status = "cancelled"
	case statusConcept.MatchesCode("http://hl7.org/fhir/ValueSet/v3-ActStatus", "active"):
		status = "in-progress"
	case statusConcept.MatchesCode("http://hl7.org/fhir/ValueSet/v3-ActStatus", "cancelled"):
		status = "cancelled"
	case statusConcept.MatchesCode("http://hl7.org/fhir/ValueSet/v3-ActStatus", "held"):
		status = "planned"
	case statusConcept.MatchesCode("http://hl7.org/fhir/ValueSet/v3-ActStatus", "new"):
		status = "planned"
	case statusConcept.MatchesCode("http://hl7.org/fhir/ValueSet/v3-ActStatus", "suspended"):
		status = "onleave"
	case statusConcept.MatchesCode("http://hl7.org/fhir/ValueSet/v3-ActStatus", "nullified"):
		status = "cancelled"
	case statusConcept.MatchesCode("http://hl7.org/fhir/ValueSet/v3-ActStatus", "obsolete"):
		status = "cancelled"
	// NOTE: this is not a real ActStatus, but HDS seems to use it
	case statusConcept.MatchesCode("http://hl7.org/fhir/ValueSet/v3-ActStatus", "ordered"):
		status = "planned"
	case e.MoodCode == "RQO":
		status = "planned"
	default:
		status = "finished"
	}

	return status
}
