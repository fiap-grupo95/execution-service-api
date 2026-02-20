FROM --platform=$BUILDPLATFORM  golang:1.24.4-alpine AS builder

WORKDIR /app

# Instala swag CLI para gerar docs
RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY . .

RUN go mod download && go clean -modcache

# Gera documentação Swagger (ajuste o caminho do main.go se necessário)
RUN swag init -g cmd/api/main.go -o docs

ARG TARGETOS
ARG TARGETARCH
ARG BUILDPLATFORM

WORKDIR /app/cmd/api

RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH:-amd64} go build -o /execution-service-api

FROM alpine:3.22

WORKDIR /app

COPY --from=builder /execution-service-api .

COPY --from=builder /app/docs ./docs

EXPOSE 8081

CMD ["./execution-service-api"]