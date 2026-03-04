package enrollment

import (
	"encoding/json"
	"net/http"
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

	}
}

func makeGetHandler(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func makeDeleteHandler(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
