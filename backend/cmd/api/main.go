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
	"github.com/kevinenv/barberia_gzs/backend/internal/middleware"
	"github.com/kevinenv/barberia_gzs/backend/internal/repository"
	"github.com/kevinenv/barberia_gzs/backend/internal/service"
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
		log.Fatalf("❌ Error conectando a MySQL: %v", err)
	}
	defer db.Close()

	log.Println("✅ Conexión exitosa a MySQL en Docker")

	// Inyección de Dependencias: Catálogo
	catalogoRepo := repository.NewCatalogoRepository(db)
	catalogoHandler := handler.NewCatalogoHandler(catalogoRepo)

	// Inyección de Dependencias: Autenticación y Usuarios
	usuarioRepo := repository.NewUsuarioRepository(db)
	authService := service.NewAuthService(usuarioRepo)
	authHandler := handler.NewAuthHandler(authService, usuarioRepo)

	r := gin.Default()

	v1 := r.Group("/api/v1")
	{
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok", "proyecto": "Barbería GZS"})
		})

		// Catálogo Público
		v1.GET("/configuracion", catalogoHandler.GetConfiguracion)
		v1.GET("/servicios", catalogoHandler.GetServicios)
		v1.GET("/productos", catalogoHandler.GetProductos)

		// Autenticación Pública
		v1.POST("/auth/registro", authHandler.Registrar)
		v1.POST("/auth/login", authHandler.Login)

		// Rutas Protegidas (Requieren Token JWT válido)
		protegido := v1.Group("")
		protegido.Use(middleware.ProtegerRuta())
		{
			protegido.GET("/auth/perfil", authHandler.MiPerfil)
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Servidor corriendo en http://localhost:%s", port)
	r.Run(":" + port)
}
