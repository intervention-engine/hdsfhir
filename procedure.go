package hdsfhir

import (
	"encoding/json"

	"github.com/intervention-engine/fhir/models"
)

type Procedure struct {
	Entry
	Description        string                     `json:"description"`
	ResultObservations []UploadableObservation    `json:"-"`
	Report             UploadableDiagnosticReport `json:"-"`
	ThingWithResults
}

func (self *Procedure) UploadResults(baseURL string) {
	self.ProcessResultObservations()
	if len(self.ResultObservations) > 0 {
		for i := 0; i < len(self.ResultObservations); i++ {
			current := &self.ResultObservations[i]
			Upload(current, baseURL+"/Observation")
		}

		self.ProcessResultReport()
		Upload(&self.Report, baseURL+"/DiagnosticReport")
	}
}

func (self *Procedure) ProcessResultObservations() {
	fhirResultObservations := make([]UploadableObservation, 0)

	for _, value := range self.Values {
		fhirObservation := models.Observation{Reliability: "ok", Status: "final"}
		fhirObservation.Name = self.ConvertCodingToFHIR()
		fhirObservation.Name.Text = self.Description
		self.HandleValue(&fhirObservation, value)
		fhirResultObservations = append(fhirResultObservations, UploadableObservation{FhirObservation: fhirObservation})
	}

	self.ResultObservations = fhirResultObservations
}

func (self *Procedure) ProcessResultReport() {
	self.Report.FhirDiagnosticReport.Result = make([]models.Reference, 0)
	for _, observation := range self.ResultObservations {
		self.Report.FhirDiagnosticReport.Result = append(self.Report.FhirDiagnosticReport.Result, models.Reference{Reference: observation.ServerURL})
	}
}

func (self *Procedure) ToFhirModel() models.Procedure {
	fhirProcedure := models.Procedure{}
	fhirProcedure.Type = self.ConvertCodingToFHIR()
	fhirProcedure.Type.Text = self.Description
	fhirProcedure.Encounter = models.Reference{Reference: self.Patient.MatchingEncounter(self.Entry).ServerURL}
	fhirProcedure.Notes = self.Description
	fhirProcedure.Date = self.AsFHIRPeriod()
	if self.Report.ServerURL != "" {
		fhirProcedure.Report = append(fhirProcedure.Report, models.Reference{Reference: self.Report.ServerURL})
	}

	fhirProcedure.Subject = models.Reference{Reference: self.Patient.ServerURL}
	return fhirProcedure
}

func (self *Procedure) ToJSON() []byte {
	fhirProcedure := self.ToFhirModel()
	json, _ := json.Marshal(fhirProcedure)
	return json
}
