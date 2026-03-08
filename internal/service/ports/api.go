package ports

import (
	"context"

	"github.com/flexer2006/case-person-enrichment-go/internal/service/domain"

	"github.com/google/uuid"
)

type API interface {
	Person() interface {
		GetByID(ctx context.Context, id uuid.UUID) (*domain.Person, error)
		GetPersons(ctx context.Context, filter map[string]any, offset, limit int) ([]*domain.Person, int, error)
		CreatePerson(ctx context.Context, person *domain.Person) error
		UpdatePerson(ctx context.Context, person *domain.Person) error
		DeletePerson(ctx context.Context, id uuid.UUID) error
		EnrichPerson(ctx context.Context, id uuid.UUID) (*domain.Person, error)
	}

	Age() interface {
		GetAgeByName(ctx context.Context, name string) (int, float64, error)
	}

	Gender() interface {
		GetGenderByName(ctx context.Context, name string) (string, float64, error)
	}

	Nationality() interface {
		GetNationalityByName(ctx context.Context, name string) (string, float64, error)
	}
}
