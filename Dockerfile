# Builder
FROM golang:1.18-alpine AS build

ARG VERSION=dev

WORKDIR /tmp/kantoku

COPY . .

RUN apk add --no-cache git && \
    go mod download && \
    go mod verify && \
    go build -ldflags="-X 'main.version=${VERSION}'" -o kantoku

# Runner
FROM alpine:latest

WORKDIR /opt/kantoku

COPY --from=build /tmp/kantoku/kantoku /opt/kantoku/kantoku

EXPOSE 80

CMD [ "./kantoku" ]
