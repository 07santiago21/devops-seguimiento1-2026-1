package course

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	Repository interface {
		Create(course *Course) error
		GetAll() ([]Course, error)
		Get(id string) (*Course, error)
		Delete(id string) error
		Patch(id string, Name *string, Code *string, Credits *int32) error
		Put(id string, Name string, Code string, Credits int32) error
	}

	repository struct {
		db *gorm.DB
	}
)

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(course *Course) error {
	course.ID = uuid.New().String()
	return r.db.Create(course).Error
}

func (r *repository) GetAll() ([]Course, error) {
	var list []Course
	result := r.db.Order("created_at desc").Find(&list)
	return list, result.Error
}

func (r *repository) Get(id string) (*Course, error) {
	course := Course{ID: id}
	result := r.db.First(&course)
	if result.Error != nil {
		return nil, result.Error
	}
	return &course, nil
}

func (r *repository) Delete(id string) error {
	result := r.db.Delete(&Course{ID: id})
	return result.Error
}

func (r *repository) Patch(id string, Name *string, Code *string, Credits *int32) error {
	values := make(map[string]interface{})
	if Name != nil {
		values["name"] = *Name
	}
	if Code != nil {
		values["code"] = *Code
	}
	if Credits != nil {
		values["credits"] = *Credits
	}

	return r.db.Model(&Course{}).Where("id = ?", id).Updates(values).Error
}

func (r *repository) Put(id string, Name string, Code string, Credits int32) error {
	course := Course{
		ID:      id,
		Name:    Name,
		Credits: Credits,
	}
	result := r.db.Model(&Course{}).Where("id = ?", id).Updates(course)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
