FROM golang:1.27-alpine@sha256:eb61c3873c97416fff0116553fff65d67a3549f72e11d9758849a689cdd8b523 AS builder

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
