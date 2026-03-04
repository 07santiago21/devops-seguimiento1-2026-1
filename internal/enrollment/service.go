package enrollment

import "errors"

type Service interface {
	Create(studentID, courseID string, amount float64) (*Enrollment, error)
	GetAll() ([]Enrollment, error)
	Get(id string) (*Enrollment, error)
	Delete(id string) error
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
	return s.repo.Delete(id)
}
