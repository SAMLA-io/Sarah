FROM golang:1.24-alpine

# Install git for getting commit hash
RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download && go mod verify

# Copy the entire repository including .git directory
COPY . .

# Get git commit hash and build with it
RUN git rev-parse --short HEAD > .git_commit_hash

RUN go build -ldflags="-X main.GitCommitHash=$(cat .git_commit_hash)" -o main .

EXPOSE 8080
CMD ["./main"]