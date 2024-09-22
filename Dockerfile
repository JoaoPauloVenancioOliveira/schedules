# Usar uma imagem base do Go
FROM golang:1.23.1 AS builder

# Definir o diretório de trabalho
WORKDIR /app

# Copiar o go.mod e go.sum para o diretório de trabalho
COPY go.mod go.sum ./

# Baixar as dependências
RUN go mod download

# Copiar o código da aplicação
COPY . .

# Compilar a aplicação
RUN go build -o schedules schedules.go

# Usar uma imagem mais leve para rodar a aplicação
FROM gcr.io/distroless/base

# Definir o diretório de trabalho
WORKDIR /app

# Copiar o executável da imagem builder
COPY --from=builder /app/schedules .


# Expor a porta da aplicação
EXPOSE 8080

# Comando para rodar a aplicação
CMD ["./schedule"]