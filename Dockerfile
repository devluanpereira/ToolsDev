# Etapa de build
FROM golang:1.20-alpine AS builder

# Definir diretório de trabalho
WORKDIR /app

# Copiar arquivos go.mod e go.sum e instalar dependências
COPY go.mod go.sum ./
RUN go mod download

# Copiar o código-fonte
COPY . .

# Compilar a aplicação
RUN go build -o main .

# Etapa de execução
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/main .

# Definir a porta usada pela aplicação
EXPOSE 8000

# Comando de execução
CMD ["./main"]
