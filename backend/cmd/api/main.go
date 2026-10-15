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
		log.Fatalf(" Error conectando a MySQL: %v", err)
	}
	defer db.Close()

	log.Println(" Conexión exitosa a MySQL en Docker")

	// 1. Inyección de Dependencias: Catálogo
	catalogoRepo := repository.NewCatalogoRepository(db)
	catalogoHandler := handler.NewCatalogoHandler(catalogoRepo)

	// 2. Inyección de Dependencias: Autenticación
	usuarioRepo := repository.NewUsuarioRepository(db)
	authService := service.NewAuthService(usuarioRepo)
	authHandler := handler.NewAuthHandler(authService, usuarioRepo)

	// 3. Inyección de Dependencias: Motor de Citas Híbrido
	citaRepo := repository.NewCitaRepository(db)
	citaService := service.NewCitaService(citaRepo, catalogoRepo)
	citaHandler := handler.NewCitaHandler(citaService, citaRepo)

	// 4. Inyección de Dependencias: Mantenimiento Hidráulico, Stock y Comisiones
	operacionesRepo := repository.NewOperacionesRepository(db)
	operacionesHandler := handler.NewOperacionesHandler(operacionesRepo)

	r := gin.Default()

	// Habilitar CORS para que Vue 3 pueda comunicarse con Go
	r.Use(middleware.ConfigurarCORS())

	v1 := r.Group("/api/v1")
	{
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok", "proyecto": "Barbería GZS"})
		})

		// Rutas Públicas
		v1.GET("/configuracion", catalogoHandler.GetConfiguracion)
		v1.GET("/servicios", catalogoHandler.GetServicios)
		v1.GET("/productos", catalogoHandler.GetProductos)
		v1.GET("/barberos", citaHandler.GetBarberos)

		// Autenticación
		v1.POST("/auth/registro", authHandler.Registrar)
		v1.POST("/auth/login", authHandler.Login)

		// Rutas Protegidas (Requieren Token JWT)
		protegido := v1.Group("")
		protegido.Use(middleware.ProtegerRuta())
		{
			protegido.GET("/auth/perfil", authHandler.MiPerfil)

			// Citas de Barbería
			protegido.POST("/citas", citaHandler.CrearCita)
			protegido.GET("/citas/mis-citas", citaHandler.MisCitas)

			// Mantenimiento de Sillas Hidráulicas
			protegido.POST("/mantenimiento", operacionesHandler.SolicitarMantenimiento)
			protegido.GET("/mantenimiento/mis-solicitudes", operacionesHandler.MisMantenimientos)

			// Inventario (Venta de productos) y Reporte de Comisiones (50/50)
			protegido.POST("/ventas-productos", operacionesHandler.VenderProducto)
			protegido.GET("/caja/comisiones", operacionesHandler.ReporteComisiones)
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Servidor corriendo en http://localhost:%s", port)
	r.Run(":" + port)
}
