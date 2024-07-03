package routes

import (
	"cs-go/handlers"
	"net/http"

	"github.com/gorilla/mux"
)

func SetUserRoutes(r *mux.Router) {
	r.HandleFunc("/api/user/", handlers.GetUsers).Methods("GET")
	r.Handle("/api/user/{id:[0-9]+}", handlers.AuthMiddleware(http.HandlerFunc(handlers.GetUser))).Methods("GET")
	r.HandleFunc("/api/signup/", handlers.SignUp).Methods("POST")
	r.HandleFunc("/api/login/", handlers.LogIn).Methods("POST")
	r.Handle("/api/token/", handlers.AuthMiddleware(http.HandlerFunc(handlers.VerifyToken))).Methods("GET")
	r.HandleFunc("/api/logout/{id:[0-9]+}", handlers.LogOut).Methods("POST")
	r.Handle("/api/user/{id:[0-9]+}", handlers.AuthMiddleware(http.HandlerFunc(handlers.UpdateUser))).Methods("PUT")
}
