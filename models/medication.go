package models

import (
	"encoding/json"
	"fmt"
	"net/http"

	fhir "github.com/intervention-engine/fhir/models"
)

type Medication struct {
	Entry
	Route   map[string]string
	BaseUrl string
}

type FHIRMedicationWrapper struct {
	Medication fhir.Medication
	ServerURL  string
}

func (f FHIRMedicationWrapper) ToJSON() []byte {
	json, _ := json.Marshal(f.Medication)
	return json
}

func (f *FHIRMedicationWrapper) SetServerURL(url string) {
	f.ServerURL = url
}

func NewMedicationWrapper(med Medication) *FHIRMedicationWrapper {
	fhirMedication := fhir.Medication{}
	fhirMedication.Code = med.Codes.FHIRCodeableConcept(med.Description)
	fhirMedication.Name = med.Description
	return &FHIRMedicationWrapper{Medication: fhirMedication}
}

func (m *Medication) FHIRModel() fhir.MedicationStatement {
	fhirMedication := fhir.MedicationStatement{}
	fhirMedication.WhenGiven = m.GetFHIRPeriod()
	fhirMedication.Patient = &fhir.Reference{Reference: m.Patient.ServerURL}
	medUrl := m.FindOrCreateFHIRMed()
	fhirMedication.Medication = &fhir.Reference{Reference: medUrl}
	return fhirMedication
}

func (m *Medication) ToJSON() []byte {
	json, _ := json.Marshal(m.FHIRModel())
	return json
}

func (m *Medication) FindOrCreateFHIRMed() string {
	var medicationQueryUrl string
	_, rxNormPresent := m.Codes["RxNorm"]
	if rxNormPresent {
		firstRxNormCode := m.Codes["RxNorm"][0]
		medicationQueryUrl = fmt.Sprintf("%s/Medication?code=%s|%s", m.BaseUrl, "http://www.nlm.nih.gov/research/umls/rxnorm/", firstRxNormCode)
	}

	_, ndcPresent := m.Codes["NDC"]
	if ndcPresent {
		firstNDCCode := m.Codes["NDC"]
		medicationQueryUrl = fmt.Sprintf("%s/Medication?code=%s|%s", m.BaseUrl, "http://www.fda.gov/Drugs/InformationOnDrugs", firstNDCCode)
	}

	resp, err := http.Get(medicationQueryUrl)
	defer resp.Body.Close()
	if err != nil {
		panic(err)
	}
	decoder := json.NewDecoder(resp.Body)
	medicationBundle := &fhir.MedicationBundle{}
	err = decoder.Decode(medicationBundle)
	if err != nil {
		panic(err)
	}
	if medicationBundle.TotalResults > 0 {
		return fmt.Sprintf("%s/Medication/%s", m.BaseUrl, medicationBundle.Entry[0].Id)
	} else {
		fhirWrapper := NewMedicationWrapper(*m)
		Upload(fhirWrapper, fmt.Sprintf("%s/Medication", m.BaseUrl))
		return fhirWrapper.ServerURL
	}
}
