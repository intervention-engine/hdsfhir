package models

import (
	"encoding/json"

	fhir "github.com/intervention-engine/fhir/models"
)

type Procedure struct {
	Entry
	Description string                     `json:"description"`
	Values      []ResultValue              `json:"values"`
	Report      UploadableDiagnosticReport `json:"-"`
}

func (self *Procedure) UploadResults(baseURL string) {
	observations := self.FHIRObservationModels()
	if len(observations) > 0 {
		report := fhir.DiagnosticReport{Result: make([]fhir.Reference, len(observations))}
		for i, observation := range observations {
			uploadable := UploadableObservation{FHIRObservation: observation}
			Upload(&uploadable, baseURL+"/Observation")
			report.Result[i] = fhir.Reference{Reference: uploadable.ServerURL}
		}

		self.Report = UploadableDiagnosticReport{FHIRDiagnosticReport: report}
		Upload(&self.Report, baseURL+"/DiagnosticReport")
	}
}

func (self *Procedure) FHIRModel() fhir.Procedure {
	fhirProcedure := fhir.Procedure{}
	fhirProcedure.Type = self.Codes.FHIRCodeableConcept(self.Description)
	fhirProcedure.Encounter = &fhir.Reference{Reference: self.Patient.MatchingEncounter(self.Entry).ServerURL}
	fhirProcedure.Notes = self.Description
	fhirProcedure.Date = self.GetFHIRPeriod()
	if self.Report.ServerURL != "" {
		fhirProcedure.Report = append(fhirProcedure.Report, fhir.Reference{Reference: self.Report.ServerURL})
	}

	fhirProcedure.Subject = &fhir.Reference{Reference: self.Patient.ServerURL}
	return fhirProcedure
}

func (self *Procedure) FHIRObservationModels() []fhir.Observation {
	observations := make([]fhir.Observation, len(self.Values))

	for i, value := range self.Values {
		observation := value.FHIRModel()
		observation.Name = self.Codes.FHIRCodeableConcept(self.Description)
		observation.AppliesPeriod = self.GetFHIRPeriod()
		observation.Subject = &fhir.Reference{Reference: self.Patient.ServerURL}
		observations[i] = observation
	}

	return observations
}

func (self *Procedure) ToJSON() []byte {
	fhirProcedure := self.FHIRModel()
	json, _ := json.Marshal(fhirProcedure)
	return json
}
