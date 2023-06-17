FROM golang:1.19

WORKDIR /go/src

COPY ./go.mod ./go.mod

RUN go mod download

# add source code
COPY . .

# build the source
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o stock-notify-amd64

ENTRYPOINT ["./stock-notify-amd64"]