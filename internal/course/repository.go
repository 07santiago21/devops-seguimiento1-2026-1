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
		Patch(id string, Name *string, Code *string, Credits, capacity *int32) error
		Put(id string, Name string, Code string, Credits, capacity int32) error
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
	c := Course{ID: id}
	if err := r.db.First(&c).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *repository) Delete(id string) error {
	result := r.db.Delete(&Course{ID: id})
	return result.Error
}

func (r *repository) Patch(id string, name *string, description *string, credits *int32, capacity *int32) error {
	vals := map[string]interface{}{}
	if name != nil {
		vals["name"] = *name
	}
	if description != nil {
		vals["description"] = *description
	}
	if credits != nil {
		vals["credits"] = *credits
	}
	if capacity != nil {
		vals["capacity"] = *capacity
	}
	return r.db.Model(&Course{}).Where("id = ?", id).Updates(vals).Error
}

func (r *repository) Put(id string, name string, description string, credits, capacity int32) error {
	c := Course{
		ID:          id,
		Name:        name,
		Description: description,
		Credits:     credits,
		Capacity:    capacity,
	}
	res := r.db.Model(&Course{}).Where("id = ?", id).Updates(c)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
