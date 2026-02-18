# Etapa 1: Builder (Compilação)
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Instala git para baixar dependências
RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Compila o binário
RUN CGO_ENABLED=0 GOOS=linux go build -o hackaton-service-api cmd/main.go

# Etapa 2: Runner (Imagem Final)
FROM alpine:latest

WORKDIR /app

# [CORREÇÃO] Instala certificados (para AWS/SSL) E tzdata (para TimeZone)
RUN apk --no-cache add ca-certificates tzdata

# Copia o binário
COPY --from=builder /app/hackaton-service-api .

# Copia o frontend
COPY --from=builder /app/web ./web

EXPOSE 8080

CMD ["./hackaton-service-api"]