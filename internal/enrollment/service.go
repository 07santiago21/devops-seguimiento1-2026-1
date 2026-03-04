package enrollment

import "errors"

type Service interface {
	Create(studentID, courseID string, amount float64) (*Enrollment, error)
	GetAll() ([]Enrollment, error)
	Get(id string) (*Enrollment, error)
	Delete(id string) error
	Put(id string, studentID, courseID string, amount float64) error
	Patch(id string, amount *float64) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(studentID, courseID string, amount float64) (*Enrollment, error) {
	if studentID == "" || courseID == "" {
		return nil, errors.New("student_id and course_id are required")
	}
	if amount < 0 {
		return nil, errors.New("amount cannot be negative")
	}

	e := &Enrollment{
		StudentID:   studentID,
		CourseID:    courseID,
		TotalAmount: amount,
	}

	if err := s.repo.Create(e); err != nil {
		return nil, errors.New("invalid student_id or course_id: entity not found")
	}
	return e, nil
}

func (s *service) GetAll() ([]Enrollment, error) {
	return s.repo.GetAll()
}

func (s *service) Get(id string) (*Enrollment, error) {
	return s.repo.Get(id)
}

func (s *service) Delete(id string) error {
	if id == "" {
		return errors.New("enrollment ID is required")
	}
	return s.repo.Delete(id)
}
func (s *service) Put(id string, studentID, courseID string, amount float64) error {
	if studentID == "" || courseID == "" || amount < 0 {
		return errors.New("invalid data for full update")
	}
	return s.repo.Put(id, studentID, courseID, amount)
}

func (s *service) Patch(id string, amount *float64) error {
	if amount != nil && *amount < 0 {
		return errors.New("amount cannot be negative")
	}
	return s.repo.Patch(id, amount)
}
