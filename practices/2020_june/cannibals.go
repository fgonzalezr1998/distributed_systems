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

/*
type ChefType struct {
	wakeup chan bool
}

type CannibalType struct {
	last_ate time.Time	// Last time that cannibal ate
}
*/
type TribeType struct {
	chef_ch chan struct{}	// Channel to wakeup the chef
	n_explorers int32
	portions int32		// Portions of the explorer
	chef_isworking bool
}

func (t * TribeType) init(mutex * sync.RWMutex) {
	t.n_explorers = NExplorers
	t.chef_ch = make(chan struct{})
	t.portions = 0
	t.chef_isworking = false
	go t.chefDoWork(mutex)
}

func (t * TribeType) chefDoWork(mutex * sync.RWMutex) {
	for {
		select {
			case <-t.chef_ch:
				// Decrease the number of explorers:

				mutex.Lock()
				t.n_explorers--
				t.chef_isworking = true
				mutex.Unlock()

				for i := 0; i < NPortionsByExplorer; i++ {
					// Wait randomly time between 1.0 and 2.0 seconds

					time.Sleep(time.Duration(rand.Float32() + 1.0) * time.Second)
					fmt.Println("[Chef] One portion made!")
				}

				mutex.Lock()
				t.portions = 5
				t.chef_isworking = false
				mutex.Unlock()
		}
	}
}

func (t * TribeType) potIsEmpty(mutex * sync.RWMutex) bool {
	mutex.RLock()
	portions := t.portions
	mutex.RUnlock()

	return portions == 0
}

func (t * TribeType) callToChef(mutex * sync.RWMutex) {
	var chef_isworking bool

	mutex.RLock()
	chef_isworking = t.chef_isworking
	mutex.RUnlock()

	if (!chef_isworking){
		t.chef_ch <- struct{}{}
	}
}

func cannibal(tribe * TribeType, mutex * sync.RWMutex) {
	for {
		if (tribe.potIsEmpty(mutex)) {
			tribe.callToChef(mutex)
		}
		fmt.Println("I am a Cannibal!")
		mutex.Lock()
		tribe.portions--
		mutex.Unlock()
	}
}

func launchCannibals(tribe * TribeType, mutex * sync.RWMutex) {
	for i := 0; i < NCannibals; i++ {
		go cannibal(tribe, mutex)
	}
}

func waitForEnd(tribe * TribeType, mutex * sync.RWMutex) {
	var n_explorers int32
	for {
		mutex.RLock()
		n_explorers = tribe.n_explorers
		mutex.RUnlock()

		if (n_explorers == 0) {
			os.Exit(0)
		}
	}
}

func main() {
	var mutex sync.RWMutex	// Read/Write mutex
	var tribe *TribeType = new(TribeType)

	tribe.init(&mutex)

	go launchCannibals(tribe, &mutex)

	waitForEnd(tribe, &mutex)

	fmt.Println("Hello World")
}