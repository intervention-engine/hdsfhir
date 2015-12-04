package hdsfhir

import fhir "github.com/intervention-engine/fhir/models"

type VitalSign struct {
	Entry
	Description    string        `json:"description"`
	Interpretation CodeMap       `json:"interpretation"`
	Values         []ResultValue `json:"values"`
}

func (v *VitalSign) FHIRModels() []interface{} {
	var fhirObservation *fhir.Observation
	switch {
	default:
		fhirObservation = &fhir.Observation{}
	case len(v.Values) == 1:
		fhirObservation = v.Values[0].FHIRModels()[0].(*fhir.Observation)
	case len(v.Values) > 1:
		panic("FHIR Observations cannot have more than one value")
	}
	fhirObservation.Code = v.Codes.FHIRCodeableConcept(v.Description)
	fhirObservation.Encounter = v.Patient.MatchingEncounterReference(v.Entry)
	fhirObservation.EffectivePeriod = v.GetFHIRPeriod()
	if len(v.Interpretation) > 0 {
		fhirObservation.Interpretation = v.Interpretation.FHIRCodeableConcept("")
	}
	fhirObservation.Subject = v.Patient.FHIRReference()

	return []interface{}{fhirObservation}
}
