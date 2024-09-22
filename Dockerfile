# Usar uma imagem base do Golang
FROM golang:1.23.1 as builder

# Define o diretório de trabalho
WORKDIR /app

# Copia os arquivos go.mod e go.sum
COPY go.mod go.sum ./

# Baixa as dependências
RUN go mod download

# Copia o código fonte
COPY . .

# Compila a aplicação
RUN go build -o main .

# Imagem final
FROM alpine:latest

# Copia o binário da imagem builder
COPY --from=builder /app/main .

# Comando para rodar a aplicação
CMD ["./main"]
