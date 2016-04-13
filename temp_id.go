package hdsfhir

import (
	"sync"

	fhir "github.com/intervention-engine/fhir/models"
	"github.com/satori/go.uuid"
)

type TemporallyIdentified struct {
	once   sync.Once `json:"-"`
	tempID string    `json:"-"`
}

func (t *TemporallyIdentified) GetTempID() string {
	t.once.Do(func() {
		t.tempID = uuid.NewV4().String()
	})
	return t.tempID
}

func (t *TemporallyIdentified) FHIRReference() *fhir.Reference {
	return &fhir.Reference{Reference: "urn:uuid:" + t.GetTempID()}
}
