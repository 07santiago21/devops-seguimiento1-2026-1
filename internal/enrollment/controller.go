package enrollment

import (
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
