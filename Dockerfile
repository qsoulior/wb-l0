FROM golang:1.23.0-alpine3.20 AS dependencies
WORKDIR /dependencies
COPY go.mod go.sum ./
RUN go mod download && go mod verify

FROM dependencies AS build
WORKDIR /build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -v -o ./main ./cmd

FROM alpine:3.20
WORKDIR /app
COPY --from=build /build/main ./
EXPOSE 80
ENTRYPOINT ["./main"]