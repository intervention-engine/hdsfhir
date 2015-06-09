package models

import (
	"math/rand"
	"strconv"
	"time"

	fhir "github.com/intervention-engine/fhir/models"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

type TemporallyIdentified struct {
	tempID string `json:"-"`
}

// Not even remotely thread-safe
func (t *TemporallyIdentified) GetTempID() string {
	if t.tempID == "" {
		t.tempID = strconv.FormatInt(rand.Int63(), 10)
	}
	return t.tempID
}

func (t *TemporallyIdentified) FHIRReference() *fhir.Reference {
	return &fhir.Reference{Reference: "cid:" + t.GetTempID()}
}
