FROM golang:1.24-alpine AS build

RUN apk add build-base

WORKDIR /app

COPY packages/server/go.mod ./
COPY packages/server/go.sum ./

RUN go mod download

COPY packages/server/cmd ./cmd/
COPY packages/server/internal ./internal/

RUN CGO_ENABLED=1 GOOS=linux go build -o /blingpot ./cmd/blingpot/main.go

FROM alpine:latest AS release

WORKDIR /

COPY --from=build /blingpot /blingpot

ENV PORT=8888
EXPOSE $PORT

ENV READ_ENV_FILE=skip

ENTRYPOINT ["/blingpot"]
