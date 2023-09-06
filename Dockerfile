FROM golang:1.21-alpine

EXPOSE 3000

WORKDIR /app

COPY go.mod ./
RUN go mod download

ADD . .

RUN go build -o /etuedt-go

CMD [ "/etuedt-go" ]
