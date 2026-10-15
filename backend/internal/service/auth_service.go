package service

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kevinenv/barberia_gzs/backend/internal/domain"
	"github.com/kevinenv/barberia_gzs/backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo *repository.UsuarioRepository
}

func NewAuthService(repo *repository.UsuarioRepository) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) RegistrarCliente(req domain.RegistroRequest) (*domain.AuthResponse, error) {
	// 1. Encriptar la contraseña con Bcrypt (costo 10 estándar de producción)
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("error al procesar la contraseña")
	}

	nuevoUsuario := &domain.Usuario{
		Nombre:       req.Nombre,
		Telefono:     req.Telefono,
		Email:        req.Email,
		PasswordHash: string(hashBytes),
		Rol:          "cliente",
		Activo:       true,
	}

	// 2. Guardar en MySQL
	if err := s.repo.Crear(nuevoUsuario); err != nil {
		return nil, errors.New("el correo o número de teléfono ya está registrado")
	}

	// 3. Generar su Token JWT automáticamente para que quede logueado
	token, err := s.generarJWT(nuevoUsuario.ID, nuevoUsuario.Rol)
	if err != nil {
		return nil, err
	}

	return &domain.AuthResponse{
		Token:   token,
		Usuario: *nuevoUsuario,
	}, nil
}

func (s *AuthService) Login(req domain.LoginRequest) (*domain.AuthResponse, error) {
	usuario, err := s.repo.ObtenerPorEmail(req.Email)
	if err != nil || !usuario.Activo {
		return nil, errors.New("credenciales incorrectas")
	}

	// Comparar contraseña ingresada con el hash guardado en MySQL
	if err := bcrypt.CompareHashAndPassword([]byte(usuario.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.New("credenciales incorrectas")
	}

	token, err := s.generarJWT(usuario.ID, usuario.Rol)
	if err != nil {
		return nil, err
	}

	return &domain.AuthResponse{
		Token:   token,
		Usuario: *usuario,
	}, nil
}

func (s *AuthService) generarJWT(usuarioID int, rol string) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	claims := jwt.MapClaims{
		"usuario_id": usuarioID,
		"rol":        rol,
		"exp":        time.Now().Add(24 * time.Hour * 7).Unix(), // Dura 7 días
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
