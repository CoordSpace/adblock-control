FROM golang:1.14-alpine AS builder

RUN apk update && apk add --no-cache git=2.27.0-r0 ca-certificates=20191127-r4 tzdata=2020a-r0 && update-ca-certificates

# Create unpriv user to run bot
RUN adduser -D -g '' app

WORKDIR $GOPATH/src/adblock-control

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -a -installsuffix cgo -o /go/bin/adblock-control .

FROM scratch

LABEL maintainer="CoordSpace"
LABEL build_version="Linuxserver.io version:- ${VERSION} Build-date:- ${BUILD_DATE}"

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /go/bin/adblock-control /go/bin/adblock-control
# Copy over the root certs from ca-certificates for SSL
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
USER appuser

CMD ["/go/bin/adblock-control"]
