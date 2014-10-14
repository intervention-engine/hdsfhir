package hdsfhir

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

type Patient struct {
	FirstName     string      `json:"first"`
	LastName      string      `json:"last"`
	UnixBirthTime int64       `json:"birthdate"`
	Gender        string      `json:"gender"`
	Encounters    []Entry     `json:"encounters"`
	Conditions    []Condition `json:"conditions"`
	ServerURL     string
}

func (p *Patient) BirthTime() time.Time {
	return time.Unix(p.UnixBirthTime, 0)
}

func (p *Patient) ToJSON() []byte {
	f := map[string]interface{}{
		"name": map[string][]string{
			"family": []string{
				p.LastName,
			},
			"given": []string{
				p.FirstName,
			},
		},
		"gender": map[string]interface{}{
			"coding": []FHIRCoding{
				FHIRCoding{System: "http://hl7.org/fhir/v3/AdministrativeGender", Code: p.Gender},
			},
		},
		"birthDate": p.BirthTime().Format("2006-01-02"),
	}
	json, _ := json.Marshal(f)
	return json
}

func (p *Patient) Upload(url string) {
	body := bytes.NewReader(p.ToJSON())
	response, err := http.Post(url, "application/json+fhir", body)
	defer response.Body.Close()
	if err != nil {
		panic("HTTP request failed")
	}
	p.ServerURL = response.Header.Get("Location")
}
