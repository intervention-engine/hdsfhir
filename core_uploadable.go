package hdsfhir

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/intervention-engine/fhir/models"
)

type Uploadable interface {
	ToJSON() []byte
	SetServerURL(string)
}

func Upload(thing Uploadable, url string) {
	body := bytes.NewReader(thing.ToJSON())
	response, err := http.Post(url, "application/json+fhir", body)
	if err != nil {
		panic("HTTP request failed")
	}
	thing.SetServerURL(response.Header.Get("Location"))
}

type UploadableObservation struct {
	ServerURL       string
	FhirObservation models.Observation
}

func (self *UploadableObservation) SetServerURL(url string) {
	self.ServerURL = url
}

func (self *UploadableObservation) ToJSON() []byte {
	json, _ := json.Marshal(self.FhirObservation)
	return json
}

type UploadableDiagnosticReport struct {
	ServerURL            string
	FhirDiagnosticReport models.DiagnosticReport
}

func (self *UploadableDiagnosticReport) SetServerURL(url string) {
	self.ServerURL = url
}

func (self *UploadableDiagnosticReport) ToJSON() []byte {
	json, _ := json.Marshal(self.FhirDiagnosticReport)
	return json
}
