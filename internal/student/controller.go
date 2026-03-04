package student

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type (
	Controller func(w http.ResponseWriter, r *http.Request)

	Endpoints struct {
		Create Controller
		Get    Controller
		GetAll Controller
		Delete Controller
		Patch  Controller
	}

	CreateRequest struct {
		Name     string `json:"name"`
		LastName string `json:"last_name"`
		Age      int32  `json:"age"`
	}

	ErrorResponse struct {
		Error string `json:"error"`
	}

	PatchRequest struct {
		Name     *string `json:"name"`
		LastName *string `json:"last_name"`
		Age      *int32  `json:"age"`
	}
)

func Handler(s Service) *Endpoints {
	return &Endpoints{
		Create: makeCreateHandler(s),
		Get:    makeGetHandler(s),
		GetAll: makeGetAllHandler(s),
		Delete: makeDeleteHandler(s),
		Patch:  makePatchHandler(s),
	}
}

func makeCreateHandler(s Service) Controller {

	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		var req CreateRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid request body"})
			return
		}

		student, err := s.Create(req.Name, req.LastName, req.Age)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(student)
	}
}

func makeGetAllHandler(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {

		users, err := s.GetAll()
		if err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(ErrorResponse{err.Error()})
			return

		}

		json.NewEncoder(w).Encode(users)

	}
}

func makeGetHandler(s Service) Controller {

	return func(w http.ResponseWriter, r *http.Request) {

		path := mux.Vars(r)
		id := path["id"]
		user, err := s.Get(id)
		if err != nil {
			w.WriteHeader(404)
			json.NewEncoder(w).Encode(ErrorResponse{"user does not exist"})
			return
		}

		json.NewEncoder(w).Encode(user)

	}

}

func makeDeleteHandler(s Service) Controller {

	return func(w http.ResponseWriter, r *http.Request) {

		path := mux.Vars(r)
		id := path["id"]
		if err := s.Delete(id); err != nil {
			w.WriteHeader(404)
			json.NewEncoder(w).Encode(ErrorResponse{"user does not exist"})
			return
		}

		json.NewEncoder(w).Encode(map[string]string{"data": "success"})

	}
}

func makePatchHandler(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		var req PatchRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(ErrorResponse{"invalid request format"})
			return
		}

		if req.Name != nil && *req.Name == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(ErrorResponse{"name is required"})
			return
		}

		path := mux.Vars(r)
		id := path["id"]
		if err := s.Patch(id, req.Name, req.LastName, req.Age); err != nil {
			w.WriteHeader(404)
			json.NewEncoder(w).Encode(ErrorResponse{"user does not exist"})
			return
		}

		json.NewEncoder(w).Encode(map[string]string{"data": "success"})

	}

}
