# BUILDER
FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.* ./

RUN go mod download

COPY . .

RUN go build -o main main.go

# RUNNER
FROM alpine:latest AS runner

WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 3000

CMD [ "./main", "--production" ]

# build command: docker build -t go-bank:latest .
# run command: docker run -p 3000:3000 -e DATABASE_URL="postgres://postgres:gobank@host.docker.internal:5432/postgres?sslmode=disable" -e JWT_SECRET="secret" go-bank:latest