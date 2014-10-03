package hdsfhir

import (
	"encoding/json"
	"time"
)

type Patient struct {
	FirstName     string `json:"first"`
	LastName      string `json:"last"`
	UnixBirthTime int64  `json:"birthdate"`
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
	}
	json, _ := json.Marshal(f)
	return json
}
