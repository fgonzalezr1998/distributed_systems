package santa_claus

import (
	"sync"
	"fmt"
)

type SantaClausType struct {
	is_working bool;
	Elv_ch chan struct{}	// Channel so that the elves can warn Santa
	Reind_ch chan struct{}	// Channel so that the reindeer can warn Santa
	mutex * sync.RWMutex
}

func (sc * SantaClausType) Init() {
	sc.is_working = false
	sc.Elv_ch = make(chan struct{})
	sc.Reind_ch = make(chan struct{})
	sc.mutex = new(sync.RWMutex)
}

func (sc * SantaClausType) IsWorking() bool {
	defer sc.mutex.RUnlock()
	
	sc.mutex.RLock()
	return sc.is_working
}

func (sc * SantaClausType) WakeUp() {
	if (!sc.IsWorking()) {
		fmt.Println("Wake up to Santa!")
		sc.Elv_ch <- struct{}{}
	}
}

func (sc * SantaClausType) WaitForFinish() {
	/*
	 * ¡BLOCKING CALL! Wait until Santa finish all his tasks
	 */

	var working bool

	for {
		sc.mutex.RLock()
		working = sc.is_working
		sc.mutex.RUnlock()
		if (!working) {
			break
		}
	}
}

func (sc * SantaClausType) SetWorking(isworking bool) {
	sc.mutex.Lock()
	sc.is_working = isworking
	sc.mutex.Unlock()
}