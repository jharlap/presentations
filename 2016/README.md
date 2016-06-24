# Presentations

* [Testing](#testing)
* [Distributed Correctness](#distributed-correctness)
* [Chorography at Wattpad](#choreography-at-wattpad)
* [Growing by Heaps and Bounds](#growing-by-heaps-and-bounds)

## Testing

This presentation on testing, part philosophy, part Go-specific, will work best if run locally with both mountebank and a GOPATH containing github.com/golang/groupcache and golang.org/x/net/context

But if you don't care about playground working, you can also just browse to http://go-talks.appspot.com/github.com/jharlap/presentations/2016/testing.slide

### Setup

```
npm install -g mountebank
go get -u golang.org/x/tools/cmd/present
go get -u github.com/golang/groupcache
go get -u golang.org/x/net/context

mb --configfile testing/origin.json --allowInjection &
GOPATH=$PWD:$GOPATH present
```

Then browse to http://localhost:3999/testing.slide

If you want Go syntax highlighting, use https://github.com/JosephBuchma/Go-Present-code-highlighter

## Distributed Correctness

Uses cleaver, so `npm install -g cleaver` and `cleaver distributed_correctness.md`

[Slides](distributed_correctness.html)

## Choreography at Wattpad

`cleaver choreography.md`

[Slides](choreography.html)

## Growing by Heaps and Bounds

http://go-talks.appspot.com/github.com/jharlap/presentations/2016/memory_management.slide

