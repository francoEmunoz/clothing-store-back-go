package main

import (
	"cs-go/routes"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {

	// models.MigrarUser()
	// models.MigrarProduct()
	// models.MigrarQuestion()

	router := mux.NewRouter()

	routes.SetUserRoutes(router)
	routes.SetProductRoutes(router)
	routes.SetQuestionRoutes(router)

	// Configuración CORS personalizada
	corsOptions := cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}

	c := cors.New(corsOptions).Handler(router)

	// Configurar servidor HTTP
	srv := &http.Server{
		Handler: c,
		Addr:    ":4002",
	}

	log.Fatal(srv.ListenAndServe())
}
