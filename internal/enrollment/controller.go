package enrollment

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type (
	Controller func(w http.ResponseWriter, r *http.Request)
	Endpoints  struct {
		Create Controller
		GetAll Controller
		Get    Controller
		Delete Controller
	}
)

func MakeEndpoints(s Service) Endpoints {
	return Endpoints{
		Create: makeCreateHandler(s),
		GetAll: makeGetAllHandler(s),
		Get:    makeGetHandler(s),
		Delete: makeDeleteHandler(s),
	}
}

func makeCreateHandler(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			StudentID   string  `json:"student_id"`
			CourseID    string  `json:"course_id"`
			TotalAmount float64 `json:"total_amount"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid json body"})
			return
		}

		e, err := s.Create(req.StudentID, req.CourseID, req.TotalAmount)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(e)
	}
}

func makeGetAllHandler(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		enrollments, err := s.GetAll()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "could not fetch enrollments"})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(enrollments)
	}
}

func makeGetHandler(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		e, err := s.Get(id)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "enrollment not found"})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(e)
	}
}

func makeDeleteHandler(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]

		err := s.Delete(id)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(map[string]string{"error": "enrollment not found"})
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "could not delete enrollment"})
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
