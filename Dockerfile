FROM golang:1.13 AS build

COPY . /app
WORKDIR /app

ARG VERSION=0.0.0
RUN make VERSION=${VERSION}

FROM alpine:3.9

COPY --from=build /app/bin/* /usr/local/bin/

EXPOSE 5000
ENV HOST=0.0.0.0
ENV PORT=5000
ENTRYPOINT [ "takanawa" ]
