# Go IRC Bot

[![Build Status][1]][2] [![Godoc][3]][4]

## Acquire and build
`go get github.com/StalkR/goircbot`

If it fails, it's likely because [goircbot][5] is based on [fluffle/goirc][6]
`master` branch and Go takes `go1` branch by default.

Solution:

1.  `cd "$GOPATH/src/github.com/fluffle/goirc"`

2.  `git checkout master`

3.  `cd -`

4.  `go install github.com/fluffle/goirc`

5.  and try again: `go get github.com/StalkR/goircbot`

## Hierarchy
* `.` is for `package main` and contains an example bot (`examplebot.go`).

* `bot` directory contains the package itself.

* `plugins` is for plugins; inspire from them to create new plugins.

* `lib` is for little libraries used by plugins.

## IRC library
It uses [fluffle/goirc][6] ([doc][7]). Very good!

## Bugs, comments, questions
Create a [new issue][8] or email [goircbot@stalkr.net][9].

[1]: https://api.travis-ci.org/StalkR/goircbot.png?branch=master
[2]: https://travis-ci.org/StalkR/goircbot
[3]: https://godoc.org/github.com/StalkR/goircbot?status.png
[4]: https://godoc.org/github.com/StalkR/goircbot
[5]: https://github.com/StalkR/goircbot
[6]: https://github.com/fluffle/goirc
[7]: https://godoc.org/github.com/fluffle/goirc
[8]: https://github.com/StalkR/goircbot/issues/new
[9]: mailto:goircbot@stalkr.net
