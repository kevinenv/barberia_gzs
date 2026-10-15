package domain

import "time"

// ConfiguracionBarberia representa los horarios generales y garantía del negocio
type ConfiguracionBarberia struct {
	ID                        int       `db:"id" json:"id"`
	NombreNegocio             string    `db:"nombre_negocio" json:"nombre_negocio"`
	HoraApertura              string    `db:"hora_apertura" json:"hora_apertura"`
	HoraCierre                string    `db:"hora_cierre" json:"hora_cierre"`
	HoraInicioDescanso        string    `db:"hora_inicio_descanso" json:"hora_inicio_descanso"`
	HoraFinDescanso           string    `db:"hora_fin_descanso" json:"hora_fin_descanso"`
	PorcentajeGarantiaReserva float64   `db:"porcentaje_garantia_reserva" json:"porcentaje_garantia_reserva"`
	Direccion                 *string   `db:"direccion" json:"direccion"` // *string permite valores NULL de MySQL
	TelefonoContacto          *string   `db:"telefono_contacto" json:"telefono_contacto"`
	UpdatedAt                 time.Time `db:"updated_at" json:"updated_at"`
}

// Servicio representa un corte, arreglo de barba o perfilado de cejas
type Servicio struct {
	ID              int       `db:"id" json:"id"`
	Nombre          string    `db:"nombre" json:"nombre"`
	Descripcion     *string   `db:"descripcion" json:"descripcion"`
	Precio          float64   `db:"precio" json:"precio"`
	DuracionMinutos int       `db:"duracion_minutos" json:"duracion_minutos"`
	Activo          bool      `db:"activo" json:"activo"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
}

// Producto representa ceras, pomadas, shampoo o bebidas en inventario
type Producto struct {
	ID          int       `db:"id" json:"id"`
	Nombre      string    `db:"nombre" json:"nombre"`
	Categoria   string    `db:"categoria" json:"categoria"`
	Descripcion *string   `db:"descripcion" json:"descripcion"`
	PrecioVenta float64   `db:"precio_venta" json:"precio_venta"`
	StockActual int       `db:"stock_actual" json:"stock_actual"`
	StockMinimo int       `db:"stock_minimo" json:"stock_minimo"`
	Activo      bool      `db:"activo" json:"activo"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}
