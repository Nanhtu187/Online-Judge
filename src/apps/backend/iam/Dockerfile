FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy the entire workspace context to handle dependencies
COPY go.work ./
COPY src/packages/proto ./src/packages/proto
COPY src/apps/backend/iam/go.mod ./src/apps/backend/iam/
COPY src/apps/backend/server/go.mod ./src/apps/backend/server/
COPY src/apps/backend/runner/go.mod ./src/apps/backend/runner/

# Copy the actual source for IAM
COPY src/apps/backend/iam ./src/apps/backend/iam

# Tidy and build
WORKDIR /app/src/apps/backend/iam
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/src/apps/backend/iam/main .
RUN chmod +x ./main
EXPOSE 8084 50051
CMD ["./main"]
