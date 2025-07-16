package people

import (
	"context"

	"github.com/olegdayo/omniconv"
	"github.com/vladlim/auth-service-practice/auth/internal/repository/models"
)

type Repository interface {
	GetPeople(ctx context.Context) ([]models.Person, error)
}

type Provider struct {
	repository Repository
}

func New(repository Repository) Provider {
	return Provider{
		repository: repository,
	}
}

func (p Provider) GetPeople(ctx context.Context) ([]Person, error) {
	people, err := p.repository.GetPeople(ctx)
	if err != nil {
		return nil, err
	}
	return omniconv.ConvertSlice(people, DBPerson2ProviderPerson), nil
}

func (p Provider) GetPerson(ctx context.Context, id int) (Person, error) {
	return Person{}, nil
}

func (p Provider) AddPerson(ctx context.Context, person Person) (Person, error) {
	return Person{}, nil
}

func (p Provider) UpdatePerson(ctx context.Context, person Person) (Person, error) {
	return Person{}, nil
}

func (p Provider) DeletePerson(ctx context.Context, id int) error {
	return nil
}
