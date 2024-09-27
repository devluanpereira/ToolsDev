# Usar uma imagem base do Golang
FROM golang:1.23.1-alpine

# Definir o diretório de trabalho dentro do container
WORKDIR /app

# Copiar os arquivos go.mod e go.sum para o container
COPY go.mod go.sum ./

# Baixar as dependências
RUN go mod download

# Copiar o código-fonte e o arquivo .env para o container
COPY . .

# Compilar o aplicativo Go
RUN go build -o main .

# Definir a porta em que o container vai rodar
EXPOSE 8000

# Comando para rodar a aplicação
CMD ["./main"]
