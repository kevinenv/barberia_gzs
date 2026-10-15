package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"

	"github.com/kevinenv/barberia_gzs/backend/internal/handler"
	"github.com/kevinenv/barberia_gzs/backend/internal/repository"
)

func main() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Println("Aviso: Usando variables de entorno del sistema")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(127.0.0.1:%s)/%s?parseTime=true&charset=utf8mb4",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		log.Fatalf(" Error conectando a MySQL: %v", err)
	}
	defer db.Close()

	log.Println(" Conexión exitosa a MySQL en Docker")

	// =========================================================================
	// INYECCIÓN DE DEPENDENCIAS (Conexión de Capas: DB -> Repo -> Handler)
	// =========================================================================
	catalogoRepo := repository.NewCatalogoRepository(db)
	catalogoHandler := handler.NewCatalogoHandler(catalogoRepo)

	// Inicializar Router de Gin
	r := gin.Default()

	// Agrupación de rutas versión 1 (/api/v1)
	v1 := r.Group("/api/v1")
	{
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok", "proyecto": "Barbería GZS"})
		})

		// Rutas públicas del catálogo
		v1.GET("/configuracion", catalogoHandler.GetConfiguracion)
		v1.GET("/servicios", catalogoHandler.GetServicios)
		v1.GET("/productos", catalogoHandler.GetProductos)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf(" Servidor corriendo en http://localhost:%s", port)
	r.Run(":" + port)
}
