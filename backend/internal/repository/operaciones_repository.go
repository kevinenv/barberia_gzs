package repository

import (
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/kevinenv/barberia_gzs/backend/internal/domain"
)

type OperacionesRepository struct {
	db *sqlx.DB
}

func NewOperacionesRepository(db *sqlx.DB) *OperacionesRepository {
	return &OperacionesRepository{db: db}
}

// CrearSolicitudMantenimiento guarda el formulario de falla de una silla hidráulica
func (r *OperacionesRepository) CrearSolicitudMantenimiento(clienteID int, req domain.CrearMantenimientoRequest) (*domain.SolicitudMantenimiento, error) {
	query := `INSERT INTO solicitudes_mantenimiento (cliente_id, marca_modelo_silla, descripcion_falla, estado)
	          VALUES (?, ?, ?, 'pendiente_revision')`

	res, err := r.db.Exec(query, clienteID, req.MarcaModeloSilla, req.DescripcionFalla)
	if err != nil {
		return nil, err
	}

	id, _ := res.LastInsertId()
	marca := req.MarcaModeloSilla
	return &domain.SolicitudMantenimiento{
		ID:               int(id),
		ClienteID:        clienteID,
		MarcaModeloSilla: &marca,
		DescripcionFalla: req.DescripcionFalla,
		Estado:           "pendiente_revision",
	}, nil
}

// ObtenerMisSolicitudesMantenimiento lista las sillas hidráulicas enviadas por el cliente
func (r *OperacionesRepository) ObtenerMisSolicitudesMantenimiento(clienteID int) ([]domain.SolicitudMantenimiento, error) {
	lista := []domain.SolicitudMantenimiento{}
	query := `SELECT id, cliente_id, marca_modelo_silla, descripcion_falla, cotizacion_estimada,
	                 DATE_FORMAT(fecha_programada, '%Y-%m-%d') AS fecha_programada, estado, created_at
	          FROM solicitudes_mantenimiento
	          WHERE cliente_id = ?
	          ORDER BY created_at DESC`
	err := r.db.Select(&lista, query, clienteID)
	return lista, err
}

// RegistrarVentaProducto descuenta el stock del inventario y registra la venta en caja (ACID)
func (r *OperacionesRepository) RegistrarVentaProducto(usuarioID int, req domain.VentaProductoRequest) (float64, int, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return 0, 0, err
	}
	defer tx.Rollback()

	var prod struct {
		PrecioVenta float64 `db:"precio_venta"`
		StockActual int     `db:"stock_actual"`
	}

	// Bloqueamos la fila del producto mientras verificamos el stock disponible
	err = tx.Get(&prod, `SELECT precio_venta, stock_actual FROM productos WHERE id = ? AND activo = TRUE FOR UPDATE`, req.ProductoID)
	if err != nil {
		return 0, 0, errors.New("producto no encontrado")
	}

	if prod.StockActual < req.Cantidad {
		return 0, 0, errors.New("stock insuficiente para realizar la venta")
	}

	subtotal := prod.PrecioVenta * float64(req.Cantidad)
	nuevoStock := prod.StockActual - req.Cantidad

	// 1. Descontar del inventario
	_, err = tx.Exec(`UPDATE productos SET stock_actual = ? WHERE id = ?`, nuevoStock, req.ProductoID)
	if err != nil {
		return 0, 0, err
	}

	// 2. Registrar la venta en mostrador
	_, err = tx.Exec(`INSERT INTO ventas_productos (producto_id, usuario_id, cantidad, precio_unitario, subtotal, metodo_pago)
	                  VALUES (?, ?, ?, ?, ?, ?)`,
		req.ProductoID, usuarioID, req.Cantidad, prod.PrecioVenta, subtotal, req.MetodoPago)
	if err != nil {
		return 0, 0, err
	}

	return subtotal, nuevoStock, tx.Commit()
}

// ObtenerReporteComisiones calcula el 50% para el barbero y el 50% para la barbería
func (r *OperacionesRepository) ObtenerReporteComisiones() ([]domain.ComisionBarberoResumen, error) {
	reporte := []domain.ComisionBarberoResumen{}
	query := `SELECT b.id AS barbero_id, b.apodo, b.silla_numero,
	                 COUNT(cs.id) AS total_servicios,
	                 COALESCE(SUM(cs.precio_aplicado), 0) AS total_generado,
	                 COALESCE(SUM(cs.comision_barbero), 0) AS comision_ganada,
	                 COALESCE(SUM(cs.precio_aplicado - cs.comision_barbero), 0) AS ganancia_barberia
	          FROM barberos b
	          LEFT JOIN citas c ON c.barbero_id = b.id AND c.estado IN ('confirmada', 'completada')
	          LEFT JOIN cita_servicios cs ON cs.cita_id = c.id
	          WHERE b.activo = TRUE
	          GROUP BY b.id, b.apodo, b.silla_numero
	          ORDER BY b.silla_numero ASC`
	err := r.db.Select(&reporte, query)
	return reporte, err
}
