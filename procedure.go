package hdsfhir

import fhir "github.com/intervention-engine/fhir/models"

type Procedure struct {
	Entry
	AnatomicalTarget *CodeObject   `json:"anatomical_target"`
	Values           []ResultValue `json:"values"`
}

func (p *Procedure) FHIRModels() []interface{} {
	if p.isProcedureRequest() {
		return p.convertProcedureRequest()
	}

	return p.convertProcedure()
}

func (p *Procedure) isProcedureRequest() bool {
	statusConcept := p.StatusCode.FHIRCodeableConcept("")
	switch {
	case statusConcept.MatchesCode("http://hl7.org/fhir/ValueSet/v3-ActStatus", "cancelled"):
		return true
	case statusConcept.MatchesCode("http://hl7.org/fhir/ValueSet/v3-ActStatus", "held"):
		return true
	case statusConcept.MatchesCode("http://hl7.org/fhir/ValueSet/v3-ActStatus", "new"):
		return true
	case statusConcept.MatchesCode("http://hl7.org/fhir/ValueSet/v3-ActStatus", "suspended"):
		return true
	case statusConcept.MatchesCode("http://hl7.org/fhir/ValueSet/v3-ActStatus", "nullified"):
		return true
	// NOTE: this is not a real ActStatus, but HDS seems to use it
	case statusConcept.MatchesCode("http://hl7.org/fhir/ValueSet/v3-ActStatus", "ordered"):
		return true
	// NOTE: this is not a real ActStatus, but HDS seems to use it
	case statusConcept.MatchesCode("http://hl7.org/fhir/ValueSet/v3-ActStatus", "recommended"):
		return true
	case p.MoodCode == "RQO":
		return true
	}

	return false
}

func (p *Procedure) convertProcedure() []interface{} {
	fhirProcedure := &fhir.Procedure{}
	fhirProcedure.Id = p.GetTempID()
	fhirProcedure.Subject = p.Patient.FHIRReference()
	fhirProcedure.Status = p.convertProcedureStatus()
	fhirProcedure.Code = p.Codes.FHIRCodeableConcept(p.Description)
	if p.NegationInd {
		t := true
		fhirProcedure.NotPerformed = &t
	}
	if p.NegationReason != nil {
		cc := p.NegationReason.FHIRCodeableConcept("")
		fhirProcedure.ReasonNotPerformed = []fhir.CodeableConcept{*cc}
	}
	if p.AnatomicalTarget != nil {
		cc := p.AnatomicalTarget.FHIRCodeableConcept("")
		fhirProcedure.BodySite = []fhir.CodeableConcept{*cc}
	}
	fhirProcedure.PerformedPeriod = p.GetFHIRPeriod()
	fhirProcedure.Encounter = p.Patient.MatchingEncounterReference(p.Entry)

	models := []interface{}{fhirProcedure}
	if len(p.Values) > 0 {
		// Create the diagnostic report model with its own ID and slots for results
		internalReportID := &TemporallyIdentified{}
		fhirReport := &fhir.DiagnosticReport{}
		fhirReport.Id = internalReportID.GetTempID()
		fhirReport.Status = "final"
		fhirReport.Code = &fhir.CodeableConcept{
			Coding: []fhir.Coding{
				{
					System:  "http://loinc.org",
					Code:    "59776-5",
					Display: "Procedure findings narrative",
				},
			},
			Text: "Procedure findings narrative",
		}
		fhirReport.Subject = p.Patient.FHIRReference()
		fhirReport.Encounter = p.Patient.MatchingEncounterReference(p.Entry)
		fhirReport.EffectivePeriod = p.GetFHIRPeriod()
		fhirReport.Issued = p.EndTime.FHIRDateTime() // Not perfect, but it's a required field
		// TODO: Technically, "performer" is required, but we don't want to make up data
		fhirReport.Result = make([]fhir.Reference, len(p.Values))
		models = append(models, fhirReport)

		// Link the procedure to the report
		fhirProcedure.Report = []fhir.Reference{*internalReportID.FHIRReference()}

		// Create the observation values
		for i := range p.Values {
			observation := p.Values[i].FHIRModels()[0].(*fhir.Observation)
			observation.Code = p.Codes.FHIRCodeableConcept(p.Description)
			observation.Subject = p.Patient.FHIRReference()
			observation.Encounter = p.Patient.MatchingEncounterReference(p.Entry)
			observation.EffectivePeriod = p.GetFHIRPeriod()
			models = append(models, observation)

			// Link the report results to the observation
			fhirReport.Result[i] = *p.Values[i].FHIRReference()
		}
	}

	return models
}

