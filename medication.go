package hdsfhir

import (
	"encoding/json"
	"fmt"
	"github.com/intervention-engine/fhir/models"
	"net/http"
)

type Medication struct {
	Entry
	Route   map[string]string
	BaseUrl string
}

type FHIRMedicationWrapper struct {
	Medication models.Medication
	ServerURL  string
}

func (f FHIRMedicationWrapper) ToJSON() []byte {
	json, _ := json.Marshal(f.Medication)
	return json
}

func (f FHIRMedicationWrapper) SetServerURL(url string) {
	f.ServerURL = url
}

func NewMedicationWrapper(med Medication) *FHIRMedicationWrapper {
	fhirMedication := models.Medication{}
	fhirMedication.Code = med.ConvertCodingToFHIR()
	fhirMedication.Name = med.Description
	return &FHIRMedicationWrapper{Medication: fhirMedication}
}

func (m *Medication) ToFhirModel() models.MedicationStatement {
	fhirMedication := models.MedicationStatement{}
	fhirMedication.WhenGiven = m.AsFHIRPeriod()
	fhirMedication.Patient = models.Reference{Reference: m.Patient.ServerURL}
	medUrl := m.FindOrCreateFHIRMed()
	fhirMedication.Medication = models.Reference{Reference: medUrl}
	return fhirMedication
}

func (m *Medication) ToJSON() []byte {
	json, _ := json.Marshal(m.ToFhirModel())
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
	medicationBundle := &models.MedicationBundle{}
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
