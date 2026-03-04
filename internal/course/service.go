package course

import "errors"

type Service interface {
	Create(name, code string, credits int32) (*Course, error)
	GetAll() ([]Course, error)
	Get(id string) (*Course, error)
	Delete(id string) error
	Patch(id string, name *string, code *string, credits *int32) error
	Put(id string, name string, code string, credits int32) (*Course, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(name, code string, credits int32) (*Course, error) {
	if name == "" {
		return nil, errors.New("name is required")
	}
	if code == "" {
		return nil, errors.New("code is required")
	}
	if credits <= 0 {
		return nil, errors.New("credits must be greater than zero")
	}

	course := &Course{
		Name:    name,
		Credits: credits,
	}

	if err := s.repo.Create(course); err != nil {
		return nil, err
	}
	return course, nil
}

func (s *service) GetAll() ([]Course, error) {
	return s.repo.GetAll()
}

func (s *service) Get(id string) (*Course, error) {
	return s.repo.Get(id)
}

func (s *service) Delete(id string) error {
	return s.repo.Delete(id)
}

func (s *service) Patch(id string, name *string, code *string, credits *int32) error {
	return s.repo.Patch(id, name, code, credits)
}

func (s *service) Put(id string, name string, code string, credits int32) (*Course, error) {
	if name == "" {
		return nil, errors.New("name is required")
	}
	if code == "" {
		return nil, errors.New("code is required")
	}

	if err := s.repo.Put(id, name, code, credits); err != nil {
		return nil, err
	}
	return s.Get(id)
}
