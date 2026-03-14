FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /etuedt ./cmd/etuedt.go

FROM alpine:3.21

RUN apk add --no-cache ca-certificates && \
    addgroup -S etuedt && adduser -S etuedt -G etuedt

EXPOSE 3000

WORKDIR /app

COPY --from=builder /etuedt /etuedt
COPY --from=builder /app/config.json /app/config.json

USER etuedt

CMD [ "/etuedt" ]
