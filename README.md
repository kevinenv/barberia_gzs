# 💈 Barbería GZS — Sistema Full-Stack de Reservas, Inventario y Caja

Sistema integral en tiempo real diseñado para la gestión operativa de una barbería mediana. Incluye motor híbrido de citas con carrito de múltiples servicios, garantía anti-inasistencias por QR, control de stock de mostrador, liquidación de comisiones (50/50) y módulo técnico de mantenimiento de sillas hidráulicas.

## 🛠️ Stack Tecnológico y Arquitectura

- **Backend:** Go (Golang) + Gin Framework + `sqlx` (Arquitectura en Capas: _Domain, Repository, Service, Handler, Middleware_).
- **Seguridad:** Autenticación stateless con **JWT** y encriptación de contraseñas con **bcrypt**.
- **Base de Datos:** MySQL 8.0 con transacciones **ACID** y bloqueo de filas (`FOR UPDATE`) para control de stock y cruces de horario.
- **Frontend:** Vue 3 (Vite + Pinia + Vue Router + Axios Interceptors).
- **Infraestructura:** Docker Multi-Stage Builds, Nginx como _Reverse Proxy_ y CI/CD automatizado con GitHub Actions hacia Docker Hub.

## 🚀 Despliegue Rápido en Producción (Docker)

1. Clonar el repositorio y crear el archivo `.env` basado en `.env.example`.
2. Levantar la infraestructura completa en el puerto 80:
   ```bash
   docker compose -f docker-compose.prod.yml up -d --build
   ```
