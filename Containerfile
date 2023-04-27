FROM golang:1.20-alpine AS build

WORKDIR /app

# sets env vars for host
ARG TARGETOS
ARG TARGETARCH

COPY go.mod ./

COPY *.go ./

RUN GOARCH=$TARGETARCH GOOS=$TARGETOS go build -o /main

FROM scratch

WORKDIR /

COPY --from=build /main /main
COPY src/facts.txt src/facts.txt
COPY html html
COPY assets assets

ENTRYPOINT ["/main"]
