package people

import "github.com/vladlim/auth-service-practice/auth/internal/repository/models"

type Person struct {
	ID       int
	Username string
}

func ProviderPerson2DBPerson(person Person) models.Person {
	return models.Person{
		ID:       person.ID,
		Username: person.Username,
	}
}

func DBPerson2ProviderPerson(person models.Person) Person {
	return Person{
		ID:       person.ID,
		Username: person.Username,
	}
}
