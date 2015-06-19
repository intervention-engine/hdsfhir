package hdsfhir

import fhir "github.com/intervention-engine/fhir/models"

type Encounter struct {
	Entry
	DischargeDisposition CodeMap `json:"dischargeDisposition"`
}

func (e *Encounter) FHIRModels() []interface{} {
	fhirEncounter := &fhir.Encounter{Id: e.GetTempID()}
	cc := e.Codes.FHIRCodeableConcept(e.Description)
	fhirEncounter.Type = []fhir.CodeableConcept{*cc}
	fhirEncounter.Period = e.GetFHIRPeriod()
	fhirEncounter.Patient = e.Patient.FHIRReference()

	return []interface{}{fhirEncounter}
}
