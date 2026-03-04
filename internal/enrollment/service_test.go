package enrollment

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

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

	t.Run("Create - Success", func(t *testing.T) {
		// Verificamos que se llame a Create con cualquier puntero de Enrollment
		m.On("Create", mock.AnythingOfType("*enrollment.Enrollment")).Return(nil).Once()

		res, err := s.Create("s1", "c1", 100.0)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, 100.0, res.TotalAmount)
		m.AssertExpectations(t)
	})

	t.Run("Create - Negative Amount Error", func(t *testing.T) {
		_, err := s.Create("s1", "c1", -10.0)
		assert.Error(t, err)
		assert.Equal(t, "amount cannot be negative", err.Error())
	})

	t.Run("GetAll - Success", func(t *testing.T) {
		list := []Enrollment{{ID: "1"}, {ID: "2"}}
		m.On("GetAll").Return(list, nil).Once()

		res, err := s.GetAll()
		assert.NoError(t, err)
		assert.Len(t, res, 2)
	})

	t.Run("Get - Success", func(t *testing.T) {
		m.On("Get", "123").Return(&Enrollment{ID: "123"}, nil).Once()

		res, err := s.Get("123")
		assert.NoError(t, err)
		assert.Equal(t, "123", res.ID)
	})

	t.Run("Patch - Success", func(t *testing.T) {
		amt := 200.0
		m.On("Patch", "123", &amt).Return(nil).Once()

		err := s.Patch("123", &amt)
		assert.NoError(t, err)
	})

	t.Run("Patch - Negative Error", func(t *testing.T) {
		amt := -5.0
		err := s.Patch("123", &amt)
		assert.Error(t, err)
		assert.Equal(t, "amount cannot be negative", err.Error())
	})

	t.Run("Delete - Success", func(t *testing.T) {
		m.On("Delete", "123").Return(nil).Once()

		err := s.Delete("123")
		assert.NoError(t, err)
	})
}
