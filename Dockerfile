FROM golang:1.12 AS build

COPY . /app
WORKDIR /app

RUN make

FROM scratch

COPY --from=build \
  /app/build/takanawa \
  /usr/local/bin/takanawa

EXPOSE 5000
ENV HOST=0.0.0.0
ENV PORT=5000
ENTRYPOINT [ "takanawa" ]
