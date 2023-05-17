# Build Stage
# First pull Golang image
FROM golang:1.19-alpine as build-env

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 go build /app/cmd/app.go

CMD ["./app"]

EXPOSE 8000