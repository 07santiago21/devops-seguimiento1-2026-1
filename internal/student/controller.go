package student

import (
	"encoding/json"
	"net/http"
)

type (
	Controller func(w http.ResponseWriter, r *http.Request)

	Endpoints struct {
		Create Controller
	}

	createRequest struct {
		Name     string `json:"name"`
		LastName string `json:"last_name"`
		Age      int32  `json:"age"`
	}

	errorResponse struct {
		Error string `json:"error"`
	}
)

func handler(s Service) *Endpoints {
	return &Endpoints{
		Create: makeCreateHandler(s),
	}
}

func makeCreateHandler(s Service) Controller {

	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		var req createRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errorResponse{Error: "invalid request body"})
			return
		}

		student, err := s.Create(req.Name, req.LastName, req.Age)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errorResponse{Error: err.Error()})
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(student)
	}
}
