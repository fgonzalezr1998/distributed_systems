package main

import (
	"sync"
	"os"
	"fmt"
	"./toy_factory"
)

func wait_for_end(time_to_deal_ch chan struct{}) {
	for {
		select {
			case <- time_to_deal_ch:
				fmt.Println("[INFO] TIME TO DEAL!")
				os.Exit(0)
		}
	}
}

func main() {
	var toy_factory *toy_factory.ToyFactoryType = 
		new(toy_factory.ToyFactoryType)
	var mutex sync.RWMutex

	toy_factory.Init(&mutex)	// Initialize the Toy Factory

	wait_for_end(toy_factory.Time_to_deal_ch)
}