package student

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockSvc struct {
	OnCreate func(n, l string, a int32) (*Student, error)
}

func (m *MockSvc) Create(n, l string, a int32) (*Student, error) {
	return m.OnCreate(n, l, a)
}

func TestController_Create(t *testing.T) {
	ast := assert.New(t)

	// Caso éxito
	svc := &MockSvc{OnCreate: func(n, l string, a int32) (*Student, error) {
		return &Student{ID: "1", Name: n, LastName: l, Age: a}, nil
	}}
	handler := makeCreateHandler(svc)

	body, _ := json.Marshal(createRequest{Name: "Ana", LastName: "G", Age: 25})
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	handler(rr, req)

	ast.Equal(http.StatusCreated, rr.Code)
}

func TestController_Errors(t *testing.T) {
	ast := assert.New(t)
	svc := &MockSvc{}
	handler := makeCreateHandler(svc)

	body := []byte(`{"name": "Juan", "last_name": "Perez", "age": "veinte"}`)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	handler(rr, req)
	ast.Equal(http.StatusBadRequest, rr.Code)

}
