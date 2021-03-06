# spdy.go: a SPDY Golang implementation for humans


## About

spdy.go is a SPDY implementation for humans, written in Golang. It offers a higher-level API to the buit-in package available at http://code.google.com/p/go.net/spdy


## Status

spdy.go is under active development. It is alpha software and should not be used in production. Contributions and patches are welcome!

Author: Solomon Hykes <solomon@dotcloud.com>
URL: http://github.com/shykes/spdy-go


## Installation

0. Install GO on your computer (http://golang.org/doc/install)

1. Setup your GO environment:

    $ mkdir ~/go
    
    $ export GOPATH=~/go
    
    $ export PATH=$PATH:$GOPATH/bin

2. Install the library

    $ go get github.com/shykes/spdy-go

3. Install the spdycat command

    $ go get github.com/shykes/spdy-go/spdycat


## Examples


### Serve a web application over spdy

    [Shell A]   $ go run examples/webapp.go -t :8080

    [Chrome]    https://localhost:8080


### Netcat over spdy

    [Shell A]   $ echo "Hi from server" | spdycat -l :4242

    [Shell B]   $ echo "Hi from client" | spdycat :4242


### Stream lots of files to the same recipient

    [Shell A]   $ spdycat -l :5555

    [Shell B]   $ tail -f /var/log/system.log | spdycat :5555 filename=/var/log/system.log

    [Shell C]   $ < ~/.bashrc spdycat :5555 filename=~/.bashrc


## Bugs & missing features

* Doesn't send protocol errors in places where it should

* Barely any testing

* No support for SETTINGS

* No support for GOAWAY

* Stream data is buffered with no watermark limit

