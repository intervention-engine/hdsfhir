package hdsfhir

import fhir "github.com/intervention-engine/fhir/models"

type Immunization struct {
	Entry
	SeriesNumber *uint32 `json:"seriesNumber"`
}

func (i *Immunization) FHIRModels() []interface{} {
	fhirImmunization := &fhir.Immunization{Id: i.GetTempID()}
	fhirImmunization.Status = "completed"
	fhirImmunization.Date = i.Time.FHIRDateTime()
	fhirImmunization.VaccineCode = i.Codes.FHIRCodeableConcept(i.Description)
	fhirImmunization.Patient = i.Patient.FHIRReference()
	if i.NegationInd {
		t := true
		fhirImmunization.WasNotGiven = &t
	}
	if len(i.NegationReason) > 0 {
		cc := i.NegationReason.FHIRCodeableConcept("")
		fhirImmunization.Explanation = &fhir.ImmunizationExplanationComponent{
			ReasonNotGiven: []fhir.CodeableConcept{*cc},
		}
	}
	if i.SeriesNumber != nil {
		fhirImmunization.VaccinationProtocol = []fhir.ImmunizationVaccinationProtocolComponent{
			{DoseSequence: i.SeriesNumber},
		}
	}

	// Ignoring dosage, route, etc.

	return []interface{}{fhirImmunization}
}
