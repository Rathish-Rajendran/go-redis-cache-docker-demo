FROM golang:1.23-alpine

WORKDIR /app

# Copy go.mod & go.sum first
COPY go.mod ./
RUN go mod download

# Copy the entire project
COPY . .

# Build the app
RUN go build -o main .

EXPOSE 8080

CMD ["./main"]
