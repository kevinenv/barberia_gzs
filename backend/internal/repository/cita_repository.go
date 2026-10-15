package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/kevinenv/barberia_gzs/backend/internal/domain"
)

type CitaRepository struct {
	db *sqlx.DB
}

func NewCitaRepository(db *sqlx.DB) *CitaRepository {
	return &CitaRepository{db: db}
}

// ObtenerBarberosActivos devuelve todos los sillones disponibles
func (r *CitaRepository) ObtenerBarberosActivos() ([]domain.Barbero, error) {
	barberos := []domain.Barbero{}
	query := `SELECT id, usuario_id, apodo, especialidad, silla_numero, comision_porcentaje, activo 
	          FROM barberos WHERE activo = TRUE ORDER BY silla_numero ASC`
	err := r.db.Select(&barberos, query)
	return barberos, err
}

// ObtenerServiciosPorIDs busca los precios y duraciones de los servicios seleccionados
func (r *CitaRepository) ObtenerServiciosPorIDs(ids []int) ([]domain.Servicio, error) {
	query, args, err := sqlx.In(`SELECT id, nombre, precio, duracion_minutos FROM servicios WHERE id IN (?) AND activo = TRUE`, ids)
	if err != nil {
		return nil, err
	}
	query = r.db.Rebind(query)
	var servicios []domain.Servicio
	err = r.db.Select(&servicios, query, args...)
	return servicios, err
}

// BarberoTieneCruce verifica si el barbero ya tiene otra cita que se solape en ese horario
func (r *CitaRepository) BarberoTieneCruce(barberoID int, fecha string, horaInicio string, horaFin string) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM citas 
	          WHERE barbero_id = ? 
	            AND fecha = ? 
	            AND estado IN ('pendiente_pago', 'confirmada', 'completada')
	            AND (? < hora_fin AND ? > hora_inicio)`
	err := r.db.Get(&count, query, barberoID, fecha, horaInicio, horaFin)
	if err != nil {
		return true, err
	}
	return count > 0, nil
}

// GuardarCitaCompleta usa una Transacción ACID para guardar la cabecera y todos sus servicios
func (r *CitaRepository) GuardarCitaCompleta(cita *domain.Cita, detalles []domain.CitaServicio) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	queryCita := `INSERT INTO citas (cliente_id, barbero_id, fecha, hora_inicio, hora_fin, origen, monto_total, monto_garantia_pagado, metodo_pago_garantia, estado)
	              VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	res, err := tx.Exec(queryCita, cita.ClienteID, cita.BarberoID, cita.Fecha, cita.HoraInicio, cita.HoraFin,
		cita.Origen, cita.MontoTotal, cita.MontoGarantiaPagado, cita.MetodoPagoGarantia, cita.Estado)
	if err != nil {
		return err
	}

	citaID, err := res.LastInsertId()
	if err != nil {
		return err
	}
	cita.ID = int(citaID)

	queryDetalle := `INSERT INTO cita_servicios (cita_id, servicio_id, precio_aplicado, duracion_minutos, comision_barbero)
	                 VALUES (?, ?, ?, ?, ?)`

	for _, d := range detalles {
		_, err := tx.Exec(queryDetalle, cita.ID, d.ServicioID, d.PrecioAplicado, d.DuracionMinutos, d.ComisionBarbero)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// ObtenerMisCitas trae el historial de reservas del cliente logueado
func (r *CitaRepository) ObtenerMisCitas(clienteID int) ([]domain.Cita, error) {
	citas := []domain.Cita{}
	query := `SELECT c.id, c.cliente_id, c.barbero_id, b.apodo AS barbero_apodo, b.silla_numero,
	                 DATE_FORMAT(c.fecha, '%Y-%m-%d') AS fecha, c.hora_inicio, c.hora_fin, 
	                 c.origen, c.monto_total, c.monto_garantia_pagado, c.metodo_pago_garantia, c.estado, c.created_at
	          FROM citas c
	          INNER JOIN barberos b ON c.barbero_id = b.id
	          WHERE c.cliente_id = ?
	          ORDER BY c.fecha DESC, c.hora_inicio DESC`
	err := r.db.Select(&citas, query, clienteID)
	return citas, err
}
