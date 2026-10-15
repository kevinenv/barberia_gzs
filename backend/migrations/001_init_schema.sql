-- ============================================================================
-- MIGRACIÓN 001: Esquema Oficial de Producción - Barbería GZS
-- ============================================================================

-- 1. Configuración Global de la Barbería (Mismo horario y descanso para todos)
CREATE TABLE IF NOT EXISTS configuracion_barberia (
    id INT PRIMARY KEY DEFAULT 1,
    nombre_negocio VARCHAR(100) NOT NULL DEFAULT 'Barbería GZS',
    hora_apertura TIME NOT NULL DEFAULT '09:00:00',
    hora_cierre TIME NOT NULL DEFAULT '20:00:00',
    hora_inicio_descanso TIME NOT NULL DEFAULT '13:00:00',
    hora_fin_descanso TIME NOT NULL DEFAULT '14:30:00',
    porcentaje_garantia_reserva DECIMAL(5,2) NOT NULL DEFAULT 50.00, -- Adelanto obligatorio para confirmar cita
    direccion VARCHAR(255) NULL,
    telefono_contacto VARCHAR(20) NULL,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 2. Usuarios (Registro obligatorio para historial y seguridad contra citas falsas)
CREATE TABLE IF NOT EXISTS usuarios (
    id INT AUTO_INCREMENT PRIMARY KEY,
    nombre VARCHAR(100) NOT NULL,
    telefono VARCHAR(20) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    rol ENUM('admin', 'barbero', 'cliente') NOT NULL DEFAULT 'cliente',
    activo BOOLEAN NOT NULL DEFAULT TRUE,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 3. Barberos (Comisión del 50% sobre el coste del servicio)
CREATE TABLE IF NOT EXISTS barberos (
    id INT AUTO_INCREMENT PRIMARY KEY,
    usuario_id INT NOT NULL UNIQUE,
    apodo VARCHAR(50) NOT NULL,
    especialidad VARCHAR(100) NOT NULL,
    silla_numero INT NOT NULL UNIQUE,
    comision_porcentaje DECIMAL(5,2) NOT NULL DEFAULT 50.00,
    activo BOOLEAN NOT NULL DEFAULT TRUE,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_barbero_usuario FOREIGN KEY (usuario_id) REFERENCES usuarios(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 4. Catálogo de Servicios de Barbería (Cortes clásicos, actuales, cejas, barba)
CREATE TABLE IF NOT EXISTS servicios (
    id INT AUTO_INCREMENT PRIMARY KEY,
    nombre VARCHAR(100) NOT NULL,
    descripcion VARCHAR(255) NULL,
    precio DECIMAL(10,2) NOT NULL,
    duracion_minutos INT NOT NULL,
    activo BOOLEAN NOT NULL DEFAULT TRUE,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 5. Cabecera de Citas (Soporta Modelo Híbrido y Pago de Garantía)
CREATE TABLE IF NOT EXISTS citas (
    id INT AUTO_INCREMENT PRIMARY KEY,
    cliente_id INT NOT NULL,
    barbero_id INT NOT NULL,
    fecha DATE NOT NULL,
    hora_inicio TIME NOT NULL,
    hora_fin TIME NOT NULL, -- Calculada sumando la duración de todos los servicios elegidos
    origen ENUM('web', 'local_walkin') NOT NULL DEFAULT 'web',
    monto_total DECIMAL(10,2) NOT NULL,
    monto_garantia_pagado DECIMAL(10,2) NOT NULL DEFAULT 0.00,
    metodo_pago_garantia ENUM('qr', 'efectivo', 'transferencia') NOT NULL DEFAULT 'qr',
    comprobante_url VARCHAR(255) NULL,
    estado ENUM('pendiente_pago', 'confirmada', 'completada', 'cancelada', 'no_asistio') NOT NULL DEFAULT 'pendiente_pago',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_cita_cliente FOREIGN KEY (cliente_id) REFERENCES usuarios(id),
    CONSTRAINT fk_cita_barbero FOREIGN KEY (barbero_id) REFERENCES barberos(id),
    INDEX idx_agenda_barbero (barbero_id, fecha, estado)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 6. Detalle de Servicios por Cita (Permite elegir varios servicios en la misma reserva)
CREATE TABLE IF NOT EXISTS cita_servicios (
    id INT AUTO_INCREMENT PRIMARY KEY,
    cita_id INT NOT NULL,
    servicio_id INT NOT NULL,
    precio_aplicado DECIMAL(10,2) NOT NULL,
    duracion_minutos INT NOT NULL,
    comision_barbero DECIMAL(10,2) NOT NULL, -- Guarda el 50% exacto para el barbero
    CONSTRAINT fk_detalle_cita FOREIGN KEY (cita_id) REFERENCES citas(id) ON DELETE CASCADE,
    CONSTRAINT fk_detalle_servicio FOREIGN KEY (servicio_id) REFERENCES servicios(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 7. Inventario de Productos (Ceras, pomadas, shampoo, bebidas)
CREATE TABLE IF NOT EXISTS productos (
    id INT AUTO_INCREMENT PRIMARY KEY,
    nombre VARCHAR(100) NOT NULL,
    categoria ENUM('cera_pomada', 'shampoo_cuidado', 'bebida', 'otro') NOT NULL,
    descripcion VARCHAR(255) NULL,
    precio_venta DECIMAL(10,2) NOT NULL,
    stock_actual INT NOT NULL DEFAULT 0,
    stock_minimo INT NOT NULL DEFAULT 5,
    activo BOOLEAN NOT NULL DEFAULT TRUE,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 8. Ventas de Productos en Mostrador (Descuenta del stock)
CREATE TABLE IF NOT EXISTS ventas_productos (
    id INT AUTO_INCREMENT PRIMARY KEY,
    producto_id INT NOT NULL,
    usuario_id INT NULL, -- Quién registró o compró
    cantidad INT NOT NULL,
    precio_unitario DECIMAL(10,2) NOT NULL,
    subtotal DECIMAL(10,2) NOT NULL,
    metodo_pago ENUM('efectivo', 'qr', 'tarjeta') NOT NULL DEFAULT 'efectivo',
    fecha_venta DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_venta_producto FOREIGN KEY (producto_id) REFERENCES productos(id),
    CONSTRAINT fk_venta_usuario FOREIGN KEY (usuario_id) REFERENCES usuarios(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 9. Solicitudes de Mantenimiento de Sillas Hidráulicas (Formulario de cotización)
CREATE TABLE IF NOT EXISTS solicitudes_mantenimiento (
    id INT AUTO_INCREMENT PRIMARY KEY,
    cliente_id INT NOT NULL,
    marca_modelo_silla VARCHAR(100) NULL,
    descripcion_falla TEXT NOT NULL,
    foto_url VARCHAR(255) NULL,
    cotizacion_estimada DECIMAL(10,2) NULL,
    fecha_programada DATE NULL,
    notas_tecnicas TEXT NULL,
    estado ENUM('pendiente_revision', 'cotizado', 'aprobado', 'en_reparacion', 'completado', 'cancelado') NOT NULL DEFAULT 'pendiente_revision',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT fk_mantenimiento_cliente FOREIGN KEY (cliente_id) REFERENCES usuarios(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================================
-- DATOS INICIALES (SEED) DE BARBERÍA GZS
-- ============================================================================

INSERT INTO configuracion_barberia (id, nombre_negocio, hora_apertura, hora_cierre, hora_inicio_descanso, hora_fin_descanso, porcentaje_garantia_reserva)
VALUES (1, 'Barbería GZS', '09:00:00', '20:00:00', '13:00:00', '14:30:00', 50.00)
ON DUPLICATE KEY UPDATE nombre_negocio = VALUES(nombre_negocio);

INSERT INTO servicios (nombre, descripcion, precio, duracion_minutos) VALUES
('Corte Clásico', 'Corte tradicional a tijera y máquina con acabado limpio.', 40.00, 40),
('Corte Actual / Degradado (Fade)', 'Fade bajo, medio o alto, taper o diseño moderno.', 45.00, 45),
('Arreglo y Ritual de Barba', 'Perfilado de barba con navaja, toalla caliente y aceites.', 25.00, 25),
('Perfilado de Cejas', 'Diseño y limpieza de cejas con navaja.', 15.00, 15);

INSERT INTO productos (nombre, categoria, descripcion, precio_venta, stock_actual, stock_minimo) VALUES
('Cera Mate Fijación Fuerte', 'cera_pomada', 'Acabado sin brillo de larga duración 150g.', 35.00, 20, 5),
('Pomada Clásica Efecto Húmedo', 'cera_pomada', 'Brillo medio ideal para peinados clásicos.', 35.00, 15, 5),
('Shampoo Anticaída y Mentol', 'shampoo_cuidado', 'Limpieza profunda para cuero cabelludo 400ml.', 50.00, 10, 3),
('Gaseosa / Agua Fría Personal', 'bebida', 'Bebida fría en botella.', 10.00, 30, 8);