package enrollment

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock del servicio (mantenemos el que ya tienes)
type svcMock struct{ mock.Mock }

func (m *svcMock) Create(sID, cID string, amt float64) (*Enrollment, error) {
	args := m.Called(sID, cID, amt)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Enrollment), args.Error(1)
}
func (m *svcMock) GetAll() ([]Enrollment, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]Enrollment), args.Error(1)
}
func (m *svcMock) Get(id string) (*Enrollment, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Enrollment), args.Error(1)
}
func (m *svcMock) Put(id, sID, cID string, amt float64) error {
	return m.Called(id, sID, cID, amt).Error(0)
}
func (m *svcMock) Patch(id string, amt *float64) error { return m.Called(id, amt).Error(0) }
func (m *svcMock) Delete(id string) error              { return m.Called(id).Error(0) }

func TestController_Endpoints(t *testing.T) {
	m := new(svcMock)

	// --- FLUJOS DE CREATE ---
	t.Run("POST Create - Success", func(t *testing.T) {
		m.On("Create", "s1", "c1", 150.0).Return(&Enrollment{ID: "new-id"}, nil).Once()
		handler := makeCreateHandler(m)
		body, _ := json.Marshal(map[string]interface{}{
			"student_id": "s1", "course_id": "c1", "total_amount": 150.0,
		})
		req := httptest.NewRequest("POST", "/enrollments", bytes.NewBuffer(body))
		rr := httptest.NewRecorder()
		handler(rr, req)
		assert.Equal(t, http.StatusCreated, rr.Code)
	})

	t.Run("POST Create - Invalid JSON", func(t *testing.T) {
		handler := makeCreateHandler(m)
		req := httptest.NewRequest("POST", "/enrollments", bytes.NewBuffer([]byte(`{invalid json`)))
		rr := httptest.NewRecorder()
		handler(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	// --- FLUJOS DE GET ALL ---
	t.Run("GET GetAll - Success", func(t *testing.T) {
		m.On("GetAll").Return([]Enrollment{{ID: "1"}, {ID: "2"}}, nil).Once()
		handler := makeGetAllHandler(m)
		req := httptest.NewRequest("GET", "/enrollments", nil)
		rr := httptest.NewRecorder()
		handler(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
		var resp []Enrollment
		json.NewDecoder(rr.Body).Decode(&resp)
		assert.Len(t, resp, 2)
	})

	t.Run("GET GetAll - Server Error", func(t *testing.T) {
		m.On("GetAll").Return(nil, errors.New("db error")).Once()
		handler := makeGetAllHandler(m)
		req := httptest.NewRequest("GET", "/enrollments", nil)
		rr := httptest.NewRecorder()
		handler(rr, req)
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})

	// --- FLUJOS DE GET BY ID ---
	t.Run("GET GetByID - Success", func(t *testing.T) {
		m.On("Get", "123").Return(&Enrollment{ID: "123"}, nil).Once()
		handler := makeGetHandler(m)
		req := httptest.NewRequest("GET", "/enrollments/123", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "123"})
		rr := httptest.NewRecorder()
		handler(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("GET GetByID - Not Found", func(t *testing.T) {
		m.On("Get", "404").Return(nil, errors.New("not found")).Once()
		handler := makeGetHandler(m)
		req := httptest.NewRequest("GET", "/enrollments/404", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "404"})
		rr := httptest.NewRecorder()
		handler(rr, req)
		assert.Equal(t, http.StatusNotFound, rr.Code)
	})

	// --- FLUJOS DE PATCH ---
	t.Run("PATCH Amount - Success", func(t *testing.T) {
		amt := 100.0
		m.On("Patch", "123", &amt).Return(nil).Once()
		handler := makePatchHandler(m)
		body, _ := json.Marshal(map[string]interface{}{"total_amount": 100.0})
		req := httptest.NewRequest("PATCH", "/enrollments/123", bytes.NewBuffer(body))
		req = mux.SetURLVars(req, map[string]string{"id": "123"})
		rr := httptest.NewRecorder()
		handler(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	// --- FLUJOS DE DELETE ---
	t.Run("DELETE - Success", func(t *testing.T) {
		m.On("Delete", "123").Return(nil).Once()
		handler := makeDeleteHandler(m)
		req := httptest.NewRequest("DELETE", "/enrollments/123", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "123"})
		rr := httptest.NewRecorder()
		handler(rr, req)
		assert.Equal(t, http.StatusNoContent, rr.Code)
	})

	t.Run("DELETE - Error", func(t *testing.T) {
		m.On("Delete", "500").Return(errors.New("db error")).Once()
		handler := makeDeleteHandler(m)
		req := httptest.NewRequest("DELETE", "/enrollments/500", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "500"})
		rr := httptest.NewRecorder()
		handler(rr, req)
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}
