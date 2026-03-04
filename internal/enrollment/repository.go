package enrollment

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	Create(enrollment *Enrollment) error
	GetAll() ([]Enrollment, error)
	Get(id string) (*Enrollment, error)
	Delete(id string) error
}

type repo struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repo{db: db}
}

func (r *repo) Create(e *Enrollment) error {
	e.ID = uuid.New().String()
	return r.db.Create(e).Error
}

func (r *repo) GetAll() ([]Enrollment, error) {
	var enrollments []Enrollment
	err := r.db.Preload("Student").Preload("Course").Order("created_at desc").Find(&enrollments).Error
	return enrollments, err
}

func (r *repo) Get(id string) (*Enrollment, error) {
	var e Enrollment
	if err := r.db.Preload("Student").Preload("Course").Where("id = ?", id).First(&e).Error; err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *repo) Delete(id string) error {
	result := r.db.Delete(&Enrollment{}, "id = ?", id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
