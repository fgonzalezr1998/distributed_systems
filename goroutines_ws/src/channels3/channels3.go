package main

import(
	"fmt"
	"time"
	"os"
)

const N_ITERATIONS = 10

func iterator(delay time.Duration, ch chan int) {
	for i := 0; i < N_ITERATIONS; i++ {
		ch <- i
		time.Sleep(delay)
	}
}

func main() {
	var ch chan int
	var finish bool

	ch = make(chan int)

	go iterator(200 * time.Millisecond, ch)
	go iterator(100 * time.Millisecond, ch)

	finish = false

	for (!finish) {
		n := <-ch

		fmt.Println(n)

		finish = n == N_ITERATIONS - 1
	}

	close(ch)
	os.Exit(0)
}