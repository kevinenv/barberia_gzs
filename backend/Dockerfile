# ETAPA 1: Compilar el binario de Go
FROM golang:alpine AS builder
WORKDIR /app

# Descargar dependencias primero para aprovechar la caché de Docker
COPY go.mod go.sum ./
RUN go mod download

# Copiar el código fuente y compilar un binario estático ultraligero
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o barberia-api ./cmd/api/main.go

# ETAPA 2: Imagen final de Producción (Ultraligera y segura)
FROM alpine:3.20
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app

# Copiamos únicamente el binario compilado desde la Etapa 1
COPY --from=builder /app/barberia-api .

EXPOSE 8080
CMD ["./barberia-api"]