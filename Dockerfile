# Etapa 1: Builder (Compilação)
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Instala ferramentas básicas necessárias
RUN apk add --no-cache git

# Baixa as dependências (Cache Layer)
COPY go.mod go.sum ./
RUN go mod download

# Copia o código fonte
COPY . .

# Compila o binário estático
# CGO_ENABLED=0 garante que não dependa de bibliotecas C do sistema
RUN CGO_ENABLED=0 GOOS=linux go build -o hackaton-service-api cmd/main.go

# Etapa 2: Runner (Imagem Final Leve)
FROM alpine:latest

WORKDIR /app

# Instala certificados CA (necessário para AWS/HTTPS)
RUN apk --no-cache add ca-certificates

# Copia o binário da etapa anterior
COPY --from=builder /app/hackaton-service-api .

# IMPORTANTE: Copia os arquivos de frontend (HTML/CSS)
COPY --from=builder /app/web ./web

# Expondo a porta padrão da API
EXPOSE 8080

# Comando de inicialização
CMD ["./hackaton-service-api"]