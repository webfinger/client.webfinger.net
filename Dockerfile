# syntax=docker/dockerfile:1.4
FROM cgr.dev/chainguard/go:latest as build
LABEL maintainer="Will Norris <will@willnorris.com>"

WORKDIR /work
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -v -o webfinger

FROM cgr.dev/chainguard/static:latest

COPY --from=build /work/webfinger /webfinger

CMD ["/webfinger"]
