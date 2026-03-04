package course

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

type mockCourseService struct {
	createFunc func(name, description string, credits, capacity int32) (*Course, error)
	getAllFunc func() ([]Course, error)
	getFunc    func(id string) (*Course, error)
	deleteFunc func(id string) error
	patchFunc  func(id string, name *string, description *string, credits *int32, capacity *int32) error
	putFunc    func(id, name, description string, credits, capacity int32) (*Course, error)
}

func (m *mockCourseService) Create(name, description string, credits, capacity int32) (*Course, error) {
	if m.createFunc == nil {
		return nil, errors.New("not implemented")
	}
	return m.createFunc(name, description, credits, capacity)
}

func (m *mockCourseService) GetAll() ([]Course, error) {
	if m.getAllFunc == nil {
		return nil, errors.New("not implemented")
	}
	return m.getAllFunc()
}

func (m *mockCourseService) Get(id string) (*Course, error) {
	if m.getFunc == nil {
		return nil, errors.New("not implemented")
	}
	return m.getFunc(id)
}

func (m *mockCourseService) Delete(id string) error {
	if m.deleteFunc == nil {
		return errors.New("not implemented")
	}
	return m.deleteFunc(id)
}

func (m *mockCourseService) Patch(id string, name *string, description *string, credits *int32, capacity *int32) error {
	if m.patchFunc == nil {
		return errors.New("not implemented")
	}
	return m.patchFunc(id, name, description, credits, capacity)
}

func (m *mockCourseService) Put(id, name, description string, credits, capacity int32) (*Course, error) {
	if m.putFunc == nil {
		return nil, errors.New("not implemented")
	}
	return m.putFunc(id, name, description, credits, capacity)
}

// ========== CREATE TESTS ==========
func TestMakeCourseCreateHandler_Success(t *testing.T) {
	mock := &mockCourseService{
		createFunc: func(name, description string, credits, capacity int32) (*Course, error) {
			return &Course{ID: "1", Name: name, Description: description, Credits: credits, Capacity: capacity}, nil
		},
	}

	handler := makeCreateHandler(mock)
	body := CreateRequest{Name: "Math", Description: "Mathematics course", Credits: 4, Capacity: 30}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/courses", bytes.NewReader(b))
	w := httptest.NewRecorder()
	handler(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected %d got %d", http.StatusCreated, w.Code)
	}

	var resp Course
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.ID != "1" {
		t.Errorf("expected ID '1' got '%s'", resp.ID)
	}
}

func TestMakeCourseCreateHandler_InvalidBody(t *testing.T) {
	mock := &mockCourseService{}
	handler := makeCreateHandler(mock)
	req := httptest.NewRequest("POST", "/courses", bytes.NewReader([]byte("invalid json")))
	w := httptest.NewRecorder()
	handler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected %d got %d", http.StatusBadRequest, w.Code)
	}
}

func TestMakeCourseCreateHandler_ServiceError(t *testing.T) {
	mock := &mockCourseService{
		createFunc: func(name, description string, credits, capacity int32) (*Course, error) {
			return nil, errors.New("validation error")
		},
	}
	handler := makeCreateHandler(mock)
	body := CreateRequest{Name: "", Description: "", Credits: 0, Capacity: 0}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/courses", bytes.NewReader(b))
	w := httptest.NewRecorder()
	handler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected %d got %d", http.StatusBadRequest, w.Code)
	}
}

// ========== GETALL TESTS ==========
func TestMakeCourseGetAllHandler_Success(t *testing.T) {
	mock := &mockCourseService{
		getAllFunc: func() ([]Course, error) {
			return []Course{
				{ID: "1", Name: "Math", Description: "Math course", Credits: 4, Capacity: 30},
				{ID: "2", Name: "Physics", Description: "Physics course", Credits: 3, Capacity: 25},
			}, nil
		},
	}

	handler := makeGetAllHandler(mock)
	req := httptest.NewRequest("GET", "/courses", nil)
	w := httptest.NewRecorder()
	handler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected %d got %d", http.StatusOK, w.Code)
	}

	var resp []Course
	json.Unmarshal(w.Body.Bytes(), &resp)
	if len(resp) != 2 {
		t.Errorf("expected 2 courses got %d", len(resp))
	}
}

