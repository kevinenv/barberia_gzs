package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kevinenv/barberia_gzs/backend/internal/domain"
	"github.com/kevinenv/barberia_gzs/backend/internal/repository"
	"github.com/kevinenv/barberia_gzs/backend/internal/service"
)

type CitaHandler struct {
	citaService *service.CitaService
	citaRepo    *repository.CitaRepository
}

func NewCitaHandler(citaService *service.CitaService, citaRepo *repository.CitaRepository) *CitaHandler {
	return &CitaHandler{citaService: citaService, citaRepo: citaRepo}
}

// GetBarberos lista los barberos activos para que el cliente elija en la web
func (h *CitaHandler) GetBarberos(c *gin.Context) {
	barberos, err := h.citaRepo.ObtenerBarberosActivos()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener barberos"})
		return
	}
	c.JSON(http.StatusOK, barberos)
}

// CrearCita agenda una cita nueva (Requiere estar logueado con JWT)
func (h *CitaHandler) CrearCita(c *gin.Context) {
	clienteID := c.GetInt("usuario_id")

	var req domain.CrearCitaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos incompletos: envía fecha, hora_inicio, servicios_ids y metodo_pago_garantia"})
		return
	}

	cita, err := h.citaService.AgendarCita(clienteID, req)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"mensaje": "Cita reservada con éxito",
		"cita":    cita,
	})
}

// MisCitas devuelve el historial de cortes del cliente logueado
func (h *CitaHandler) MisCitas(c *gin.Context) {
	clienteID := c.GetInt("usuario_id")
	citas, err := h.citaRepo.ObtenerMisCitas(clienteID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener tu historial de citas"})
		return
	}
	c.JSON(http.StatusOK, citas)
}
