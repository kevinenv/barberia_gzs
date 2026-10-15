package domain

import "time"

// SolicitudMantenimiento representa el formulario para reparar sillas hidráulicas
type SolicitudMantenimiento struct {
	ID                 int       `db:"id" json:"id"`
	ClienteID          int       `db:"cliente_id" json:"cliente_id"`
	MarcaModeloSilla   *string   `db:"marca_modelo_silla" json:"marca_modelo_silla"`
	DescripcionFalla   string    `db:"descripcion_falla" json:"descripcion_falla"`
	CotizacionEstimada *float64  `db:"cotizacion_estimada" json:"cotizacion_estimada"`
	FechaProgramada    *string   `db:"fecha_programada" json:"fecha_programada"`
	Estado             string    `db:"estado" json:"estado"`
	CreatedAt          time.Time `db:"created_at" json:"created_at"`
}

type CrearMantenimientoRequest struct {
	MarcaModeloSilla string `json:"marca_modelo_silla" binding:"required"`
	DescripcionFalla string `json:"descripcion_falla" binding:"required,min=10"`
}

// VentaProductoRequest es el DTO para vender ceras, pomadas, shampoos o bebidas descontando stock
type VentaProductoRequest struct {
	ProductoID int    `json:"producto_id" binding:"required"`
	Cantidad   int    `json:"cantidad" binding:"required,min=1"`
	MetodoPago string `json:"metodo_pago" binding:"required"` // "efectivo", "qr", "tarjeta"
}

// ComisionBarberoResumen muestra cuánto produjo cada barbero y su ganancia del 50%
type ComisionBarberoResumen struct {
	BarberoID          int     `db:"barbero_id" json:"barbero_id"`
	Apodo              string  `db:"apodo" json:"apodo"`
	SillaNumero        int     `db:"silla_numero" json:"silla_numero"`
	TotalServicios     int     `db:"total_servicios" json:"total_servicios"`
	TotalGenerado      float64 `db:"total_generado" json:"total_generado"`
	ComisionGanada50   float64 `db:"comision_ganada" json:"comision_ganada"`
	GananciaBarberia50 float64 `db:"ganancia_barberia" json:"ganancia_barberia"`
}
