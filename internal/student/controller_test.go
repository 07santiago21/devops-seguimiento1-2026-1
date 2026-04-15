package student

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

type mockService struct {
	createFunc func(name, lastName string, age int32) (*Student, error)
	cuervoFunc func(payload any) ([]byte, error)
	getAllFunc func() ([]Student, error)
	getFunc    func(id string) (*Student, error)
	deleteFunc func(id string) error
	patchFunc  func(id string, name *string, lastName *string, age *int32) error
	putFunc    func(id, name, lastName string, age int32) (*Student, error)
}

func (m *mockService) Create(name, lastName string, age int32) (*Student, error) {
	if m.createFunc == nil {
		return nil, errors.New("not implemented")
	}
	return m.createFunc(name, lastName, age)
}
func (m *mockService) Cuervo(payload any) ([]byte, error) {
	if m.cuervoFunc == nil {
		return nil, errors.New("not implemented")
	}
	return m.cuervoFunc(payload)
}

func (m *mockService) GetAll() ([]Student, error) {
	if m.getAllFunc == nil {
		return nil, errors.New("not implemented")
	}
	return m.getAllFunc()
}

func (m *mockService) Get(id string) (*Student, error) {
	if m.getFunc == nil {
		return nil, errors.New("not implemented")
	}
	return m.getFunc(id)
}

func (m *mockService) Delete(id string) error {
	if m.deleteFunc == nil {
		return errors.New("not implemented")
	}
	return m.deleteFunc(id)
}

func (m *mockService) Patch(id string, name *string, lastName *string, age *int32) error {
	if m.patchFunc == nil {
		return errors.New("not implemented")
	}
	return m.patchFunc(id, name, lastName, age)
}

func (m *mockService) Put(id, name, lastName string, age int32) (*Student, error) {
	if m.putFunc == nil {
		return nil, errors.New("not implemented")
	}
	return m.putFunc(id, name, lastName, age)
}

func TestMakeCreateHandler(t *testing.T) {
	mock := &mockService{
		createFunc: func(name, lastName string, age int32) (*Student, error) {
			return &Student{ID: "1", Name: name, LastName: lastName, Age: age}, nil
		},
	}

	handler := makeCreateHandler(mock)
	body := CreateRequest{Name: "John", LastName: "Doe", Age: 25}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/students", bytes.NewReader(b))
	w := httptest.NewRecorder()
	handler(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected %d got %d", http.StatusCreated, w.Code)
	}
}

func TestMakeCreateHandler_InvalidBody(t *testing.T) {
	mock := &mockService{}
	handler := makeCreateHandler(mock)

	req := httptest.NewRequest("POST", "/students", bytes.NewReader([]byte("x")))
	w := httptest.NewRecorder()
	handler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected %d got %d", http.StatusBadRequest, w.Code)
	}
}

func TestMakeCreateV2Handler(t *testing.T) {
	mock := &mockService{
		createFunc: func(name, lastName string, age int32) (*Student, error) {
			return &Student{ID: "1", Name: name, LastName: lastName, Age: age}, nil
		},
		cuervoFunc: func(payload any) ([]byte, error) {
			return []byte(`{"status":"ok"}`), nil
		},
	}

	handler := makeCreateV2Handler(mock)
	body := CreateRequestV2{
		Lopez: CreateRequest{Name: "John", LastName: "Doe", Age: 25},
		Cuervo: CuervoRequest{
			FullName:    "Juan Cuervo",
			Email:       "juan@test.com",
			PhoneNumber: "3001234567",
		},
		Lasso: LasoRequest{
			Nombre:    "Hamburguesa",
			Direccion: "Calle 1",
			Telefono:  "1234567",
		},
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/students/api/v2", bytes.NewReader(b))
	w := httptest.NewRecorder()
	handler(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected %d got %d", http.StatusCreated, w.Code)
	}
	if !bytes.Contains(w.Body.Bytes(), []byte(`"lopez"`)) {
		t.Fatalf("expected response to include lopez data: %s", w.Body.String())
	}
	if !bytes.Contains(w.Body.Bytes(), []byte(`"cuervo"`)) {
		t.Fatalf("expected response to include cuervo data: %s", w.Body.String())
	}
}

func TestMakeCreateV2Handler_InvalidBody(t *testing.T) {
	mock := &mockService{}
	handler := makeCreateV2Handler(mock)
	req := httptest.NewRequest("POST", "/students/api/v2", bytes.NewReader([]byte("{")))
	w := httptest.NewRecorder()
	handler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected %d got %d", http.StatusBadRequest, w.Code)
	}
}

func TestMakeGetAllHandler(t *testing.T) {
	expected := []Student{{ID: "1", Name: "John", LastName: "Doe", Age: 25}}
	mock := &mockService{
		getAllFunc: func() ([]Student, error) { return expected, nil },
	}

	handler := makeGetAllHandler(mock)
	req := httptest.NewRequest("GET", "/students", nil)
	w := httptest.NewRecorder()
	handler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected %d got %d", http.StatusOK, w.Code)
	}
}

