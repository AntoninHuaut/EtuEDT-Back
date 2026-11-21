FROM golang:1.27-alpine@sha256:70bd7ddea8b6cbe5d6943cde279d6b67feb10cc924fb9100d1d8757c1575b280 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /etuedt ./cmd/etuedt.go

FROM alpine:3.24@sha256:f59abe656b692d901e184be3445e084786602108599d2a51d229c0f5a48c71bf

RUN apk add --no-cache ca-certificates && \
    addgroup -S etuedt && adduser -S etuedt -G etuedt

EXPOSE 3000

WORKDIR /app

COPY --from=builder /etuedt /etuedt
COPY --from=builder /app/config.json /app/config.json

USER etuedt

CMD [ "/etuedt" ]
