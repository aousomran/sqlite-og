# sqlite-og

![build workflow](https://github.com/aousomran/sqlite-og/actions/workflows/build.yml/badge.svg)

## Overview

**SQLite Over GRPC**, is an experimental tool to
enable the separation of app & database server when
using sqlite as the database.

Essentially sqliteog is a golang proxy for sqlite and
as the name suggest it uses grpc as the database wire protocol.

## Quickstart

### Installation

- Run via docker
    ```shell
    docker run -p 9091:9091 aousomran/sqlite-og:latest
    ```
