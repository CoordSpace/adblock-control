FROM golang:1.14-alpine AS builder

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
# Copy over all the templates and icons
COPY --from=builder $GOPATH/src/adblock-control/assets/* /go/bin/adblock-control/assets
USER app

CMD ["/go/bin/adblock-control"]
