FROM golang:1.21-alpine AS build

WORKDIR /app

COPY . .

RUN go build -o ./bin/connect4 .

FROM alpine:latest

WORKDIR /app

COPY --from=build /app/bin/connect4 .

CMD [ "./discord-bot" ]