package student

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Mock para simular el repositorio
type MockRepo struct {
	OnCreate func(s *Student) error
}

func (m *MockRepo) Create(s *Student) error {
	return m.OnCreate(s)
}

func TestCreateService(t *testing.T) {
	ast := assert.New(t)

	tests := []struct {
		name        string
		inputName   string
		inputLast   string
		inputAge    int32
		mockErr     error
		expectedErr string
	}{
		{"Éxito", "Juan", "Perez", 20, nil, ""},
		{"Error Nombre", "", "Perez", 20, nil, "name is required"},
		{"Error Apellido", "Juan", "", 20, nil, "last_name is required"},
		{"Error Edad", "Juan", "Perez", 0, nil, "age must be greater than zero"},
		{"Error Repo", "Juan", "Perez", 20, errors.New("db error"), "db error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &MockRepo{OnCreate: func(s *Student) error { return tt.mockErr }}
			srv := NewService(repo)

			res, err := srv.Create(tt.inputName, tt.inputLast, tt.inputAge)

			if tt.expectedErr != "" {
				ast.Error(err)
				ast.Contains(err.Error(), tt.expectedErr)
			} else {
				ast.NoError(err)
				ast.Equal(tt.inputName, res.Name)
			}
		})
	}
}
