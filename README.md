# 高輪 − Takanawa

[![CircleCI](https://circleci.com/gh/kou64yama/takanawa.svg?style=svg)](https://circleci.com/gh/kou64yama/takanawa)

Takanawa is a reverse proxy for HTTP services for development.

## Usage

For example, execute the following command to proxy `/` to UI server
and `/api` to the API server:

```shell
$ takanawa /api:http://localhost:8080/v1 http://localhost:3000
```

Takanawa runs on port 5000 by default.

## Requirements

- Go 1.11+
