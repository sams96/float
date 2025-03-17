FROM golang:1.24-alpine AS builder

WORKDIR ${GOPATH}/src/github.com/sams96/float

# re-add this if an external dependency is added
# COPY go.mod go.sum ${GOPATH}/src/github.com/sams96/float/
# RUN go mod download

COPY . ${GOPATH}/src/github.com/sams96/float/

RUN go build -o /go/bin/float .

FROM docker

COPY --from=builder /go/bin/float /usr/bin/float

ENTRYPOINT ["/usr/bin/float"]
