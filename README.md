# Go IRC Bot

[![Build Status][1]][2] [![Godoc][3]][4]

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
It uses [fluffle/goirc][6] ([doc][7]).

## Bugs, comments, questions
Create a [new issue][8].

[1]: https://api.travis-ci.org/StalkR/goircbot.png?branch=master
[2]: https://travis-ci.org/StalkR/goircbot
[3]: https://godoc.org/github.com/StalkR/goircbot?status.png
[4]: https://godoc.org/github.com/StalkR/goircbot
[5]: https://github.com/StalkR/goircbot
[6]: https://github.com/fluffle/goirc
[7]: https://godoc.org/github.com/fluffle/goirc
[8]: https://github.com/StalkR/goircbot/issues/new
