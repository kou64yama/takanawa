# 高輪 − Takanawa

[![CircleCI](https://circleci.com/gh/kou64yama/takanawa.svg?style=svg)](https://circleci.com/gh/kou64yama/takanawa)
[![codecov](https://codecov.io/gh/kou64yama/takanawa/branch/master/graph/badge.svg)](https://codecov.io/gh/kou64yama/takanawa)

Takanawa is a reverse proxy for HTTP services for development.

## Usage

For example, execute the following command to proxy `/` to UI server
and `/api` to the API server:

```shell
$ takanawa --access-log=common \
    http://localhost:3000 --change-origin \
    http://localhost:8080/v1 --change-origin --path=/api
```

Takanawa runs on port 5000 by default.

## Requirements

- Go 1.11+

## Installation

```shell
$ go get -u github.com/kou64yama/takanawa
```

## Docker

```shell
$ docker run -p 5000:5000 kou64yama/takanawa --access-log=common \
    http://localhost:3000 --change-origin \
    http://localhost:8080/v1 --change-origin --path=/api
```
