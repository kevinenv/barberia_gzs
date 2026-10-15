package domain

import "time"

// Barbero representa a cada especialista y su sillón en la barbería
type Barbero struct {
	ID                 int     `db:"id" json:"id"`
	UsuarioID          int     `db:"usuario_id" json:"usuario_id"`
	Apodo              string  `db:"apodo" json:"apodo"`
	Especialidad       string  `db:"especialidad" json:"especialidad"`
	SillaNumero        int     `db:"silla_numero" json:"silla_numero"`
	ComisionPorcentaje float64 `db:"comision_porcentaje" json:"comision_porcentaje"`
	Activo             bool    `db:"activo" json:"activo"`
}

// Cita representa la cabecera de la reserva
type Cita struct {
	ID                  int       `db:"id" json:"id"`
	ClienteID           int       `db:"cliente_id" json:"cliente_id"`
	BarberoID           int       `db:"barbero_id" json:"barbero_id"`
	BarberoApodo        string    `db:"barbero_apodo" json:"barbero_apodo,omitempty"`
	SillaNumero         int       `db:"silla_numero" json:"silla_numero,omitempty"`
	Fecha               string    `db:"fecha" json:"fecha"`
	HoraInicio          string    `db:"hora_inicio" json:"hora_inicio"`
	HoraFin             string    `db:"hora_fin" json:"hora_fin"`
	Origen              string    `db:"origen" json:"origen"`
	MontoTotal          float64   `db:"monto_total" json:"monto_total"`
	MontoGarantiaPagado float64   `db:"monto_garantia_pagado" json:"monto_garantia_pagado"`
	MetodoPagoGarantia  string    `db:"metodo_pago_garantia" json:"metodo_pago_garantia"`
	Estado              string    `db:"estado" json:"estado"`
	CreatedAt           time.Time `db:"created_at" json:"created_at"`
}

// CitaServicio representa cada servicio individual dentro del carrito de la cita
type CitaServicio struct {
	ID              int     `db:"id" json:"id"`
	CitaID          int     `db:"cita_id" json:"cita_id"`
	ServicioID      int     `db:"servicio_id" json:"servicio_id"`
	PrecioAplicado  float64 `db:"precio_aplicado" json:"precio_aplicado"`
	DuracionMinutos int     `db:"duracion_minutos" json:"duracion_minutos"`
	ComisionBarbero float64 `db:"comision_barbero" json:"comision_barbero"`
}

// CrearCitaRequest es el JSON que envía el cliente desde Vue 3
type CrearCitaRequest struct {
	BarberoID          int    `json:"barbero_id"`                              // 0 = Cualquier barbero disponible (Híbrido)
	Fecha              string `json:"fecha" binding:"required"`                // Formato: "2026-10-20"
	HoraInicio         string `json:"hora_inicio" binding:"required"`          // Formato: "15:00"
	ServiciosIDs       []int  `json:"servicios_ids" binding:"required,min=1"`  // Ej: [1, 4] (Corte + Cejas)
	MetodoPagoGarantia string `json:"metodo_pago_garantia" binding:"required"` // "qr", "transferencia"
}
