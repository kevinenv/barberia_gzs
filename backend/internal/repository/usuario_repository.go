package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/kevinenv/barberia_gzs/backend/internal/domain"
)

type UsuarioRepository struct {
	db *sqlx.DB
}

func NewUsuarioRepository(db *sqlx.DB) *UsuarioRepository {
	return &UsuarioRepository{db: db}
}

func (r *UsuarioRepository) Crear(u *domain.Usuario) error {
	query := `INSERT INTO usuarios (nombre, telefono, email, password_hash, rol) 
	          VALUES (?, ?, ?, ?, ?)`

	result, err := r.db.Exec(query, u.Nombre, u.Telefono, u.Email, u.PasswordHash, u.Rol)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	u.ID = int(id)
	return nil
}

func (r *UsuarioRepository) ObtenerPorEmail(email string) (*domain.Usuario, error) {
	var u domain.Usuario
	query := `SELECT id, nombre, telefono, email, password_hash, rol, activo, created_at, updated_at 
	          FROM usuarios WHERE email = ? LIMIT 1`

	err := r.db.Get(&u, query, email)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UsuarioRepository) ObtenerPorID(id int) (*domain.Usuario, error) {
	var u domain.Usuario
	query := `SELECT id, nombre, telefono, email, password_hash, rol, activo, created_at, updated_at 
	          FROM usuarios WHERE id = ? LIMIT 1`

	err := r.db.Get(&u, query, id)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
