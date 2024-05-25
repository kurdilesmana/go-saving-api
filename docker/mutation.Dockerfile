FROM golang:1.22-alpine AS builder

ARG ENVIRONMENT=development

ENV ENV=${ENVIRONMENT}
ENV PATH="/usr/local/go/bin:${PATH}"

WORKDIR /backend-collection-api

COPY . .

# COPY environment-specific .env*
COPY .env.${ENV} .env.application

RUN apk update && apk add --no-cache gcc libc-dev && \
    go version && go mod download && go mod verify && \
    CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags "-X main.version=1.0.0 -X main.buildTime=$(date +%Y-%m-%d) -s -w" -o ./backend-collection-api

#
FROM alpine:3.18

ENV TZ=Asia/Jakarta

WORKDIR /backend-collection-api

RUN apk update && \
    apk add --no-cache tzdata

COPY --from=builder /backend-collection-api/backend-collection-api .
COPY --from=builder /backend-collection-api/.env.application .env
COPY --from=builder /backend-collection-api/go.mod go.mod

RUN chmod +x ./backend-collection-api

EXPOSE 80
EXPOSE 443 

CMD ["./backend-collection-api"]