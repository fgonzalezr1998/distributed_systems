package main

import (
	"fmt"
	"sync"
	"math/rand"
	"time"
	"os"
)

const NElves = 12
const NReindeer = 9

// Interval time that reindeer arrives

const MinReindeerInterval = 5.0
const MaxReindeerInterval = 8.0  // My choice
const FailurePercentage = 33  // 33%

// Interval time Santa Claus spend helping with a toy

const MinHelpTime = 2.0
const MaxHelpTime = 5.0

type SantaClausType struct {
	is_working bool;
	elf_ch chan struct{}	// Channel so that the elves can warn Santa
	reindeer_ch chan struct{}	// Channel so that the reindeer can warn Santa
}

type ElveType struct {
	problems bool
	is_working bool
}

type ToyFactoryType struct {
	santa_claus SantaClausType
	elves [NElves]ElveType
	reindeer_available int32
	time_to_deal_ch chan struct{}
}

func (sc SantaClausType) init() {
	sc.is_working = false
	sc.elf_ch = make(chan struct{})
	sc.reindeer_ch = make(chan struct{})
}

func (tf * ToyFactoryType) initElves() {
	for i := 0; i < NElves; i++ {
		tf.elves[i].problems = false
		tf.elves[i].is_working = false
	}
}

func (tf * ToyFactoryType) waitForReindeer(m * sync.RWMutex) {
	var n_reindeer int32
	for {
		m.RLock()
		n_reindeer = tf.reindeer_available
		m.RUnlock()
		if (n_reindeer == NReindeer) {
			tf.time_to_deal_ch <- struct{}{}
			break
		}
	}
}

func (tf * ToyFactoryType) init(mutex * sync.RWMutex) {
	tf.santa_claus.init()
	tf.initElves()
	tf.reindeer_available = 0
	tf.time_to_deal_ch = make(chan struct{}, 1)

	go tf.waitForReindeer(mutex)
}

func reindeer_behavior(tf * ToyFactoryType, mutex * sync.RWMutex) {
	// When 'waiting_time' has finished, the reindeer arrives

	var waiting_time float64
	var t0 time.Time
	var finish bool

	waiting_time = MinReindeerInterval +
		rand.Float64() * (MaxReindeerInterval - MinReindeerInterval)

	finish = false
	t0 = time.Now()
	for !finish {
		finish = time.Since(t0).Seconds() >= waiting_time
	}
	fmt.Println("***Reindeer Arrived!***")
	mutex.Lock()
	tf.reindeer_available++
	mutex.Unlock()
}

func run_reindeer_behavior(tf * ToyFactoryType, mutex * sync.RWMutex) {
	for i := 0; i < NReindeer; i++ {
		reindeer_behavior(tf, mutex)
	}
}

func wait_for_end(time_to_deal_ch chan struct{}, mutex * sync.RWMutex) {
	for {
		select {
			case <- time_to_deal_ch:
				os.Exit(0)
		}
	}
}

func main() {
	var toy_factory *ToyFactoryType = new(ToyFactoryType)
	var mutex sync.RWMutex

	toy_factory.init(&mutex)	// Initialize the Toy Factory

	go run_reindeer_behavior(toy_factory, &mutex)

	wait_for_end(toy_factory.time_to_deal_ch, &mutex)

	fmt.Println("Hello World")
}