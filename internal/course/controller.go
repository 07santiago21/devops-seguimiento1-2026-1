package course

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
		Put    Controller
	}

	ErrorResponse struct {
		Error string `json:"error"`
	}
)

func Handler(s Service) *Endpoints {
	return &Endpoints{
		Create: makeCreateHandler(s),
		Get:    makeGetHandler(s),
		GetAll: makeGetAllHandler(s),
		Delete: makeDeleteHandler(s),
		Patch:  makePatchHandler(s),
		Put:    makePutHandler(s),
	}
}

type CreateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Credits     int32  `json:"credits"`
	Capacity    int32  `json:"capacity"`
}

func makeCreateHandler(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var req CreateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid request format"})
			return
		}

		course, err := s.Create(req.Name, req.Description, req.Credits, req.Capacity)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(course)
	}
}

func makeGetAllHandler(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		courses, err := s.GetAll()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(courses)
	}
}

func makeGetHandler(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		id := mux.Vars(r)["id"]

		course, err := s.Get(id)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "course not found"})
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(course)
	}
}

func makeDeleteHandler(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		if err := s.Delete(id); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "course not found"})
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

type PatchRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Credits     *int32  `json:"credits"`
	Capacity    *int32  `json:"capacity"`
}

func makePatchHandler(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var req PatchRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid request format"})
			return
		}
		id := mux.Vars(r)["id"]
		if err := s.Patch(id, req.Name, req.Description, req.Credits, req.Capacity); err != nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "course not found"})
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "course updated successfully"})
	}
}

type PutRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Credits     int32  `json:"credits"`
	Capacity    int32  `json:"capacity"`
}

func makePutHandler(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var req PutRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid request format"})
			return
		}
		id := mux.Vars(r)["id"]
		course, err := s.Put(id, req.Name, req.Description, req.Credits, req.Capacity)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(course)
	}
}
