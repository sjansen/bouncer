FROM golang:1.17-alpine as builder
ARG FN=rekey-fn
ARG MODULE=github.com/sjansen/bouncer
RUN apk add --update \
    ca-certificates \
    git
ADD go.mod go.sum /go/src/${MODULE}/
RUN cd /go/src/${MODULE} && \
    go mod download

ARG GITSHA="(missing)"
ARG TIMESTAMP="(missing)"
ADD cmd /go/src/${MODULE}/cmd
ADD internal /go/src/${MODULE}/internal
RUN cd /go/src/${MODULE} && \
    echo GITSHA=${GITSHA} && \
    echo TIMESTAMP="${TIMESTAMP}" && \
    CGO_ENABLED=0 GOARCH=amd64 GOOS=linux \
    go build \
        -o /lambda-handler \
        -a -installsuffix cgo \
        -ldflags="-s -w -X ${MODULE}/internal/build.GitSHA=${GITSHA} -X '${MODULE}/internal/build.Timestamp=${TIMESTAMP}'" \
        ./cmd/${FN}/

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /lambda-handler /lambda-handler
EXPOSE 8000
WORKDIR /
ENTRYPOINT ["/lambda-handler"]
