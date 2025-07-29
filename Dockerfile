# Etapa 1: build
FROM golang:1.24 AS builder

WORKDIR /app

# Copiar los archivos necesarios
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Compilar el binario
RUN CGO_ENABLED=0 go build -o app ./cmd/main.go

# Etapa 2: imagen final
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copiar solo el binario desde la etapa anterior
COPY --from=builder /app/app .

# Ejecutar la app por defecto
ENTRYPOINT ["./app"]
