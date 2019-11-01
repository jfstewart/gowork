package main

import (
	"fmt"
	"sync"
)

type Item = int

type Queue struct {
	queue chan Item
}

func NewQueue() *Queue {
	q := new(Queue)
	q.queue = make(chan Item, 100)
	return q
}

func (q *Queue) Put(item Item) {
	q.queue <- item
}

func (q *Queue) GetMany(n int) []Item {
	buff := make([]Item, 0, n)

	for {
		i := <-q.queue
		buff = append(buff, i)

		if len(buff) == n {
			return buff
		}
	}
}

func main() {
	q := NewQueue()

	var wg sync.WaitGroup
	for n := 10; n > 0; n-- {
		wg.Add(1)
		go func(n int) {
			items := q.GetMany(n)
			fmt.Printf("%2d: %2d\n", n, items)
			wg.Done()
		}(n)
	}

	for i := 0; i < 100; i++ {
		q.Put(i)
	}
	defer close(q.queue)

	wg.Wait()
}
