FROM golang:1.19

WORKDIR /app
# Retrieve application dependencies.
# This allows the container build to reuse cached dependencies.
# Expecting to copy go.mod and if present go.sum.
COPY go.* ./
RUN go mod download

RUN go mod download golang.org/x/term

COPY . .

# Build the binary.
RUN go build -o main

CMD ["/app/main"]