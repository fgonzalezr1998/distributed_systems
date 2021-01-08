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

const MinBuildTime = 5.0
const MaxBuildTime = 8.0

const ElvesGroup = 3

type SantaClausType struct {
	is_working bool;
	elv_ch chan struct{}	// Channel so that the elves can warn Santa
	reind_ch chan struct{}	// Channel so that the reindeer can warn Santa
	mutex * sync.RWMutex
}

type ElveType struct {
	problems bool
	is_working bool
	mutex * sync.RWMutex
}

type ToyFactoryType struct {
	santa_claus SantaClausType
	elves [NElves]ElveType
	reindeer_available int32
	elves_with_problems int32
	time_to_deal_ch chan struct{}
}

func (sc * SantaClausType) init() {
	sc.is_working = false
	sc.elv_ch = make(chan struct{})
	sc.reind_ch = make(chan struct{})
	sc.mutex = new(sync.RWMutex)
}

func (tf * ToyFactoryType) helpToElves(mutex * sync.RWMutex) {
	var t0 time.Time
	var t2s float64
	var elves [NElves]ElveType
	var counter int32

	// Wait for finish the current job:

	for {
		tf.santa_claus.mutex.RLock()
		working := tf.santa_claus.is_working
		tf.santa_claus.mutex.RUnlock()
		if (!working) {
			break
		}
	}

	t2s = MinHelpTime + rand.Float64() * (MaxHelpTime - MinHelpTime)
	t0 = time.Now()
	fmt.Println("[SANTA] Working!")
	for {
		if (time.Since(t0).Seconds() >= t2s) {
			fmt.Println("[SANTA] Toy checked!")
			break
		}
	}

	mutex.RLock()
	elves = tf.elves
	mutex.RUnlock()

	counter = 0
	for i, elf := range elves {
		if (counter == ElvesGroup) {
			break
		}
		if (elf.problems) {
			elf.mutex.Lock()
			tf.elves[i].problems = false
			elf.mutex.Unlock()
			counter++
		}
	}

	mutex.Lock()
	tf.elves_with_problems -= 3
	tf.santa_claus.is_working = false
	fmt.Println(tf.elves_with_problems)
	mutex.Unlock()
}

func (tf * ToyFactoryType) run_santa_behavior(mutex * sync.RWMutex) {
	for {
		select{
		case <- tf.santa_claus.reind_ch:
			// If Santa is not helping to the elves, the program finish

			tf.santa_claus.mutex.RLock()
			if (!tf.santa_claus.is_working) {
				fmt.Println("[SANTA] Time to deal!")
				tf.time_to_deal_ch <- struct{}{}
			}
			tf.santa_claus.mutex.RUnlock()
		case <- tf.santa_claus.elv_ch:
			// Help to the elves
			
			go tf.helpToElves(mutex)
		}
	}
}

func (sc SantaClausType) wakeUp(mutex * sync.RWMutex) {
	var working bool

	sc.mutex.RLock()
	working = sc.is_working
	sc.mutex.RUnlock()

	if (!working) {
		sc.mutex.Lock()
		sc.is_working = true
		sc.mutex.Unlock()

		fmt.Println("Wake up to Santa!")
		sc.elv_ch <- struct{}{}
	}
}

func (elf * ElveType) toyFails() (failure bool) {
	failure = rand.Int31n(100) <= FailurePercentage
	return failure
}

func (tf * ToyFactoryType) ElfDoWork(elf * ElveType, mutex * sync.RWMutex) bool {
	if (elf.toyFails()) {
		fmt.Println("[WARN] Toy Failed!")
		elf.is_working = false
		elf.mutex.Lock()
		elf.problems = true
		elf.mutex.Unlock()

		mutex.Lock()
		tf.elves_with_problems++
		mutex.Unlock()
		mutex.RLock()
		problems := tf.elves_with_problems
		mutex.RUnlock()
		if (problems % ElvesGroup == 0) {
			tf.santa_claus.wakeUp(mutex)
		}

		return false
	}

	return true
}

func (elf * ElveType) waitForHelp(mutex * sync.RWMutex) {
	var problems bool

	elf.mutex.RLock()
	problems = elf.problems
	elf.mutex.RUnlock()

	if (!problems) {
		elf.is_working = true
		fmt.Println("[ELF] Vuelta al trabajo!")
	}
}

func (tf * ToyFactoryType) run_elf_behavior(elf * ElveType,
	mutex * sync.RWMutex) {
	var working bool
	var t0 time.Time
	var t2s float64

	elf.is_working = true
	working = false
	for {
		if (elf.is_working) {
			if (!working) {
				t2s = MinBuildTime + rand.Float64() *
					(MaxBuildTime - MinBuildTime)
				t0 = time.Now()
				working = tf.ElfDoWork(elf, mutex)
			}
			if (working) {
				if (time.Since(t0).Seconds() >= t2s) {
					fmt.Println("[INFO] Toy Success!")
					working = false
				}
			}
		} else {
			elf.waitForHelp(mutex)
		}
	}
}

func (tf * ToyFactoryType) toysListener(mutex * sync.RWMutex) {
	var elves_with_problems int32

	for {
		mutex.RLock()
		elves_with_problems = tf.elves_with_problems
		mutex.RUnlock()
		fmt.Println(elves_with_problems)
		if (elves_with_problems >= ElvesGroup) {
			// Wake up to santa if it is sleeping:

			tf.santa_claus.wakeUp(mutex)
		}
	}
}

func (tf * ToyFactoryType) initElves(mutex * sync.RWMutex) {
	for i := 0; i < NElves; i++ {
		// Init one elve:

		tf.elves[i].mutex = new(sync.RWMutex)
		tf.elves[i].problems = false
		tf.elves[i].is_working = false

		// Run its behavior:

		go tf.run_elf_behavior(&tf.elves[i], mutex)
	}
}

func (tf * ToyFactoryType) waitForReindeer(m * sync.RWMutex) {
	var n_reindeer int32
	for {
		m.RLock()
		n_reindeer = tf.reindeer_available
		m.RUnlock()
		if (n_reindeer == NReindeer) {
			tf.santa_claus.reind_ch <- struct{}{}
			break
		}
	}
}

func (tf * ToyFactoryType)reindeer_behavior(mutex * sync.RWMutex) {
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

func (tf * ToyFactoryType)run_reindeer_behavior(mutex * sync.RWMutex) {
	for i := 0; i < NReindeer; i++ {
		tf.reindeer_behavior(mutex)
	}
}

func (tf * ToyFactoryType) init(mutex * sync.RWMutex) {
	tf.reindeer_available = 0
	tf.time_to_deal_ch = make(chan struct{}, 1)
	tf.elves_with_problems = 0
	tf.santa_claus.init()
	tf.initElves(mutex)

	go tf.run_santa_behavior(mutex)
	go tf.run_reindeer_behavior(mutex)
	go tf.waitForReindeer(mutex)
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

	wait_for_end(toy_factory.time_to_deal_ch, &mutex)

	fmt.Println("Hello World")
}