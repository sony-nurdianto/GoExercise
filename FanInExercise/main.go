package main

import (
	"fmt"
	"sort"
)

func Generator(done <-chan struct{}, num [3][3]int) <-chan [3]int {
	stream := make(chan [3]int)

	go func() {
		defer close(stream)

		for _, n := range num {
			select {
			case <-done:
				return
			case stream <- n:
			}
		}
	}()

	return stream
}

func Indexing(done <-chan struct{}, num <-chan [3]int) <-chan int {
	stream := make(chan int)

	go func() {
		defer close(stream)

		for n := range num {
			for _, v := range n {
				select {
				case <-done:
					return
				case stream <- v:
				}
			}
		}
	}()

	return stream
}

func OrDone(done <-chan struct{}, c <-chan int) <-chan int {
	valueStream := make(chan int)
	go func() {
		defer close(valueStream)
		for {
			select {
			case <-done:
				return
			case v, ok := <-c:
				if !ok {
					return
				}
				select {
				case valueStream <- v:
				case <-done:
				}
			}
		}
	}()

	return valueStream
}

func Grouping(done <-chan struct{}, numStream <-chan int) <-chan []int {
	out := make(chan []int)

	go func() {
		defer close(out)

		var arr []int

		for n := range numStream {
			arr = append(arr, n)
		}

		sort.Ints(arr)

		select {
		case <-done:
			return
		case out <- arr:
		}
	}()

	return out
}

func main() {
	done := make(chan struct{})
	defer close(done)
	twoDArray := [3][3]int{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}

	grouping := Grouping(done, Indexing(done, Generator(done, twoDArray)))

	fmt.Println(<-grouping)
}
