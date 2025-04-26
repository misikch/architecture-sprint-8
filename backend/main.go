package main

import (
	"log"
	"net/http"
	"os"

	"github.com/rs/cors"

	"github.com/misikch/architecture-sprint-8/backend/api/handler/report"
	"github.com/misikch/architecture-sprint-8/backend/api/middleware"
)

func main() {
	cfg := middleware.Config{
		KeycloakURL:  os.Getenv("KEYCLOAK_URL"),
		Realm:        os.Getenv("KEYCLOAK_REALM"),
		ClientId:     os.Getenv("KEYCLOAK_CLIENT_ID"),
		ClientSecret: os.Getenv("KEYCLOAK_CLIENT_SECRET"),
	}

	// debug
	//cfg := middleware.Config{
	//	KeycloakURL:  "http://localhost:8080",
	//	Realm:        "reports-realm",
	//	ClientId:     "reports-api",
	//	ClientSecret: "oNwoLQdvJAvRcL89SydqCWCe5ry1jMgq",
	//}

	authMiddleware := middleware.NewAuthMiddleware(cfg)

	// Настройки CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
		Debug:            true,
	})

	mux := http.NewServeMux()
	mux.Handle("/reports", authMiddleware.Middleware(
		http.HandlerFunc(report.ReportsHandler),
	))

	handler := c.Handler(mux)

	port := "8000"
	log.Printf("Starting server on port %s", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
