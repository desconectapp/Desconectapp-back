# Build stage
FROM golang:1.24.2-alpine3.20 AS builder
LABEL intermediateStageToBeDeleted=true

WORKDIR /home/app

COPY /src/go.mod /src/go.sum ./
RUN go mod download

COPY /src ./

RUN go build -o desconectapp ./main.go

# Test stage
# FROM builder AS desconectapp-test-stage
# CMD ["go", "test", "-v", "./tests"]

# Run stage
FROM alpine:3.20

WORKDIR /home/app

COPY --from=builder /home/app/desconectapp ./

ENTRYPOINT ["./desconectapp"]
