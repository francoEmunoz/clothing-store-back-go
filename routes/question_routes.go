package routes

import (
	"cs-go/handlers"
	"net/http"

	"github.com/gorilla/mux"
)

func SetQuestionRoutes(r *mux.Router) {
	r.HandleFunc("/api/question/", handlers.GetQuestions).Methods("GET")
	r.Handle("/api/question/", handlers.AuthMiddleware(http.HandlerFunc(handlers.CreateQuestion))).Methods("POST")
	r.Handle("/api/comment/{id:[0-9]+}", handlers.AuthMiddleware(http.HandlerFunc(handlers.UpdateComment))).Methods("PUT")
	r.Handle("/api/question/{id:[0-9]+}", handlers.AuthMiddleware(http.HandlerFunc(handlers.DeleteQuestion))).Methods("DELETE")
}
