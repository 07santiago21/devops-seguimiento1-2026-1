package course

import (
	"errors"
	"testing"
)

type mockCourseRepository struct {
	createFunc func(*Course) error
	getAllFunc func() ([]Course, error)
	getFunc    func(id string) (*Course, error)
	deleteFunc func(id string) error
	patchFunc  func(id string, name *string, description *string, credits *int32, capacity *int32) error
	putFunc    func(id, name, description string, credits, capacity int32) error
}

func (m *mockCourseRepository) Create(c *Course) error {
	if m.createFunc == nil {
		return errors.New("not implemented")
	}
	return m.createFunc(c)
}

func (m *mockCourseRepository) GetAll() ([]Course, error) {
	if m.getAllFunc == nil {
		return nil, errors.New("not implemented")
	}
	return m.getAllFunc()
}

func (m *mockCourseRepository) Get(id string) (*Course, error) {
	if m.getFunc == nil {
		return nil, errors.New("not implemented")
	}
	return m.getFunc(id)
}

func (m *mockCourseRepository) Delete(id string) error {
	if m.deleteFunc == nil {
		return errors.New("not implemented")
	}
	return m.deleteFunc(id)
}

func (m *mockCourseRepository) Patch(id string, name *string, description *string, credits *int32, capacity *int32) error {
	if m.patchFunc == nil {
		return errors.New("not implemented")
	}
	return m.patchFunc(id, name, description, credits, capacity)
}

func (m *mockCourseRepository) Put(i, n, d string, cr, c int32) error {
	if m.putFunc == nil {
		return errors.New("not implemented")
	}
	return m.putFunc(i, n, d, cr, c)
}

// ========== CREATE TESTS ==========
func TestServiceCourseCreate_Success(t *testing.T) {
	mockRepo := &mockCourseRepository{
		createFunc: func(c *Course) error { return nil },
	}
	svc := NewService(mockRepo)
	course, err := svc.Create("Math", "Mathematics course", 4, 30)
	if err != nil {
		t.Fatalf("expected no error got %v", err)
	}
	if course.Name != "Math" {
		t.Errorf("expected Name 'Math' got '%s'", course.Name)
	}
}

func TestServiceCourseCreate_ValidationErrors(t *testing.T) {
	svc := NewService(&mockCourseRepository{})

	t.Run("Empty Name", func(t *testing.T) {
		_, err := svc.Create("", "Desc", 4, 30)
		if err == nil || err.Error() != "name is required" {
			t.Error("expected 'name is required'")
		}
	})

	t.Run("Invalid Credits", func(t *testing.T) {
		_, err := svc.Create("Math", "Desc", 0, 30)
		if err == nil || err.Error() != "credits must be greater than zero" {
			t.Error("expected credit error")
		}
	})
}

// ========== GET/GETALL TESTS ==========
func TestServiceCourseGetAll(t *testing.T) {
	mockRepo := &mockCourseRepository{
		getAllFunc: func() ([]Course, error) {
			return []Course{{ID: "1", Name: "Math"}}, nil
		},
	}
	svc := NewService(mockRepo)
	courses, _ := svc.GetAll()
	if len(courses) != 1 {
		t.Errorf("expected 1 course got %d", len(courses))
	}
}

// ========== DELETE TESTS ==========
func TestServiceCourseDelete(t *testing.T) {
	mockRepo := &mockCourseRepository{
		deleteFunc: func(id string) error { return nil },
	}
	svc := NewService(mockRepo)
	err := svc.Delete("1")
	if err != nil {
		t.Error("expected no error")
	}
}

// ========== PUT TESTS ==========
func TestServiceCoursePut_Success(t *testing.T) {
	var saved *Course
	mockRepo := &mockCourseRepository{
		putFunc: func(id, name, description string, credits, capacity int32) error {
			saved = &Course{ID: id, Name: name, Description: description, Credits: credits, Capacity: capacity}
			return nil
		},
		getFunc: func(id string) (*Course, error) { return saved, nil },
	}
	svc := NewService(mockRepo)

	course, err := svc.Put("1", "Physics", "Physics course", 3, 25)
	if err != nil {
		t.Fatalf("expected no error got %v", err)
	}
	if course.Name != "Physics" {
		t.Errorf("expected Name 'Physics' got '%s'", course.Name)
	}
}

func TestServiceCoursePut_Errors(t *testing.T) {
	svc := NewService(&mockCourseRepository{})
	t.Run("Empty Name", func(t *testing.T) {
		_, err := svc.Put("1", "", "Desc", 3, 25)
		if err == nil {
			t.Error("expected error")
		}
	})
}

func TestNewService(t *testing.T) {
	svc := NewService(&mockCourseRepository{})
	if svc == nil {
		t.Fatal("expected service")
	}
}
