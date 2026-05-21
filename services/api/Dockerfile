FROM golang:1.26-alpine AS builder

WORKDIR /app

# go.work + services/api の依存を先にキャッシュ
COPY go.work go.work.sum ./
COPY services/api/go.mod services/api/go.sum ./services/api/
COPY gen/go/go.mod gen/go/go.sum ./gen/go/
RUN go work sync && go mod download

COPY services/api/ ./services/api/
COPY gen/go/ ./gen/go/

RUN CGO_ENABLED=0 GOOS=linux go build -o server ./services/api/cmd/server

FROM alpine:3.21

RUN addgroup -S app && adduser -S -G app app \
  && apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/server .

RUN chown -R app:app /app

USER app

EXPOSE 8080

CMD ["./server"]
