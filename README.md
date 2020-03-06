# 高輪 − Takanawa

[![GitHub Actions](https://github.com/kou64yama/takanawa/workflows/Go/badge.svg?branch=master)](https://github.com/kou64yama/takanawa/actions?query=workflow%3AGo+branch%3Amaster)
[![Docker Cloud Build Status](https://img.shields.io/docker/cloud/build/kou64yama/takanawa)](https://hub.docker.com/r/kou64yama/takanawa)
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

Takanawa runs on port 5000.

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

## License

Copyright 2020 Yamada Koji. This source code is licensed under the MIT
license.
