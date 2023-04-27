FROM golang:1.20-alpine AS build

WORKDIR /app

# sets env vars for host
ARG TARGETOS
ARG TARGETARCH

COPY go.mod ./

RUN go mod download

COPY *.go ./

RUN go build -o /main

FROM scratch

WORKDIR /

COPY --from=build /main /main
COPY src/facts.txt src/facts.txt
COPY html html
COPY assets assets

ENTRYPOINT ["/main"]