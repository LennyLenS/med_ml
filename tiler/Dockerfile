FROM golang:alpine

WORKDIR /service

# Copy go mod files first for better caching
COPY go.mod go.sum* ./
RUN go mod download

# Copy the rest of the code
COPY . .

RUN go build -o bin/tiler cmd/service/main.go

EXPOSE 50055

CMD ["./bin/tiler"]
