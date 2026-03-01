package student

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	Repository interface {
		Create(student *Student) error
		GetAll() ([]Student, error)
		Get(id string) (*Student, error)
	}

	repository struct {
		db *gorm.DB
	}
)

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}

}
func (r *repository) Create(student *Student) error {
	student.ID = uuid.New().String()
	return r.db.Create(student).Error
}

func (r *repository) GetAll() ([]Student, error) {
	var students []Student
	result := r.db.Model(&students).Order("created_at desc").Find(&students)

	if result.Error != nil {
		return nil, result.Error
	}

	return students, nil
}

func (r *repository) Get(id string) (*Student, error) {

	student := Student{ID: id}
	result := r.db.First(&student)

	if result.Error != nil {
		return nil, result.Error
	}

	return &student, nil

}
