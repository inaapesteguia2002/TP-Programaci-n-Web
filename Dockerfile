# Etapa de compilación
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copiar dependencias
COPY go.mod go.sum ./
RUN go mod download

# Copiar el código fuente
COPY . .

# Compilar el binario ejecutable
RUN CGO_ENABLED=0 GOOS=linux go build -o /de-rabona-app .

# Etapa final de ejecución
FROM alpine:latest

WORKDIR /root/

# Copiar el binario y la carpeta estática necesaria
COPY --from=builder /de-rabona-app .
COPY --from=builder /app/static ./static

EXPOSE 8080

CMD ["./de-rabona-app"]