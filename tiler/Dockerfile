FROM golang:alpine

WORKDIR /service
COPY . .

# Copy the rest of the code
COPY . .

RUN go build -o bin/tiler cmd/service/main.go

EXPOSE 50055

CMD ["./bin/tiler"]
