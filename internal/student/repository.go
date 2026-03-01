package student

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	Repository interface {
		Create(student *Student) error
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
