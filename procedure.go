package hdsfhir

import fhir "github.com/intervention-engine/fhir/models"

type Procedure struct {
	Entry
	Description string               `json:"description"`
	Values      []ResultValue        `json:"values"`
	report      TemporallyIdentified `json:"-"`
}

func (p *Procedure) FHIRModels() []interface{} {
	fhirProcedure := &fhir.Procedure{Id: p.GetTempID()}
	fhirProcedure.Type = p.Codes.FHIRCodeableConcept(p.Description)
	fhirProcedure.Encounter = p.Patient.MatchingEncounterReference(p.Entry)
	fhirProcedure.Notes = p.Description
	fhirProcedure.Date = p.GetFHIRPeriod()
	fhirProcedure.Subject = p.Patient.FHIRReference()

	models := []interface{}{fhirProcedure}
	if len(p.Values) > 0 {
		// Create the diagnostic report model with its own ID and slots for results
		internalReportID := &TemporallyIdentified{}
		fhirReport := &fhir.DiagnosticReport{Id: internalReportID.GetTempID()}
		fhirReport.Result = make([]fhir.Reference, len(p.Values))
		fhirReport.Subject = p.Patient.FHIRReference()
		models = append(models, fhirReport)

		// Link the procedure to the report
		fhirProcedure.Report = []fhir.Reference{*internalReportID.FHIRReference()}

		// Create the observation values
		for i := range p.Values {
			observation := p.Values[i].FHIRModels()[0].(*fhir.Observation)
			observation.Name = p.Codes.FHIRCodeableConcept(p.Description)
			observation.AppliesPeriod = p.GetFHIRPeriod()
			observation.Subject = p.Patient.FHIRReference()
			models = append(models, observation)

			// Link the report results to the observation
			fhirReport.Result[i] = *p.Values[i].FHIRReference()
		}
	}

	return models
}
