package hdsfhir

import fhir "github.com/intervention-engine/fhir/models"

type Medication struct {
	Entry
	Route map[string]string `json:"route"`
}

func (m *Medication) FHIRModels() []interface{} {
	_, isImmunization := m.Codes["CVX"]
	if isImmunization {
		fhirImmunization := &fhir.Immunization{Id: m.GetTempID()}
		fhirImmunization.Date = m.StartTime.FHIRDateTime()
		fhirImmunization.VaccineType = m.Codes.FHIRCodeableConcept(m.Description)
		fhirImmunization.Patient = m.Patient.FHIRReference()
		// Ignoring Route

		return []interface{}{fhirImmunization}
	} else {
		internalMedID := &TemporallyIdentified{}
		fhirMedication := &fhir.Medication{Id: internalMedID.GetTempID()}
		fhirMedication.Code = m.Codes.FHIRCodeableConcept(m.Description)
		fhirMedication.Name = m.Description

		fhirMedicationStatement := &fhir.MedicationStatement{Id: m.GetTempID()}
		fhirMedicationStatement.EffectivePeriod = m.GetFHIRPeriod()
		fhirMedicationStatement.Patient = m.Patient.FHIRReference()
		fhirMedicationStatement.MedicationReference = internalMedID.FHIRReference()
		// Ignoring route

		return []interface{}{fhirMedication, fhirMedicationStatement}
	}
}
