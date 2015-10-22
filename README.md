goggles
=======

This repository contains some examples of WebGL rendering using GopherJS. To
use it, you'll first need to install Go. To do that on Mac, you could install
using Homebrew:

    brew install go
    export GOPATH=/usr/local/share/go

Then you can run use go get to download the files:

    go get github.com/noonat/goggles

And build things with:

    cd $GOPATH/src/github.com/noonat/goggles
    make

You can run a server with:

    make run

Then you can view the examples at:

- http://127.0.0.1:8080/examples/cube/
- http://127.0.0.1:8080/examples/obj/
