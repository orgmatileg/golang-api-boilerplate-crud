package routes

import (
	"fmt"
	"net/http"

	"github.com/getsentry/raven-go"

	"golang-api-boilerplate-crud/handlers"
	"golang-api-boilerplate-crud/middleware"

	"github.com/gorilla/mux"
)

// Initiation ...
func Initiation() *mux.Router {

	r := mux.NewRouter()

	// for check health purpose
	r.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Pong!")
	}).Methods("GET")

	// Versioning API
	r1 := r.PathPrefix("/v1").Subrouter()

	// Middleware
	r1.Use(middleware.CORS)

	// Pets
	r1.HandleFunc("/fruits", raven.RecoveryHandler(handlers.NewFruitsHandlers().Create)).Methods("POST")
	r1.HandleFunc("/fruits", raven.RecoveryHandler(handlers.NewFruitsHandlers().GetAll)).Methods("GET")
	r1.HandleFunc("/fruits/{id}", raven.RecoveryHandler(handlers.NewFruitsHandlers().GetOneByID)).Methods("GET")
	r1.HandleFunc("/fruits/{id}", raven.RecoveryHandler(handlers.NewFruitsHandlers().UpdateOneByID)).Methods("PUT")
	r1.HandleFunc("/fruits/{id}/uploadImage", raven.RecoveryHandler(handlers.NewFruitsHandlers().UpdatePhotoOneByID)).Methods("POST")
	r1.HandleFunc("/fruits/{id}", raven.RecoveryHandler(handlers.NewFruitsHandlers().DeleteOneByID)).Methods("DELETE")

	return r
}
