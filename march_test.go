package march

import (
	"fmt"
	"testing"
)

type item struct {
	seqID string
	order int64
}

func (i item) Order() int64 {
	return i.order
}

func Example() {
	chs := []chan Ordered{
		seq("a", 1, 2, 4, 6, 10),
		seq("b", 3, 5, 7),
		seq("c", 8, 9),
	}

	ordered := make(chan Ordered)
	go March(chs, ordered)
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

func TestSingle(t *testing.T) {
	chs := []chan Ordered{
		seq("a", 1, 2, 3, 4, 5),
	}

	ordered := make(chan Ordered)
	go March(chs, ordered)

	check(5, ordered, t)
}

func TestMultiple(t *testing.T) {
	chs := []chan Ordered{
		seq("a", 1, 2, 4, 6, 10, 11, 13),
		seq("b", 3, 5, 7, 12),
		seq("c", 8, 9),
	}

	ordered := make(chan Ordered)
	go March(chs, ordered)

	check(13, ordered, t)
}

func TestEmptyChannelMultiple(t *testing.T) {
	chs := []chan Ordered{
		seq("a", 1, 2, 4, 6, 10, 11, 13),
		seq("b", 3, 5, 7, 12),
		seq("c", 8, 9),
		seq("d"),
	}

	ordered := make(chan Ordered)
	go March(chs, ordered)

	check(13, ordered, t)
}

func TestEmptyChannelSingle(t *testing.T) {
	chs := []chan Ordered{
		seq("a"),
	}

	ordered := make(chan Ordered)
	go March(chs, ordered)

	for range ordered {
		t.Error("Did not expect any results")
	}
}

func TestEmptyChannelSingleOtherItem(t *testing.T) {
	chs := []chan Ordered{
		seq("a"),
		seq("b", 1),
	}

	ordered := make(chan Ordered)
	go March(chs, ordered)

	check(1, ordered, t)
}

// check expects each item in the sequence to be equivalent to its natural order
func check(iterations int, ordered chan Ordered, t *testing.T) {
	count := int64(1)

	for o := range ordered {
		if o.Order() != count {
			t.Errorf("Expected item #%v to equal %v", count, o.Order())
		}
		count++
	}

	if count != int64(iterations+1) {
		t.Error("Expected only 1 iteration")
	}
}

// seq returns a channel holding `vals`
func seq(id string, vals ...int64) chan Ordered {
	ch := make(chan Ordered)

	go func() {
		for _, v := range vals {
			ch <- item{id, v}
		}
		close(ch)
	}()

	return ch
}
