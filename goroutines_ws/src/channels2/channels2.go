/*
 * Program to demostrate that the lecture from
 * an opened channel is a blocking call
 */

package main

import (
	"fmt"
	"time"
	"os"
)

func iterator(c chan int) {
	for i := 0; i < 10; i++ {
		c <- i
		time.Sleep(100 * time.Millisecond)
	}
}

func main() {
	var finish bool
	var ch chan int
	var n int
	var t0 time.Time

	ch = make(chan int)

	go iterator(ch)

	finish = false
	for (!finish) {
		t0 = time.Now()

		n = <-ch

		fmt.Print("elapsed Time: ")
		fmt.Println(time.Since(t0).Milliseconds())

		if (n == 9) {
			finish = true
		}
	}
	close(ch)
	os.Exit(0)
}