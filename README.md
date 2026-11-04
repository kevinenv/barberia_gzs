# Barbería GZS

Aplicación web para la gestión de una barbería. El proyecto incluye una interfaz web, una API y una base de datos MySQL. Permite consultar servicios y productos, gestionar reservas y acceder a funciones de usuario, inventario y operaciones.

## Tecnologías

- **Frontend:** Vue 3, Vite, Pinia y Vue Router.
- **Backend:** Go, Gin y sqlx.
- **Base de datos:** MySQL 8.
- **Contenedores:** Docker y Docker Compose.

## Requisitos

Para ejecutar el proyecto en desarrollo, instala:

- Git.
- Docker y Docker Compose.
- Go, en la versión indicada en `backend/go.mod` (Go 1.27.2).
- Node.js 24 y npm, para ejecutar el frontend localmente.

## Configuración inicial

Clona el repositorio y entra en el directorio del proyecto:

```bash
git clone https://github.com/kevinenv/barberia_gzs.git
cd barberia_gzs
```

Crea el archivo de variables de entorno a partir de la plantilla:

```bash
cp .env.example .env
```

Edita `.env` y sustituye los valores de ejemplo por credenciales propias. Como mínimo, configura `DB_ROOT_PASSWORD`, `DB_NAME`, `DB_USER`, `DB_PASSWORD` y `JWT_SECRET`. Conserva `DB_PORT` y `PORT` si los puertos predeterminados están disponibles.

No compartas ni subas el archivo `.env` a repositorios públicos: contiene credenciales y secretos de la aplicación.

## Ejecución completa con Docker

Desde la raíz del repositorio, ejecuta:

```bash
docker compose -f docker-compose.prod.yml up -d --build
```

Este comando inicia MySQL, la API y el frontend. Cuando los servicios estén disponibles, abre [http://localhost](http://localhost) en el navegador.

Para consultar el estado de los servicios y sus registros:

```bash
docker compose -f docker-compose.prod.yml ps
docker compose -f docker-compose.prod.yml logs -f
```

Para detener los servicios:

```bash
docker compose -f docker-compose.prod.yml down
```

La migración inicial se aplica automáticamente cuando MySQL crea una base de datos nueva. Si ya existe un volumen de datos, Docker no vuelve a ejecutar automáticamente esa migración. El comando `down` conserva los datos; para eliminarlos también, se puede usar `docker compose -f docker-compose.prod.yml down -v`. Esta última opción borra permanentemente el volumen de la base de datos.

## Ejecución en desarrollo

El modo de desarrollo ejecuta MySQL en Docker y permite iniciar la API y el frontend por separado.

### 1. Iniciar la base de datos

Desde la raíz del repositorio, inicia el servicio de base de datos:

```bash
docker compose -f docker-compose.dev.yml up -d db
```

Aplica el esquema inicial a la base de datos:

```bash
docker compose -f docker-compose.dev.yml exec -T db sh -c 'mysql -uroot -p"$MYSQL_ROOT_PASSWORD" "$MYSQL_DATABASE"' < backend/migrations/001_init_schema.sql
```

### 2. Iniciar la API

En una terminal, ejecuta:

```bash
cd backend
go mod download
go run ./cmd/api
```

La API estará disponible en `http://localhost:8080`. Puedes comprobar su estado con esta dirección: [http://localhost:8080/api/v1/health](http://localhost:8080/api/v1/health).

### 3. Iniciar el frontend

En otra terminal, desde la raíz del repositorio, ejecuta:

```bash
cd frontend
npm ci
npm run dev
```

Abre la dirección que Vite muestre en la terminal, normalmente [http://localhost:5173](http://localhost:5173). La configuración de desarrollo dirige las solicitudes de la aplicación a `http://localhost:8080/api/v1`.

Para detener el frontend y la API, pulsa `Ctrl+C` en sus respectivas terminales. Para detener MySQL:

```bash
docker compose -f docker-compose.dev.yml down
```

## Estructura del proyecto

```text
.
├── backend/                 API y migraciones de base de datos
├── frontend/                Aplicación web
├── docker-compose.dev.yml   Configuración de desarrollo para MySQL
├── docker-compose.prod.yml  Configuración de ejecución con Docker
└── .env.example             Plantilla de variables de entorno
```

## Solución de problemas

- **No se puede conectar a MySQL:** comprueba que el contenedor esté activo con `docker compose -f docker-compose.dev.yml ps` y que `DB_PORT` esté disponible.
- **La API no inicia:** verifica las credenciales de la base de datos en `.env`, que MySQL esté iniciado y que el esquema inicial se haya aplicado.
- **El frontend no conecta con la API:** confirma que la API esté escuchando en el puerto `8080` y que el archivo `frontend/.env.development` apunte a `http://localhost:8080/api/v1`.
- **El puerto ya está en uso:** libera el puerto ocupado o actualiza la configuración correspondiente antes de iniciar los servicios.
