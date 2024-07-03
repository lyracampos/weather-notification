FROM golang:1.22.2-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o ./weather-notification ./cmd

FROM alpine
COPY --from=build /app/weather-notification /usr/local/bin/app
COPY --from=build /app/configs /configs

ENTRYPOINT ["app", "-c", "/configs/config.yaml"]

EXPOSE 8080
EXPOSE 9090