# Buid app
FROM golang:1.16-alpine as builder
WORKDIR /go/src/app/

RUN apk add --no-cache git make gcc g++

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN make build

# Finish image
FROM alpine
WORKDIR /app

COPY --from=builder /go/src/app/webhook .

ENTRYPOINT [ "./webhook" ]