# march

[![GoDoc](https://godoc.org/github.com/notduncansmith/march?status.svg)](https://godoc.org/github.com/notduncansmith/march) [![Build Status](https://travis-ci.com/notduncansmith/march.svg?branch=master)](https://travis-ci.com/notduncansmith/march) [![codecov](https://codecov.io/gh/notduncansmith/march/branch/master/graph/badge.svg)](https://codecov.io/gh/notduncansmith/march)

march helps produce ordered output from multiple channels. It works by keeping track of the latest item received on each channel, and continually reading from whichever channel whose latest item has the lowest order, publishing items to a result channel in order.

## Example

```go
package main

import (
    "fmt"

    "github.com/notduncansmith/march"
)

type item struct {
	seqID string
	order int64
}

// Order implements march.Ordered
func (i item) Order() int64 {
	return i.order
}

func main() {
    chs := []chan march.Ordered{
		seq("a", 1, 2, 4, 6, 10),
		seq("b", 3, 5, 7),
		seq("c", 8, 9),
	}

	ordered := make(chan march.Ordered)
	go march.March(chs, ordered)
	for o := range ordered {
		fmt.Println(o)
	}
	// Output:
	// {a 1}
	// {a 2}
	// {b 3}
	// {a 4}
	// {b 5}
	// {a 6}
	// {b 7}
	// {c 8}
	// {c 9}
	// {a 10}
}

// seq returns a channel holding `vals`
func seq(id string, vals ...int64) chan march.Ordered {
	ch := make(chan march.Ordered)

	go func() {
		for _, v := range vals {
			ch <- item{id, v}
		}
		close(ch)
	}()

	return ch
}
```

See [tests](./march_test.go) for more examples.

## License

Released under [The MIT License](https://opensource.org/licenses/MIT) (see `LICENSE.txt`).

Copyright 2019 Duncan Smith