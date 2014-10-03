package hdsfhir

import (
	"time"
)

type Patient struct {
	FirstName     string `json:"first"`
	LastName      string `json:"last"`
	UnixBirthTime int64  `json:"birthdate"`
}

func (p *Patient) BirthTime() Time {
	return time.Unix(p.UnixBirthTime, 0)
}
