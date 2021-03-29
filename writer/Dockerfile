# Step 1
ARG BUILDER_IMAGE=golang:alpine

FROM ${BUILDER_IMAGE} as builder

# Install git + SSL ca certificates.
# Git is required for fetching the dependencies.
# Ca-certificates is required to call HTTPS endpoints.
RUN apk update && apk add --no-cache git ca-certificates tzdata && update-ca-certificates

# Create appuser
ENV USER=appuser
ENV UID=10001 
RUN adduser \    
    --disabled-password \    
    --gecos "" \    
    --home "/nonexistent" \    
    --shell "/sbin/nologin" \    
    --no-create-home \    
    --uid "${UID}" \    
    "${USER}"

WORKDIR $GOPATH/src/github.com/rmeharg/gomessage/writer

# Use modules, fetch dependencies
COPY go.mod .
ENV GO111MODULE=on
RUN go mod download
RUN go mod verify

COPY . .

RUN touch /report.pdf
RUN chown -R appuser:appuser /report.pdf

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' -a \
    -o /go/bin/writer .

# Step 2
FROM scratch

# Import from builder.
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder /report.pdf /report.pdf

# Copy our static executable
COPY --from=builder /go/bin/writer /go/bin/writer

# Use an unprivileged user.
USER appuser:appuser

# Expose on a port > 1024
EXPOSE 8080

# Run the gomessage binary.
ENTRYPOINT ["/go/bin/writer"]
