FROM golang:1.23-alpine

WORKDIR /app

RUN apk add --no-cache git bash postgresql-client

RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.30.0
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.18.3

COPY go.mod ./
COPY go.sum* ./

RUN go mod download

COPY . .

EXPOSE 8080

CMD ["go", "run", "./cmd/api"]