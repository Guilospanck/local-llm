FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o server

FROM gcr.io/distroless/base-debian12

WORKDIR /root/

COPY --from=builder /app/server .

EXPOSE 4444

CMD ["./server"]

