# ---- Etapa 1: build (con el SDK de Go) ----
FROM golang:1.22-alpine AS build
WORKDIR /src

COPY . .

# go mod tidy resuelve e instala las dependencias necesarias para compilar
# y genera el go.sum (no lo traiamos commiteado desde el vamos).
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server .

# ---- Etapa de test: parte de la imagen de build (ya tiene el SDK y el codigo) ----
# No forma parte de la imagen final: solo se construye con --target test (o
# target: test en el workflow). Al correr el contenedor ejecuta los tests, mide
# cobertura y falla si no llega al umbral. Va en el ENTRYPOINT (no en un RUN)
# para poder sacar los reportes al host con un volumen (-v ...:/out).
FROM build AS test
ENTRYPOINT ["sh", "scripts/test.sh"]

# ---- Etapa 2: runtime minimo (sin el SDK) ----
FROM alpine:3.20 AS final
WORKDIR /app
COPY --from=build /app/server .
EXPOSE 8080
ENTRYPOINT ["./server"]