func TestMakeCourseGetAllHandler_Empty(t *testing.T) {
	mock := &mockCourseService{
		getAllFunc: func() ([]Course, error) { return []Course{}, nil },
	}

	handler := makeGetAllHandler(mock)
	req := httptest.NewRequest("GET", "/courses", nil)
	w := httptest.NewRecorder()
	handler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected %d got %d", http.StatusOK, w.Code)
	}
}

func TestMakeCourseGetAllHandler_DatabaseError(t *testing.T) {
	mock := &mockCourseService{
		getAllFunc: func() ([]Course, error) { return nil, errors.New("database error") },
	}
	handler := makeGetAllHandler(mock)
	req := httptest.NewRequest("GET", "/courses", nil)
	w := httptest.NewRecorder()
	handler(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected %d got %d", http.StatusInternalServerError, w.Code)
	}
}

// ========== GET BY ID TESTS ==========
func TestMakeCourseGetHandler_Success(t *testing.T) {
	mock := &mockCourseService{
		getFunc: func(id string) (*Course, error) {
			return &Course{ID: id, Name: "Math", Description: "Mathematics", Credits: 4, Capacity: 30}, nil
		},
	}

	handler := makeGetHandler(mock)
	req := httptest.NewRequest("GET", "/courses/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	w := httptest.NewRecorder()
	handler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected %d got %d", http.StatusOK, w.Code)
	}

	var resp Course
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.ID != "1" {
		t.Errorf("expected ID '1' got '%s'", resp.ID)
	}
}

func TestMakeCourseGetHandler_NotFound(t *testing.T) {
	mock := &mockCourseService{
		getFunc: func(id string) (*Course, error) { return nil, errors.New("record not found") },
	}

	handler := makeGetHandler(mock)
	req := httptest.NewRequest("GET", "/courses/999", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "999"})
	w := httptest.NewRecorder()
	handler(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected %d got %d", http.StatusNotFound, w.Code)
	}
}

// ========== DELETE TESTS ==========
func TestMakeCourseDeleteHandler_Success(t *testing.T) {
	mock := &mockCourseService{
		deleteFunc: func(id string) error { return nil },
	}
	handler := makeDeleteHandler(mock)
	req := httptest.NewRequest("DELETE", "/courses/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	w := httptest.NewRecorder()
	handler(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected %d got %d", http.StatusNoContent, w.Code)
	}
}

func TestMakeCourseDeleteHandler_NotFound(t *testing.T) {
	mock := &mockCourseService{
		deleteFunc: func(id string) error { return errors.New("record not found") },
	}
	handler := makeDeleteHandler(mock)
	req := httptest.NewRequest("DELETE", "/courses/999", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "999"})
	w := httptest.NewRecorder()
	handler(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected %d got %d", http.StatusNotFound, w.Code)
	}
}

// ========== PATCH TESTS ==========
func TestMakeCoursePatchHandler_Success(t *testing.T) {
	mock := &mockCourseService{
		patchFunc: func(id string, name *string, description *string, credits *int32, capacity *int32) error {
			return nil
		},
	}
	handler := makePatchHandler(mock)
	name := "Advanced Math"
	body := PatchRequest{Name: &name}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest("PATCH", "/courses/1", bytes.NewReader(b))
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	w := httptest.NewRecorder()
	handler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected %d got %d", http.StatusOK, w.Code)
	}
}