// convertProcedureStatus maps the status to a code in the required FHIR value set:
//   http://hl7.org/fhir/DSTU2/valueset-procedure-status.html
func (p *Procedure) convertProcedureStatus() string {
	var status string
	statusConcept := p.StatusCode.FHIRCodeableConcept("")
	switch {
	case statusConcept.MatchesCode("http://hl7.org/fhir/ValueSet/v3-ActStatus", "aborted"):
		status = "aborted"
	case statusConcept.MatchesCode("http://hl7.org/fhir/ValueSet/v3-ActStatus", "active"):
		status = "in-progress"
	case statusConcept.MatchesCode("http://hl7.org/fhir/ValueSet/v3-ActStatus", "obsolete"):
		status = "entered-in-error"
	case len(p.StatusCode) == 0 && p.EndTime == 0:
		status = "in-progress"
	default:
		status = "completed"
	}

	return status
}

func (p *Procedure) convertProcedureRequest() []interface{} {
	fhirProcedureRequest := &fhir.ProcedureRequest{}
	fhirProcedureRequest.Id = p.GetTempID()
	fhirProcedureRequest.Subject = p.Patient.FHIRReference()
	fhirProcedureRequest.Status = p.convertProcedureRequestStatus()
	fhirProcedureRequest.Code = p.Codes.FHIRCodeableConcept(p.Description)
	if p.AnatomicalTarget != nil {
		cc := p.AnatomicalTarget.FHIRCodeableConcept("")
		fhirProcedureRequest.BodySite = []fhir.CodeableConcept{*cc}
	}
	if p.Time > 0 {
		fhirProcedureRequest.OrderedOn = p.Time.FHIRDateTime()
	} else if p.StartTime > 0 {
		fhirProcedureRequest.OrderedOn = p.StartTime.FHIRDateTime()
	} else if p.EndTime > 0 {
		fhirProcedureRequest.OrderedOn = p.EndTime.FHIRDateTime()
	}
	fhirProcedureRequest.Encounter = p.Patient.MatchingEncounterReference(p.Entry)

	return []interface{}{fhirProcedureRequest}
}

// convertProcedureRequestStatus maps the status to a code in the required FHIR value set:
//   http://hl7.org/fhir/DSTU2/valueset-procedure-request-status.html
func (p *Procedure) convertProcedureRequestStatus() string {
	var status string
	statusConcept := p.StatusCode.FHIRCodeableConcept("")
	switch {
	case p.NegationInd == true:
		status = "rejected"
	case statusConcept.MatchesCode("http://hl7.org/fhir/ValueSet/v3-ActStatus", "cancelled"):
		status = "rejected"
	case statusConcept.MatchesCode("http://hl7.org/fhir/ValueSet/v3-ActStatus", "held"):
		status = "suspended"
	case statusConcept.MatchesCode("http://hl7.org/fhir/ValueSet/v3-ActStatus", "suspended"):
		status = "suspended"
	case statusConcept.MatchesCode("http://hl7.org/fhir/ValueSet/v3-ActStatus", "nullified"):
		status = "rejected"
	// NOTE: this is not a real ActStatus, but HDS seems to use it
	case statusConcept.MatchesCode("http://hl7.org/fhir/ValueSet/v3-ActStatus", "recommended"):
		status = "proposed"
	default:
		status = "accepted"
	}

	return status
}
