# Go IRC Bot

[![Build Status][5]][6]

## Acquire and build
`go get github.com/StalkR/goircbot`

It may fail with errors because [goircbot][1] is based on [fluffle/goirc][3]
master branch but it defaults to branch `go1`.

Solution:

1.  `cd "$GOPATH/src/github.com/fluffle/goirc"`

2.  `git checkout master`

3.  `git pull`

4.  `cd -`

5.  `go get github.com/StalkR/goircbot`

## Hierarchy
* `.` is for `package main` and contains an example bot (`examplebot.go`).

* `bot` directory contains the package itself.

* `plugins` is for plugins; inspire from them to create new plugins.

* `lib` is for little libraries used by plugins.

## Documentation
On GoDoc: [godoc.org/github.com/StalkR/goircbot][2].

## IRC library
It uses [fluffle/goirc][3] ([doc][4]). Very good!

## Bugs, comments, questions
Create a [new issue][9] or email [goircbot@stalkr.net][8].

[1]: http://github.com/StalkR/goircbot
[2]: http://godoc.org/github.com/StalkR/goircbot
[3]: http://github.com/fluffle/goirc
[4]: http://godoc.org/github.com/fluffle/goirc/client
[5]: https://secure.travis-ci.org/StalkR/goircbot.png
[6]: http://www.travis-ci.org/StalkR/goircbot
[7]: http://godoc.org
[8]: mailto:goircbot@stalkr.net
[9]: https://github.com/StalkR/goircbot/issues/new
