FROM golang:1.24-alpine AS build-db

RUN apk add build-base

WORKDIR /app

COPY packages/db/tool/go.mod ./
COPY packages/db/tool/go.sum ./

RUN go mod download

COPY packages/db/tool/main.go ./main.go
COPY packages/db/tool/cmd ./cmd/
COPY packages/db/tool/internal ./internal/

RUN CGO_ENABLED=1 GOOS=linux go build -o /dbtool ./main.go

##################

FROM golang:1.24-alpine AS build-server

RUN apk add build-base

WORKDIR /app

COPY packages/server/go.mod ./
COPY packages/server/go.sum ./

RUN go mod download

COPY packages/server/cmd ./cmd/
COPY packages/server/internal ./internal/

RUN CGO_ENABLED=1 GOOS=linux go build -o /blingpot ./cmd/blingpot/main.go

##################

FROM node:24-alpine3.21 AS build-web

WORKDIR /app

COPY packages/web/package.json .
COPY packages/web/package-lock.json .

RUN npm ci

COPY packages/web ./

RUN npm run build

##################

FROM alpine:latest AS release

WORKDIR /

COPY --from=build-db /dbtool /dbtool
COPY --from=build-server /blingpot /blingpot
COPY --from=build-web /app/dist/ /app/web/

COPY packages/db/migrations ./app/migrations/

ENV PORT=8888
EXPOSE $PORT

ENV READ_ENV_FILE=skip

ENTRYPOINT ["/blingpot"]
