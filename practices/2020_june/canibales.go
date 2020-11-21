package main

import(
	"fmt"
	//"time"
	"sync"
	"os"
)

const NExplorers = 7
const NCannibals = 9
/*
type ChefType struct {
	wakeup chan bool
}

type CannibalType struct {
	last_ate time.Time	// Last time that cannibal ate
}
*/
type TribeType struct {
	chef_ch chan bool	// Channel to wakeup the chef
	n_explorers int32
}

func (t * TribeType) init() {
	t.n_explorers = NExplorers
	t.chef_ch = make(chan bool)
}

func cannibal(tribe * TribeType, mutex * sync.RWMutex) {
	fmt.Println("I am a Cannibal!")
	mutex.Lock()
	tribe.n_explorers--
	mutex.Unlock()
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

	tribe.init()

	go launchCannibals(tribe, &mutex)

	waitForEnd(tribe, &mutex)

	fmt.Println("Hello World")
}