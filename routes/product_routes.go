package routes

import (
	"cs-go/handlers"

	"github.com/gorilla/mux"
)

func SetProductRoutes(r *mux.Router) {
	r.HandleFunc("/api/product/", handlers.GetProducts).Methods("GET")
	r.HandleFunc("/api/product/{id:[0-9]+}", handlers.GetProduct).Methods("GET")
	r.HandleFunc("/api/product/", handlers.CreateProduct).Methods("POST")
	r.HandleFunc("/api/product/{id:[0-9]+}", handlers.UpdateProduct).Methods("PUT")
	r.HandleFunc("/api/product/{id:[0-9]+}", handlers.DeleteProduct).Methods("DELETE")
}
