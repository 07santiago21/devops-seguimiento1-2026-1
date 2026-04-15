package student

import (
	"errors"
	"testing"
)

type mockRepository struct {
	createFunc func(*Student) error
	getAllFunc func() ([]Student, error)
	getFunc    func(id string) (*Student, error)
	deleteFunc func(id string) error
	patchFunc  func(id string, name *string, lastName *string, age *int32) error
	putFunc    func(id, name, lastName string, age int32) error
}

func (m *mockRepository) Create(s *Student) error {
	if m.createFunc == nil {
		return errors.New("not implemented")
	}
	return m.createFunc(s)
}
func (m *mockRepository) GetAll() ([]Student, error) {
	if m.getAllFunc == nil {
		return nil, errors.New("not implemented")
	}
	return m.getAllFunc()
}
func (m *mockRepository) Get(id string) (*Student, error) {
	if m.getFunc == nil {
		return nil, errors.New("not implemented")
	}
	return m.getFunc(id)
}
func (m *mockRepository) Delete(id string) error {
	if m.deleteFunc == nil {
		return errors.New("not implemented")
	}
	return m.deleteFunc(id)
}
func (m *mockRepository) Patch(id string, n, l *string, a *int32) error {
	if m.patchFunc == nil {
		return errors.New("not implemented")
	}
	return m.patchFunc(id, n, l, a)
}
func (m *mockRepository) Put(i, n, l string, a int32) error {
	if m.putFunc == nil {
		return errors.New("not implemented")
	}
	return m.putFunc(i, n, l, a)
}

func TestServiceCreate(t *testing.T) {
	mockRepo := &mockRepository{createFunc: func(s *Student) error { return nil }}
	svc := NewService(mockRepo)

	_, err := svc.Create("John", "Doe", 25)
	if err != nil {
		t.Fatal(err)
	}

	// validaciones
	if _, err := svc.Create("", "Doe", 25); err == nil {
		t.Error("expected error for empty name")
	}
	if _, err := svc.Create("John", "", 25); err == nil {
		t.Error("expected error for empty lastName")
	}
	if _, err := svc.Create("John", "Doe", 0); err == nil {
		t.Error("expected error for age <=0")
	}
}

func TestServiceGetAllGetDeletePatchPut(t *testing.T) {
	mock := &mockRepository{
		getAllFunc: func() ([]Student, error) { return []Student{{ID: "1"}}, nil },
		getFunc:    func(id string) (*Student, error) { return &Student{ID: id}, nil },
		deleteFunc: func(id string) error { return nil },
		patchFunc:  func(id string, n, l *string, a *int32) error { return nil },
		putFunc:    func(id, n, l string, a int32) error { return nil },
	}
	svc := NewService(mock)

	if _, err := svc.GetAll(); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.Get("1"); err != nil {
		t.Fatal(err)
	}
	if err := svc.Delete("1"); err != nil {
		t.Fatal(err)
	}
	name := "A"
	if err := svc.Patch("1", &name, nil, nil); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.Put("1", "X", "Y", 20); err != nil {
		t.Fatal(err)
	}
}

func TestServiceRepoErrors(t *testing.T) {
	mock := &mockRepository{
		getFunc: func(id string) (*Student, error) { return nil, errors.New("db") },
	}
	svc := NewService(mock)
	if _, err := svc.Get("1"); err == nil {
		t.Error("expected repo error to propagate")
	}
}

func TestServicePut_ValidationErrors(t *testing.T) {
	svc := NewService(&mockRepository{})
	if _, err := svc.Put("1", "", "Smith", 30); err == nil {
		t.Error("expected empty name error")
	}
	if _, err := svc.Put("1", "Jane", "", 30); err == nil {
		t.Error("expected empty lastName error")
	}
	if _, err := svc.Put("1", "Jane", "Smith", -1); err == nil {
		t.Error("expected invalid age error")
	}
}

func TestServiceGetAll_Error(t *testing.T) {
	svc := NewService(&mockRepository{
		getAllFunc: func() ([]Student, error) {
			return nil, errors.New("db down")
		},
	})

	if _, err := svc.GetAll(); err == nil {
		t.Error("expected get all error to propagate")
	}
}

func TestServiceCuervo_MarshalError(t *testing.T) {
	svc := NewService(&mockRepository{})

	// Channels cannot be JSON-marshaled and should fail before any HTTP request is attempted.
	if _, err := svc.Cuervo(make(chan int)); err == nil {
		t.Error("expected marshal error for unsupported payload type")
	}
}
