package enrollment

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// repoMock es el mock para el repositorio de enrollment
type repoMock struct{ mock.Mock }

func (m *repoMock) Create(e *Enrollment) error { return m.Called(e).Error(0) }
func (m *repoMock) GetAll() ([]Enrollment, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]Enrollment), args.Error(1)
}
func (m *repoMock) Get(id string) (*Enrollment, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Enrollment), args.Error(1)
}
func (m *repoMock) Put(id, sID, cID string, amt float64) error {
	return m.Called(id, sID, cID, amt).Error(0)
}
func (m *repoMock) Patch(id string, amt *float64) error {
	return m.Called(id, amt).Error(0)
}
func (m *repoMock) Delete(id string) error {
	return m.Called(id).Error(0)
}

func TestService_Logic(t *testing.T) {
	m := new(repoMock)
	s := NewService(m)

	// --- PRUEBAS DE CREATE ---

	t.Run("Create - Success", func(t *testing.T) {
		m.On("Create", mock.AnythingOfType("*enrollment.Enrollment")).Return(nil).Once()
		res, err := s.Create("student-1", "course-1", 100.0)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, 100.0, res.TotalAmount)
	})

	t.Run("Create - Empty IDs Error", func(t *testing.T) {
		res, err := s.Create("", "course-1", 100.0)
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, "student_id and course_id are required", err.Error())
	})

	t.Run("Create - Negative Amount Error", func(t *testing.T) {
		res, err := s.Create("s1", "c1", -50.0)
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, "amount cannot be negative", err.Error())
	})

	t.Run("Create - Repo Error (Not Found)", func(t *testing.T) {
		m.On("Create", mock.Anything).Return(errors.New("db error")).Once()
		res, err := s.Create("invalid-s", "invalid-c", 100.0)
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "entity not found")
	})

	// --- PRUEBAS DE GET & GETALL ---

	t.Run("GetAll - Success", func(t *testing.T) {
		list := []Enrollment{{ID: "1"}, {ID: "2"}}
		m.On("GetAll").Return(list, nil).Once()
		res, err := s.GetAll()
		assert.NoError(t, err)
		assert.Len(t, res, 2)
	})

	t.Run("Get - Success", func(t *testing.T) {
		m.On("Get", "uuid-1").Return(&Enrollment{ID: "uuid-1"}, nil).Once()
		res, err := s.Get("uuid-1")
		assert.NoError(t, err)
		assert.Equal(t, "uuid-1", res.ID)
	})

	// --- PRUEBAS DE DELETE ---

	t.Run("Delete - Success", func(t *testing.T) {
		m.On("Delete", "uuid-1").Return(nil).Once()
		err := s.Delete("uuid-1")
		assert.NoError(t, err)
	})

	t.Run("Delete - Empty ID Error", func(t *testing.T) {
		err := s.Delete("")
		assert.Error(t, err)
		assert.Equal(t, "enrollment ID is required", err.Error())
	})

	// --- PRUEBAS DE PUT (FULL UPDATE) ---

	t.Run("Put - Success", func(t *testing.T) {
		m.On("Put", "id1", "s1", "c1", 300.0).Return(nil).Once()
		err := s.Put("id1", "s1", "c1", 300.0)
		assert.NoError(t, err)
	})

	t.Run("Put - Invalid Data Error", func(t *testing.T) {
		// Probamos con monto negativo
		err := s.Put("id1", "s1", "c1", -1.0)
		assert.Error(t, err)
		assert.Equal(t, "invalid data for full update", err.Error())
	})

	// --- PRUEBAS DE PATCH (PARTIAL UPDATE) ---

	t.Run("Patch - Success", func(t *testing.T) {
		amt := 450.0
		m.On("Patch", "id1", &amt).Return(nil).Once()
		err := s.Patch("id1", &amt)
		assert.NoError(t, err)
	})

	t.Run("Patch - Negative Amount Error", func(t *testing.T) {
		amt := -10.0
		err := s.Patch("id1", &amt)
		assert.Error(t, err)
		assert.Equal(t, "amount cannot be negative", err.Error())
	})

	t.Run("Patch - Nil Amount Success", func(t *testing.T) {
		m.On("Patch", "id1", (*float64)(nil)).Return(nil).Once()
		err := s.Patch("id1", nil)
		assert.NoError(t, err)
	})

	m.AssertExpectations(t)
}
