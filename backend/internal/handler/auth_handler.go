package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kevinenv/barberia_gzs/backend/internal/domain"
	"github.com/kevinenv/barberia_gzs/backend/internal/repository"
	"github.com/kevinenv/barberia_gzs/backend/internal/service"
)

type AuthHandler struct {
	authService *service.AuthService
	usuarioRepo *repository.UsuarioRepository
}

func NewAuthHandler(authService *service.AuthService, usuarioRepo *repository.UsuarioRepository) *AuthHandler {
	return &AuthHandler{authService: authService, usuarioRepo: usuarioRepo}
}

func (h *AuthHandler) Registrar(c *gin.Context) {
	var req domain.RegistroRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos: revisa nombre, teléfono, email y contraseña (mín. 6 caracteres)"})
		return
	}

	resp, err := h.authService.RegistrarCliente(req)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req domain.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Debes ingresar correo y contraseña"})
		return
	}

	resp, err := h.authService.Login(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// MiPerfil es una ruta protegida: devuelve los datos del usuario logueado
func (h *AuthHandler) MiPerfil(c *gin.Context) {
	usuarioID := c.GetInt("usuario_id")
	usuario, err := h.usuarioRepo.ObtenerPorID(usuarioID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
		return
	}
	c.JSON(http.StatusOK, usuario)
}
