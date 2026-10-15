package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/kevinenv/barberia_gzs/backend/internal/domain"
)

type CatalogoRepository struct {
	db *sqlx.DB
}

// NewCatalogoRepository es el "constructor" en Go
func NewCatalogoRepository(db *sqlx.DB) *CatalogoRepository {
	return &CatalogoRepository{db: db}
}

// ObtenerConfiguracion trae los horarios y porcentaje de garantía de la barbería
func (r *CatalogoRepository) ObtenerConfiguracion() (*domain.ConfiguracionBarberia, error) {
	var config domain.ConfiguracionBarberia
	query := `SELECT id, nombre_negocio, hora_apertura, hora_cierre, 
	                 hora_inicio_descanso, hora_fin_descanso, porcentaje_garantia_reserva, 
	                 direccion, telefono_contacto, updated_at 
	          FROM configuracion_barberia WHERE id = 1`

	err := r.db.Get(&config, query)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// ObtenerServiciosActivos lista todos los cortes y servicios disponibles
func (r *CatalogoRepository) ObtenerServiciosActivos() ([]domain.Servicio, error) {
	servicios := []domain.Servicio{} // Inicializamos slice vacío para que JSON devuelva [] y no null
	query := `SELECT id, nombre, descripcion, precio, duracion_minutos, activo, created_at 
	          FROM servicios 
	          WHERE activo = TRUE 
	          ORDER BY id ASC`

	err := r.db.Select(&servicios, query)
	if err != nil {
		return nil, err
	}
	return servicios, nil
}

// ObtenerProductosActivos lista las ceras, pomadas, shampoos y bebidas
func (r *CatalogoRepository) ObtenerProductosActivos() ([]domain.Producto, error) {
	productos := []domain.Producto{}
	query := `SELECT id, nombre, categoria, descripcion, precio_venta, stock_actual, stock_minimo, activo, created_at, updated_at 
	          FROM productos 
	          WHERE activo = TRUE 
	          ORDER BY categoria, nombre ASC`

	err := r.db.Select(&productos, query)
	if err != nil {
		return nil, err
	}
	return productos, nil
}
