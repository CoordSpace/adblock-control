FROM golang:1.24-alpine AS builder

# Create unpriv user to run bot
RUN adduser -D -g '' app
WORKDIR /app/src/adblock-control
COPY . .
# Compile the static binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -a -installsuffix cgo .

FROM scratch
LABEL maintainer="CoordSpace - http://coord.space"
# Without this, the binary can't find the template
WORKDIR /app
# This enables using the non-root user
COPY --from=builder /etc/passwd /etc/passwd
# Copy over the compiled binary
COPY --from=builder /app/src/adblock-control/adblock-control /app/src/adblock-control/index.html ./
# Copy over all the templates and icons
COPY --from=builder /app/src/adblock-control/assets/* ./assets/
USER app
CMD ["/app/adblock-control"]
