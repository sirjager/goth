FROM golang:1.22-alpine AS builder
WORKDIR /app

RUN apk update
RUN apk add libc6-compat

COPY go.mod go.sum ./

RUN go mod download
RUN go mod verify

COPY . .

ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0

RUN go build -o main main.go

FROM alpine:latest AS runner
WORKDIR /app

RUN apk add libc6-compat

EXPOSE 4420

RUN mkdir docs

COPY  start.sh .
COPY  wait-for.sh .
COPY defaults.env .

RUN chmod +x start.sh wait-for.sh

COPY --from=builder /app/main .

ENTRYPOINT [ "/app/main" ]

LABEL name="goth"
LABEL org.opencontainers.image.source="https://github.com/sirjager/goth"
