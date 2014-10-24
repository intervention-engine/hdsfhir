package hdsfhir

import (
	"encoding/json"
	"time"
)

type Patient struct {
	FirstName     string       `json:"first"`
	LastName      string       `json:"last"`
	UnixBirthTime int64        `json:"birthdate"`
	Gender        string       `json:"gender"`
	Encounters    []*Encounter `json:"encounters"`
	Conditions    []*Condition `json:"conditions"`
	ServerURL     string       `json:"-"`
}

func (p *Patient) SetServerURL(url string) {
	p.ServerURL = url
}

func (p *Patient) BirthTime() time.Time {
	return time.Unix(p.UnixBirthTime, 0)
}

func (p *Patient) PostToFHIRServer(baseURL string) {
	Upload(p, baseURL+"/Patient")
	for _, encounter := range p.Encounters {
		encounter.Patient = p
		Upload(encounter, baseURL+"/Encounter")
	}
	for _, condition := range p.Conditions {
		condition.Patient = p
		Upload(condition, baseURL+"/Condition")
	}
}

func (p *Patient) ToJSON() []byte {
	f := map[string]interface{}{
		"name": []FHIRName{
			FHIRName{FirstName: []string{p.FirstName}, LastName: []string{p.LastName}},
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