func TestMakeGetAllHandler_Error(t *testing.T) {
	mock := &mockService{
		getAllFunc: func() ([]Student, error) { return nil, errors.New("db") },
	}
	handler := makeGetAllHandler(mock)
	req := httptest.NewRequest("GET", "/students", nil)
	w := httptest.NewRecorder()
	handler(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected %d got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestMakeGetHandler(t *testing.T) {
	mock := &mockService{
		getFunc: func(id string) (*Student, error) { return &Student{ID: id, Name: "John"}, nil },
	}

	handler := makeGetHandler(mock)
	req := httptest.NewRequest("GET", "/students/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	w := httptest.NewRecorder()
	handler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected %d got %d", http.StatusOK, w.Code)
	}
}

func TestMakeGetHandler_NotFound(t *testing.T) {
	mock := &mockService{
		getFunc: func(id string) (*Student, error) { return nil, errors.New("not found") },
	}

	handler := makeGetHandler(mock)
	req := httptest.NewRequest("GET", "/students/999", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "999"})
	w := httptest.NewRecorder()
	handler(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected %d got %d", http.StatusNotFound, w.Code)
	}
}

func TestMakeDeleteHandler(t *testing.T) {
	mock := &mockService{deleteFunc: func(id string) error { return nil }}
	handler := makeDeleteHandler(mock)
	req := httptest.NewRequest("DELETE", "/students/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	w := httptest.NewRecorder()
	handler(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected %d got %d", http.StatusNoContent, w.Code)
	}
}

func TestMakeDeleteHandler_NotFound(t *testing.T) {
	mock := &mockService{deleteFunc: func(id string) error { return errors.New("notfound") }}
	handler := makeDeleteHandler(mock)
	req := httptest.NewRequest("DELETE", "/students/9", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "9"})
	w := httptest.NewRecorder()
	handler(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected %d got %d", http.StatusNotFound, w.Code)
	}
}

func TestMakePatchHandler(t *testing.T) {
	mock := &mockService{
		patchFunc: func(id string, name *string, lastName *string, age *int32) error { return nil },
	}
	handler := makePatchHandler(mock)
	name := "Jane"
	body := PatchRequest{Name: &name}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest("PATCH", "/students/1", bytes.NewReader(b))
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	w := httptest.NewRecorder()
	handler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected %d got %d", http.StatusOK, w.Code)
	}
}

func TestMakePatchHandler_InvalidBody(t *testing.T) {
	mock := &mockService{}
	handler := makePatchHandler(mock)
	req := httptest.NewRequest("PATCH", "/students/1", bytes.NewReader([]byte("{")))
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	w := httptest.NewRecorder()
	handler(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected %d got %d", http.StatusBadRequest, w.Code)
	}
}

func TestMakePatchHandler_NameRequired(t *testing.T) {
	mock := &mockService{}
	handler := makePatchHandler(mock)
	name := ""
	body := PatchRequest{Name: &name}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest("PATCH", "/students/1", bytes.NewReader(b))
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	w := httptest.NewRecorder()
	handler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected %d got %d", http.StatusBadRequest, w.Code)
	}
}

func TestMakePatchHandler_NotFound(t *testing.T) {
	mock := &mockService{
		patchFunc: func(id string, name *string, lastName *string, age *int32) error {
			return errors.New("missing")
		},
	}
	handler := makePatchHandler(mock)
	name := "Jane"
	body := PatchRequest{Name: &name}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest("PATCH", "/students/1", bytes.NewReader(b))
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	w := httptest.NewRecorder()
	handler(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected %d got %d", http.StatusNotFound, w.Code)
	}
}

func TestMakePutHandler(t *testing.T) {
	mock := &mockService{
		putFunc: func(id, name, lastName string, age int32) (*Student, error) {
			return &Student{ID: id, Name: name, LastName: lastName, Age: age}, nil
		},
	}

	handler := makePutHandler(mock)
	body := PutRequest{Name: "Jane", LastName: "Smith", Age: 30}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest("PUT", "/students/1", bytes.NewReader(b))
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	w := httptest.NewRecorder()
	handler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected %d got %d", http.StatusOK, w.Code)
	}
}

func TestMakePutHandler_InvalidBody(t *testing.T) {
	mock := &mockService{}
	handler := makePutHandler(mock)
	req := httptest.NewRequest("PUT", "/students/1", bytes.NewReader([]byte("{")))
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	w := httptest.NewRecorder()
	handler(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected %d got %d", http.StatusBadRequest, w.Code)
	}
}

func TestMakePutHandler_ServiceError(t *testing.T) {
	mock := &mockService{
		putFunc: func(id, name, lastName string, age int32) (*Student, error) {
			return nil, errors.New("bad")
		},
	}
	handler := makePutHandler(mock)
	body := PutRequest{Name: "A", LastName: "B", Age: 1}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest("PUT", "/students/1", bytes.NewReader(b))
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	w := httptest.NewRecorder()
	handler(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected %d got %d", http.StatusBadRequest, w.Code)
	}
}
