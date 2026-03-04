package course

import (
	"net/http"
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

func makePatchHandler(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {

	}

}

func makePutHandler(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
