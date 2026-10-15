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
)

func main() {
	// 1. Cargar las variables de entorno desde la raíz del proyecto (../.env)
	if err := godotenv.Load("../.env"); err != nil {
		log.Println("Aviso: No se encontró archivo ../.env, usando variables del sistema")
	}

	// 2. Armar la cadena de conexión (DSN) para MySQL
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// parseTime=true es obligatorio en Go para convertir DATETIME de MySQL a time.Time de Go
	dsn := fmt.Sprintf("%s:%s@tcp(127.0.0.1:%s)/%s?parseTime=true&charset=utf8mb4",
		dbUser, dbPass, dbPort, dbName)

	// 3. Conectar y hacer Ping real a MySQL en Docker
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		log.Fatalf(" Error conectando a MySQL en Docker: %v", err)
	}
	defer db.Close()

	log.Println(" Conexión exitosa a MySQL en Docker")

	// 4. Inicializar el servidor HTTP con Gin
	r := gin.Default()

	// Endpoint de prueba de salud (Healthcheck de producción)
	r.GET("/api/v1/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":   "ok",
			"proyecto": "Barbería GZS",
			"database": "MySQL conectada correctamente",
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf(" Servidor corriendo en http://localhost:%s", port)
	r.Run(":" + port)
}
