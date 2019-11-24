package march

// Ordered represents an item with an order
type Ordered interface {
	Order() int64
}

// March returns the contents of each channel in order
func March(chs []chan Ordered, results chan Ordered) {
	size := len(chs)
	latest := make([]Ordered, size)
	open := make([]bool, size)
	closedCount := 0

	for i := range chs {
		latest[i], open[i] = <-chs[i]
		if !open[i] {
			closedCount++
		}
	}

	if closedCount == size {
		close(results)
		return
	}

	var next Ordered
	var nextI int

	for {
		for i, o := range latest {
			if o != nil && (next == nil || o.Order() < next.Order()) {
				next = o
				nextI = i
			}
		}

		results <- next
		next = nil

		if open[nextI] {
			latest[nextI], open[nextI] = <-chs[nextI]
			if !open[nextI] {
				closedCount++
			}
		}

		if closedCount == size {
			close(results)
			return
		}
	}
}
