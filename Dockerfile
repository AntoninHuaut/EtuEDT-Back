FROM golang:1.24-alpine

EXPOSE 3000

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

ADD . .

RUN go build -o /etuedt-go

CMD [ "/etuedt-go" ]
