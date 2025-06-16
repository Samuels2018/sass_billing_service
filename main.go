package main

import (
	"database/sql"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	_ "github.com/lib/pq"

	"sass-billing-service/src/config"
	"sass-billing-service/src/controllers"
	"sass-billing-service/src/repositories"
	router "sass-billing-service/src/routes"
	"sass-billing-service/src/services"
)

func main() {
	// Cargar configuración
	cfg := config.LoadConfig()

	// Conectar a PostgreSQL
	db, err := sql.Open("postgres",
		"host="+cfg.DBHost+" port="+cfg.DBPort+" user="+cfg.DBUser+
			" password="+cfg.DBPassword+" dbname="+cfg.DBName+" sslmode=disable")
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()

	// Verificar conexión
	if err := db.Ping(); err != nil {
		log.Fatalf("Error pinging database: %v", err)
	}

	// Inicializar repositorio, servicio y controlador
	invoiceRepo := repositories.NewInvoiceRepository(db)
	invoiceService := services.NewInvoiceService(invoiceRepo)
	invoiceController := controllers.NewInvoiceController(invoiceService)

	// Crear aplicación Fiber
	app := fiber.New()
	app.Use(logger.New())

	// Rutas
	api := app.Group("/api")
	router.SetupRoutes(api, invoiceController)

	// Iniciar servidor
	port := ":" + cfg.ServerPort
	if err := app.Listen(port); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

	log.Printf("Server running on port %s", cfg.ServerPort)
}
