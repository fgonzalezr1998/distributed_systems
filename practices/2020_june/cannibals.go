package main

import(
	"fmt"
	"time"
	"sync"
	"os"
	"math/rand"
)

const NExplorers = 7
const NCannibals = 9
const NPortionsByExplorer = 5

type TribeType struct {
	chef_ch chan struct{}	// Channel to wakeup the chef
	end_ch chan struct{}
	n_explorers int32
	portions int32		// Portions in the pot
	chef_isworking bool
}

func (t * TribeType) init(mutex * sync.RWMutex) {
	t.n_explorers = NExplorers
	t.chef_ch = make(chan struct{}, 1)
	t.end_ch = make(chan struct{})
	t.portions = 0
	t.chef_isworking = false
	go t.chefDoWork(mutex)
}

func (t * TribeType) chefDoWork(mutex * sync.RWMutex) {
	var c chan struct{}

	c = make(chan struct{})
	for {
		select {
			case <-t.chef_ch:
				// The chef has been waked up

				// Decrease the number of explorers:
				
				mutex.Lock()
				t.n_explorers--
				if (t.n_explorers == 0) {
					t.end_ch <- struct{}{}
				}
				mutex.Unlock()

				go cook(c)
			case <-c:
				// The work of the chef has finished

				mutex.Lock()
				t.portions = 5
				t.chef_isworking = false
				mutex.Unlock()
		}
	}
}

func (t * TribeType) potIsEmpty(mutex * sync.RWMutex) bool {
	var portions int32

	defer mutex.RUnlock()

	mutex.RLock()
	portions = t.portions

	return portions == 0
}

func (t * TribeType) callToChef(mutex * sync.RWMutex) {
	mutex.Lock()

	if (!t.chef_isworking){
		t.chef_isworking = true
		fmt.Println("[Cannibal] Waking up to the Chef!")
		t.chef_ch <- struct{}{}
	}
	mutex.Unlock()
}

func cook(c chan struct{}) {
	for i := 0; i < NPortionsByExplorer; i++ {
		// Wait randomly time between 1.0 and 2.0 seconds

		time.Sleep(time.Duration(rand.Float32() + 1.0) * time.Second)
		fmt.Println("[Chef] One portion made!")
	}
	c <- struct{}{}
}

func eat(tribe * TribeType, mutex * sync.RWMutex, state * int32) {
	/*
	 * If pot is empty, call to the chef. Else, decrease one portion and
	 * sleep 0.5-1.0 seconds before to change the state to 'work'
	 */

	if (tribe.potIsEmpty(mutex)) {
		tribe.callToChef(mutex)
	} else {
		fmt.Println("[Cannibal] Eating!")
		mutex.Lock()
		tribe.portions--
		mutex.Unlock()
	}
}

func cannibal(tribe * TribeType, mutex * sync.RWMutex) {
	/*
	 * Cannibal behaviour is implemented as a simple states machine
	 */

	var eat_st, state int32
	var t0 time.Time
	var t2s float64
	var eating, working bool
	eat_st = 0
	
	state = eat_st
	eating = false
	working = false
	for {
		if (state == eat_st) {
			if (!eating) {
				eat(tribe, mutex, &state)

				eating = true

				// Sleep randomly time between 0.5s - 1.0s

				t2s = 0.5 + rand.Float64() * (1.0 - 0.5)
				t0 = time.Now()
			} else {
				if (time.Since(t0).Seconds() >= t2s) {
					state = 1
					eating = false
				}
			}
		} else {
			// Work

			if (!working) {
				fmt.Println("[Cannibal] Working!")
				working = true
				t2s = 0.5 + rand.Float64() * (2.0 - 0.5)
				t0 = time.Now()
			} else {
				if (time.Since(t0).Seconds() >= t2s) {
					state = eat_st
					working = false
				}
			}
		}
	}
}

func launchCannibals(tribe * TribeType, mutex * sync.RWMutex) {
	for i := 0; i < NCannibals; i++ {
		go cannibal(tribe, mutex)
	}
}

func waitForEnd(tribe * TribeType) {
	for {
		select {
			case <- tribe.end_ch:
				os.Exit(0)
		}
	}
}

func main() {
	var mutex sync.RWMutex	// Read/Write mutex
	var tribe *TribeType = new(TribeType)

	tribe.init(&mutex)

	launchCannibals(tribe, &mutex)

	waitForEnd(tribe)
}