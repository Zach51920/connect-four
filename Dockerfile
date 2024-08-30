FROM golang:1.22-alpine AS build

WORKDIR /app

COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/connect4 .

FROM alpine:latest

WORKDIR /app

COPY --from=build /app/bin/connect4 .
COPY --from=build /app/configs /app/configs
COPY --from=build /app/public /app/public

ENV CONFIG_PATH=/app/configs/config.prod.yaml

CMD ["./connect4"]
