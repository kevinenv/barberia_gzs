package domain

import "time"

type Usuario struct {
	ID           int       `db:"id" json:"id"`
	Nombre       string    `db:"nombre" json:"nombre"`
	Telefono     string    `db:"telefono" json:"telefono"`
	Email        string    `db:"email" json:"email"`
	PasswordHash string    `db:"password_hash" json:"-"` // El guion "-" oculta este campo en cualquier JSON
	Rol          string    `db:"rol" json:"rol"`
	Activo       bool      `db:"activo" json:"activo"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

// DTO para cuando un cliente se registra en la web
type RegistroRequest struct {
	Nombre   string `json:"nombre" binding:"required,min=3"`
	Telefono string `json:"telefono" binding:"required,min=7"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// DTO para iniciar sesión
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// Respuesta que recibe Vue 3 al iniciar sesión o registrarse
type AuthResponse struct {
	Token   string  `json:"token"`
	Usuario Usuario `json:"usuario"`
}
