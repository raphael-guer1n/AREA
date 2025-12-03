package domain

import "time"

type User struct {
	ID        int
	Email     string
	FirstName string
	LastName  string
	CreatedAt time.Time
}

type UserRepository interface {
	List() ([]User, error)
	Create(email, firstName, lastName string) (User, error)
}
