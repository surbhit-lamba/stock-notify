
FROM golang:1.19

WORKDIR /go/src

COPY ./go.mod ./go.mod

RUN go mod download

# add source code
COPY . .

RUN -d \
 --name newrelic-infra \
 --network=host \
 --cap-add=SYS_PTRACE \
 --privileged \
 --pid=host \
 -v "/:/host:ro" \
 -v "/var/run/docker.sock:/var/run/docker.sock" \
 -e  NRIA_LICENSE_KEY=8b902c7cb9de77e972b811d71939d9f5aec3NRAL \
 newrelic/infrastructure:latest

# build the source
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o stock-notify-amd64

ENTRYPOINT ["./stock-notify-amd64"]