func TestMakeCoursePatchHandler_PartialUpdate(t *testing.T) {
	mock := &mockCourseService{
		patchFunc: func(id string, name *string, description *string, credits *int32, capacity *int32) error {
			return nil
		},
	}
	handler := makePatchHandler(mock)
	credits := int32(5)
	capacity := int32(40)
	body := PatchRequest{Credits: &credits, Capacity: &capacity}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest("PATCH", "/courses/1", bytes.NewReader(b))
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	w := httptest.NewRecorder()
	handler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected %d got %d", http.StatusOK, w.Code)
	}
}

func TestMakeCoursePatchHandler_InvalidBody(t *testing.T) {
	mock := &mockCourseService{}
	handler := makePatchHandler(mock)
	req := httptest.NewRequest("PATCH", "/courses/1", bytes.NewReader([]byte("invalid")))
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	w := httptest.NewRecorder()
	handler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected %d got %d", http.StatusBadRequest, w.Code)
	}
}

func TestMakeCoursePatchHandler_NotFound(t *testing.T) {
	mock := &mockCourseService{
		patchFunc: func(id string, name *string, description *string, credits *int32, capacity *int32) error {
			return errors.New("record not found")
		},
	}
	handler := makePatchHandler(mock)
	name := "Test"
	body := PatchRequest{Name: &name}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest("PATCH", "/courses/999", bytes.NewReader(b))
	req = mux.SetURLVars(req, map[string]string{"id": "999"})
	w := httptest.NewRecorder()
	handler(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected %d got %d", http.StatusNotFound, w.Code)
	}
}

// ========== PUT TESTS ==========
func TestMakeCoursePutHandler_Success(t *testing.T) {
	mock := &mockCourseService{
		putFunc: func(id, name, description string, credits, capacity int32) (*Course, error) {
			return &Course{ID: id, Name: name, Description: description, Credits: credits, Capacity: capacity}, nil
		},
	}

	handler := makePutHandler(mock)
	body := PutRequest{Name: "Physics", Description: "Physics course", Credits: 3, Capacity: 25}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest("PUT", "/courses/1", bytes.NewReader(b))
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	w := httptest.NewRecorder()
	handler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected %d got %d", http.StatusOK, w.Code)
	}

	var resp Course
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Name != "Physics" {
		t.Errorf("expected Name 'Physics' got '%s'", resp.Name)
	}
}

func TestMakeCoursePutHandler_InvalidBody(t *testing.T) {
	mock := &mockCourseService{}
	handler := makePutHandler(mock)
	req := httptest.NewRequest("PUT", "/courses/1", bytes.NewReader([]byte("{")))
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	w := httptest.NewRecorder()
	handler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected %d got %d", http.StatusBadRequest, w.Code)
	}
}
func int32Ptr(i int32) *int32 {
	return &i
}
func TestMakeCoursePutHandler_ServiceValidationError(t *testing.T) {
	mock := &mockCourseService{
		putFunc: func(id, name, description string, credits, capacity int32) (*Course, error) {
			return nil, errors.New("name is required")
		},
	}
	handler := makePutHandler(mock)
	body := PutRequest{Name: "", Description: "", Credits: 0, Capacity: 0}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest("PUT", "/courses/1", bytes.NewReader(b))
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	w := httptest.NewRecorder()
	handler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected %d got %d", http.StatusBadRequest, w.Code)
	}
}

// ========== HANDLER FACTORY TEST ==========
func TestHandler_CreatesCourseEndpoints(t *testing.T) {
	mock := &mockCourseService{}
	endpoints := Handler(mock)

	if endpoints == nil {
		t.Fatal("expected endpoints not to be nil")
	}
	if endpoints.Create == nil {
		t.Error("Create endpoint is nil")
	}
	if endpoints.GetAll == nil {
		t.Error("GetAll endpoint is nil")
	}
	if endpoints.Get == nil {
		t.Error("Get endpoint is nil")
	}
	if endpoints.Delete == nil {
		t.Error("Delete endpoint is nil")
	}
	if endpoints.Patch == nil {
		t.Error("Patch endpoint is nil")
	}
	if endpoints.Put == nil {
		t.Error("Put endpoint is nil")
	}
}
