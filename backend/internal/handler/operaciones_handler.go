package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kevinenv/barberia_gzs/backend/internal/domain"
	"github.com/kevinenv/barberia_gzs/backend/internal/repository"
)

type OperacionesHandler struct {
	repo *repository.OperacionesRepository
}

func NewOperacionesHandler(repo *repository.OperacionesRepository) *OperacionesHandler {
	return &OperacionesHandler{repo: repo}
}

func (h *OperacionesHandler) SolicitarMantenimiento(c *gin.Context) {
	clienteID := c.GetInt("usuario_id")
	var req domain.CrearMantenimientoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ingresa la marca/modelo de tu silla hidráulica y describe la falla (mín. 10 caracteres)"})
		return
	}

	sol, err := h.repo.CrearSolicitudMantenimiento(clienteID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo registrar la solicitud de mantenimiento"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"mensaje":   "Solicitud de mantenimiento registrada. Te enviaremos la cotización pronto.",
		"solicitud": sol,
	})
}

func (h *OperacionesHandler) MisMantenimientos(c *gin.Context) {
	clienteID := c.GetInt("usuario_id")
	lista, err := h.repo.ObtenerMisSolicitudesMantenimiento(clienteID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al consultar solicitudes"})
		return
	}
	c.JSON(http.StatusOK, lista)
}

func (h *OperacionesHandler) VenderProducto(c *gin.Context) {
	usuarioID := c.GetInt("usuario_id")
	var req domain.VentaProductoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos de venta inválidos"})
		return
	}

	subtotal, stockRestante, err := h.repo.RegistrarVentaProducto(usuarioID, req)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"mensaje":        "Venta registrada y stock actualizado",
		"subtotal_bs":    subtotal,
		"stock_restante": stockRestante,
	})
}

func (h *OperacionesHandler) ReporteComisiones(c *gin.Context) {
	reporte, err := h.repo.ObtenerReporteComisiones()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al calcular comisiones"})
		return
	}
	c.JSON(http.StatusOK, reporte)
}
