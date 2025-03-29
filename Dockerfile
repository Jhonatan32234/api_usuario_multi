# Usar la imagen base de Go
FROM golang:1.23-alpine AS builder

# Establecer el directorio de trabajo dentro del contenedor
WORKDIR /app

# Copiar go.mod y go.sum al contenedor
COPY go.mod go.sum ./

# Descargar las dependencias
RUN go mod tidy

# Copiar todo el código fuente del proyecto al contenedor
COPY . .

# Construir la aplicación
RUN go build -o apiuser main.go

# Crear una imagen mínima para producción
FROM alpine:latest

# Establecer el directorio de trabajo en el contenedor final
WORKDIR /app

# Copiar el binario desde la fase de construcción
COPY --from=builder /app/apiuser .

# Exponer el puerto 5000
EXPOSE 5000

# Ejecutar la aplicación
CMD ["/app/apiuser"]
