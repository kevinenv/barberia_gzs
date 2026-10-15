package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kevinenv/barberia_gzs/backend/internal/repository"
)

type CatalogoHandler struct {
	repo *repository.CatalogoRepository
}

func NewCatalogoHandler(repo *repository.CatalogoRepository) *CatalogoHandler {
	return &CatalogoHandler{repo: repo}
}

// GetConfiguracion responde en GET /api/v1/configuracion
func (h *CatalogoHandler) GetConfiguracion(c *gin.Context) {
	config, err := h.repo.ObtenerConfiguracion()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo obtener la configuración de la barbería"})
		return
	}
	c.JSON(http.StatusOK, config)
}

// GetServicios responde en GET /api/v1/servicios
func (h *CatalogoHandler) GetServicios(c *gin.Context) {
	servicios, err := h.repo.ObtenerServiciosActivos()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al consultar los servicios"})
		return
	}
	c.JSON(http.StatusOK, servicios)
}

// GetProductos responde en GET /api/v1/productos
func (h *CatalogoHandler) GetProductos(c *gin.Context) {
	productos, err := h.repo.ObtenerProductosActivos()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al consultar el inventario de productos"})
		return
	}
	c.JSON(http.StatusOK, productos)
}
