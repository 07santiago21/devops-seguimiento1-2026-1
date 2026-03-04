package enrollment

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

	e := &Enrollment{
		StudentID:   studentID,
		CourseID:    courseID,
		TotalAmount: amount,
	}

	if err := s.repo.Create(e); err != nil {
		return nil, err
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
