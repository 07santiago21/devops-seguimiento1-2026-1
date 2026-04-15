package student

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type (
	Controller func(w http.ResponseWriter, r *http.Request)

	Endpoints struct {
		Create   Controller
		CreateV2 Controller
		Get      Controller
		GetAll   Controller
		Delete   Controller
		Patch    Controller
		Put      Controller
	}

	CreateRequestV2 struct {
		Lopez  CreateRequest `json:"lopez"`
		Cuervo CuervoRequest `json:"cuervo"`
		Lasso  LasoRequest   `json:"lasso"`
	}

	LasoRequest struct {
		Nombre    string `json:"nombre"`
		Direccion string `json:"direccion"`
		Telefono  string `json:"telefono"`
	}

	CuervoRequest struct {
		FullName    string `json:"fullName"`
		Email       string `json:"email"`
		PhoneNumber string `json:"phone_number"`
	}

	CreateRequest struct {
		Name     string `json:"name"`
		LastName string `json:"last_name"`
		Age      int32  `json:"age"`
	}

	ErrorResponse struct {
		Error string `json:"error"`
	}

	ValidationErrorResponse struct {
		Error   string            `json:"error"`
		Details map[string]string `json:"details"`
	}

	PatchRequest struct {
		Name     *string `json:"name"`
		LastName *string `json:"last_name"`
		Age      *int32  `json:"age"`
	}
	PutRequest struct {
		Name     string `json:"name"`
		LastName string `json:"last_name"`
		Age      int32  `json:"age"`
	}
)

func Handler(s Service) *Endpoints {
	return &Endpoints{
		Create:   makeCreateHandler(s),
		CreateV2: makeCreateV2Handler(s),
		Get:      makeGetHandler(s),
		GetAll:   makeGetAllHandler(s),
		Delete:   makeDeleteHandler(s),
		Patch:    makePatchHandler(s),
		Put:      makePutHandler(s),
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

func makeCreateV2Handler(s Service) Controller {

	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		var req CreateRequestV2

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid request body"})
			return
		}

		if details := validateCreateV2RequiredFields(req); len(details) > 0 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ValidationErrorResponse{
				Error:   "missing or invalid required fields",
				Details: details,
			})
			return
		}

		payload := map[string]interface{}{
			"lasso":  req.Lasso,
			"cuervo": req.Cuervo,
		}

		cuervoResp, err := s.Cuervo(payload)

		var cuervoData interface{}

		if err != nil {
			cuervoData = map[string]interface{}{
				"error": err.Error(),
			}
		} else {
			cuervoData = json.RawMessage(cuervoResp)
		}

		reqStudent := req.Lopez

		student, err := s.Create(reqStudent.Name, reqStudent.LastName, reqStudent.Age)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"lopez": err.Error(), "cuervo": cuervoData})
			return
		}

		response := map[string]interface{}{
			"lopez":  student,
			"cuervo": cuervoData,
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	}
}

func validateCreateV2RequiredFields(req CreateRequestV2) map[string]string {
	details := map[string]string{}

	if strings.TrimSpace(req.Lopez.Name) == "" {
		details["lopez.name"] = "is required"
	}
	if strings.TrimSpace(req.Lopez.LastName) == "" {
		details["lopez.last_name"] = "is required"
	}
	if req.Lopez.Age <= 0 {
		details["lopez.age"] = "must be greater than zero"
	}

	if req.Cuervo == (CuervoRequest{}) {
		details["cuervo"] = "is required"
	}

	if req.Lasso == (LasoRequest{}) {
		details["lasso"] = "is required"
	}

	return details
}

func makeGetAllHandler(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {

		users, err := s.GetAll()
		if err != nil {
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(ErrorResponse{err.Error()})
			return

		}

		w.WriteHeader(http.StatusOK)
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

		w.WriteHeader(http.StatusNoContent)

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

func makePutHandler(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		var req PutRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{"invalid request format"})
			return
		}

		path := mux.Vars(r)
		id := path["id"]

		student, err := s.Put(id, req.Name, req.LastName, req.Age)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(student)
	}
}
