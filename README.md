# sqlite-og

## Overview

**SQLite Over GRPC**, is an experimental tool to
enable the separation of app & database server when
using sqlite as the database.

Essentially sqliteog is a golang proxy for sqlite and 
as the name suggest it uses grpc as the database wire protocol.

## Motivation

I've been involved with many small projects where
a relational database is needed. SQLite is perfect
to get up & running quickly, as it requires
no installation & very little -if any- configuration.

However, as soon as the application grows & horizontal scaling
is required, SQLite is no longer viable.

## Goals