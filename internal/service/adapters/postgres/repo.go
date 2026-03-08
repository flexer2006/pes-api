package postgres

import (
	"context"

	"github.com/flexer2006/case-person-enrichment-go/internal/service/domain"
	"github.com/flexer2006/case-person-enrichment-go/internal/service/ports"
	"github.com/flexer2006/case-person-enrichment-go/pkg/database/postgres"

	"github.com/google/uuid"
)

var _ ports.Repositories = (*Repositories)(nil)

type Repositories struct {
	personRepo *Repository
}

func NewRepositories(db postgres.Provider) *Repositories {
	return &Repositories{
		personRepo: NewRepository(db),
	}
}

func (r *Repositories) Person() interface {
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Person, error)
	GetPersons(ctx context.Context, filter map[string]any, offset, limit int) ([]*domain.Person, int, error)
	CreatePerson(ctx context.Context, person *domain.Person) error
	UpdatePerson(ctx context.Context, person *domain.Person) error
	DeletePerson(ctx context.Context, id uuid.UUID) error
	ExistsByID(ctx context.Context, id uuid.UUID) (bool, error)
} {
	return r.personRepo
}
