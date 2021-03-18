# Step 1
FROM golang:alpine AS builder
RUN apk update && apk add --no-cache git
WORKDIR $GOPATH/src/rmeharg/gostart/
COPY . .
RUN go get -d -v
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o 
# Step 2
FROM scratch
COPY --from=builder /go/bin/gostart /go/bin/gostart
ENTRYPOINT ["/go/bin/gostart"]
