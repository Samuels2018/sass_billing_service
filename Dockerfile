# Etapa de construcción
FROM golang:1.21-alpine AS builder

# Instalar herramientas necesarias
RUN apk add --no-cache git ca-certificates

# Configurar el directorio de trabajo
WORKDIR /app

# Copiar los archivos de módulos primero para aprovechar el cache de Docker
COPY go.mod go.sum ./

# Descargar dependencias
RUN go mod download

# Copiar el resto del código fuente
COPY . .

# Construir la aplicación
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /sass-billing-service

# Etapa final (imagen ligera)
FROM alpine:3.18

# Instalar ca-certificates para conexiones TLS
RUN apk add --no-cache ca-certificates

# Copiar el binario desde el builder
COPY --from=builder /sass-billing-service /sass-billing-service

# Puerto expuesto
EXPOSE 8080

# Comando de ejecución
CMD ["/sass-billing-service"]