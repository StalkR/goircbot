# Go IRC Bot

[![godoc](https://godoc.org/github.com/StalkR/goircbot?status.png)](https://godoc.org/github.com/StalkR/goircbot)
[![build status](https://github.com/StalkR/goircbot/actions/workflows/build.yml/badge.svg)](https://github.com/StalkR/goircbot/actions/workflows/build.yml)
[![test status](https://github.com/StalkR/goircbot/actions/workflows/test.yml/badge.svg)](https://github.com/StalkR/goircbot/actions/workflows/test.yml)

## Get and build

`go get github.com/StalkR/goircbot`

## Hierarchy

* `.` contains `examplebot.go`, an example bot binary.

* `bot` directory contains the bot library.

* `plugins` is a directory with plugin libraries;
  inspire from them to create new plugins.

* `lib` is for little libraries used by plugins,
  which may be reused.

## IRC library

It uses [fluffle/goirc](https://github.com/fluffle/goirc)
([doc](https://godoc.org/github.com/fluffle/goirc)).

## Bugs, comments, questions

Create a [new issue](https://github.com/StalkR/goircbot/issues/new).
