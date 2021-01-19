# builder
ARG GO_VERSION=1.14
FROM golang:${GO_VERSION}-alpine AS builder
RUN apk update && apk add alpine-sdk git && rm -rf /var/cache/apk/*
RUN mkdir -p /api
WORKDIR /api
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o ./app ./*.go

# use builder
FROM golang:${GO_VERSION}-alpine
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
RUN mkdir -p /api
WORKDIR /api
COPY --from=builder /api/app .
EXPOSE 2000
ENTRYPOINT ["./app"]