FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.work ./
COPY src/packages/proto ./src/packages/proto
COPY src/apps/backend/iam/go.mod ./src/apps/backend/iam/
COPY src/apps/backend/server/go.mod ./src/apps/backend/server/
COPY src/apps/backend/runner/go.mod ./src/apps/backend/runner/

COPY src/apps/backend/runner ./src/apps/backend/runner

WORKDIR /app/src/apps/backend/runner
RUN go mod download
RUN go build -o main .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/src/apps/backend/runner/main .
EXPOSE 8083 50053
CMD ["./main"]
