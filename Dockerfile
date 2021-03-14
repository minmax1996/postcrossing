FROM golang:1.15.0 AS builder
WORKDIR /go/src/github.com/minmax1996/postcrossing/
COPY main.go .
COPY go.mod .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/src/github.com/minmax1996/postcrossing/app .
CMD ["./app"] 