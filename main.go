package main

import (
	"API_wrkf/config"
	"API_wrkf/handlers"
	"API_wrkf/routes"
	"API_wrkf/services"
	"API_wrkf/storage"
	"fmt"
	"log"

	"github.com/labstack/echo/v4"
)

func main() {
	// Cargar la configuración de la aplicación desde las variables de entorno
	cfg := config.LoadConfig()

	// Create a new database connection
	db, err := storage.NewConnection(cfg.DB)
	if err != nil {
		log.Fatalf("No se pudo conectar a la base de datos: %v", err)
	}

	fmt.Println("Se conecto con exito a la db")

	// Run database migrations
	if err := storage.Migrate(db); err != nil {
		log.Fatalf("No se pudo migrar la base de datos: %v", err)
	}

	fmt.Println("Database migration completed successfully!")

	// --- Initialize Layers with Dependencies ---
	userRepo := storage.NewUserRepository(db)
	userService := services.NewUserService(userRepo, cfg.JWTSecret)
	userHandler := handlers.NewUserHandler(userService)

	// --- Initialize Echo and Set Up Routes ---
	e := echo.New()
	routes.SetupRoutes(e, userHandler, cfg.JWTSecret)

	// --- Start Server ---
	fmt.Println("Starting server on port 8080...")
	if err := e.Start(":8080"); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}
