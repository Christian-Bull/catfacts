FROM golang:alpine3.16

# needed for go modules
ENV GO111MODULE=on

# sets env vars for host
ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

RUN mkdir src

COPY go.mod ./
COPY main.go .
COPY src/facts.txt src/facts.txt
COPY html html
COPY assets assets

RUN GOARCH=$TARGETARCH GOOS=$TARGETOS go build -o /app/main

CMD [ "/app/main" ]