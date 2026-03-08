package enrichment

import (
	"context"

	api "github.com/flexer2006/case-person-enrichment-go/internal/service/adapters/enrichment/services"
	"github.com/flexer2006/case-person-enrichment-go/internal/service/domain"
	apiports "github.com/flexer2006/case-person-enrichment-go/internal/service/ports"

	"github.com/google/uuid"
)

var _ apiports.API = (*Enrichment)(nil)

type Enrichment struct {
	impl *api.API
}

func NewEnrichment(apiImpl *api.API) *Enrichment {
	return &Enrichment{impl: apiImpl}
}

func NewDefaultEnrichment() *Enrichment {
	return &Enrichment{impl: api.NewDefaultAPI()}
}

func (e *Enrichment) Age() interface {
	GetAgeByName(ctx context.Context, name string) (int, float64, error)
} {
	return e.impl.Age()
}

func (e *Enrichment) Gender() interface {
	GetGenderByName(ctx context.Context, name string) (string, float64, error)
} {
	return e.impl.Gender()
}

func (e *Enrichment) Nationality() interface {
	GetNationalityByName(ctx context.Context, name string) (string, float64, error)
} {
	return e.impl.Nationality()
}

func (e *Enrichment) Person() interface {
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Person, error)
	GetPersons(ctx context.Context, filter map[string]any, offset, limit int) ([]*domain.Person, int, error)
	CreatePerson(ctx context.Context, person *domain.Person) error
	UpdatePerson(ctx context.Context, person *domain.Person) error
	DeletePerson(ctx context.Context, id uuid.UUID) error
	EnrichPerson(ctx context.Context, id uuid.UUID) (*domain.Person, error)
} {
	return e.impl.Person()
}
