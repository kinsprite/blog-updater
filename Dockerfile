# build
FROM golang:1.13.0-alpine3.10 as build

ENV PORT 8080
EXPOSE 8080

RUN mkdir /app
ADD . /app

ENV GOPROXY https://goproxy.io
WORKDIR  /app
RUN go mod vendor
RUN go build -mod=vendor -o blog-updater .


# release
FROM alpine:3.10
RUN mkdir /app
COPY --from=build /app/blog-updater /app/blog-updater

ENV GIN_MODE release
WORKDIR  /app
CMD ["/app/blog-updater"]
