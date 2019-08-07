FROM golang:alpine AS builder

RUN apk add --no-cache git curl
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

WORKDIR $GOPATH/src/github.com/kuzaxak/promalert

COPY Gopkg.toml Gopkg.lock ./
RUN dep ensure --vendor-only -v

COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o /app .

FROM alpine

COPY --from=builder /app ./
COPY config.example.yaml ./config.yaml

RUN apk add --no-cache ca-certificates && \
    update-ca-certificates

ENTRYPOINT ["./app"]
EXPOSE 8080