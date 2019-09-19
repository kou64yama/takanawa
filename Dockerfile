FROM golang:1.12 AS build

COPY . /go/src/github.com/kou64yama/takanawa
WORKDIR /go/src/github.com/kou64yama/takanawa

RUN go get -u github.com/golang/dep/cmd/dep \
  && dep ensure \
  && make

FROM scratch

COPY --from=build \
  /go/src/github.com/kou64yama/takanawa/build/takanawa \
  /usr/local/bin/takanawa

EXPOSE 5000
ENV HOST=0.0.0.0
ENV PORT=5000
ENTRYPOINT [ "takanawa" ]
