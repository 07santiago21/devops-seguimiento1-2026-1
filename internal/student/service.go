package student

import "errors"

type Service interface {
	Create(name, lastName string, age int32) (*Student, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) Create(name, lastName string, age int32) (*Student, error) {

	if name == "" {
		return nil, errors.New("name is required")
	}

	if lastName == "" {
		return nil, errors.New("last_name is required")
	}

	if age <= 0 {
		return nil, errors.New("age must be greater than zero")
	}

	student := &Student{
		Name:     name,
		LastName: lastName,
		Age:      age,
	}

	if err := s.repo.Create(student); err != nil {
		return nil, err
	}

	return student, nil
}
