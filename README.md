# go-get

A repository fetcher, forked from golang/go

## Why?

> When in doubt, simply port Go's source code, documentation, and tests.

- from *Deno Standard Modules*, https://deno.land/std@0.68.0

[Tilt](https://tilt.dev/) needs a system for importing extensions.

We love the Go package-import system.

We decided to copy it!

But when we looked at how Go's `go get` was implemented, 
we saw that it supports a lot of different repositories.

This package contains a fork of that package, to make it easier to re-use.

## How?

TK

## License

Licensed under [3-clause BSD](LICENSE)

Originally Copyright (c) 2009 The Go Authors. All rights reserved.

Modified by Windmill Engineering, Inc.
