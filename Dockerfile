FROM golang:alpine

# Working directory
WORKDIR /app

# Copy all content to the working container's directory
COPY ./ /app

# Download all required dependencies and the package needed for hot reload
RUN go mod download
RUN go get github.com/githubnemo/CompileDaemon

# Build app and run :)
ENTRYPOINT CompileDaemon --build="go build main.go" --command=./main -log-prefix=false