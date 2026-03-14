FROM golang:1.26-alpine

EXPOSE 3000

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

ADD . .

RUN go build -o /etuedt cmd/etuedt.go

CMD [ "/etuedt" ]
