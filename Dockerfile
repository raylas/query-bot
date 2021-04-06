#
# Base
#
FROM golang:1.16.3-alpine3.13 as base
RUN mkdir /build
ADD . /build
WORKDIR /build

#
# Build
#
FROM base as build
RUN go build -o query-bot .

#
# Develop
#
FROM base as develop
CMD ["go","run","main.go"]

#
# App
#
FROM alpine:3.13 as app
RUN mkdir /app
WORKDIR /app
COPY --from=build /build/query-bot .
CMD ["/app/query-bot"]
