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
		Delete(id string) error
		Patch(id string, Name *string, LastName *string, Age *int32) error
		Put(id string, Name string, LastName string, Age int32) error
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

func (r *repository) Delete(id string) error {

	student := Student{ID: id}
	result := r.db.Delete(&student)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *repository) Patch(id string, Name *string, LastName *string, Age *int32) error {

	values := make(map[string]interface{})

	if Name != nil {
		values["Name"] = *Name
	}

	if LastName != nil {
		values["LastName "] = *LastName
	}

	if Age != nil {
		values["Age"] = *Age
	}

	if result := r.db.Model(&Student{}).Where("id = ?", id).Updates(values); result.Error != nil {
		return result.Error
	}
	return nil

}
func (r *repository) Put(id string, Name string, LastName string, Age int32) error {
	student := Student{
		ID:       id,
		Name:     Name,
		LastName: LastName,
		Age:      Age,
	}

	result := r.db.Model(&Student{}).Where("id = ?", id).Updates(student)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
