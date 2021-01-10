package toy_factory

import (
	"fmt"
	"sync"
	"math/rand"
	"time"
	"./santa_claus"
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

type ElfType struct {
	problems bool
	is_working bool
	mutex * sync.RWMutex
}

type ToyFactoryType struct {
	santa_claus santa_claus.SantaClausType
	elves [NElves]ElfType
	reindeer_available int32
	elves_with_problems int32
	Time_to_deal_ch chan struct{}
}

func (elf * ElfType) toyFails() (failure bool) {
	failure = rand.Int31n(100) <= FailurePercentage
	return failure
}

func (elf * ElfType) waitForHelp(mutex * sync.RWMutex) {
	var problems bool

	elf.mutex.RLock()
	problems = elf.problems
	elf.mutex.RUnlock()

	if (!problems) {
		elf.is_working = true
		fmt.Println("[ELF] Back to work!")
	}
}

func (tf * ToyFactoryType) helpToElves(mutex * sync.RWMutex) {
	// var t0 time.Time
	var t2s float64
	var elves [NElves]ElfType
	var counter int32

	// Wait for finish the current job:

	tf.santa_claus.WaitForFinish()
	tf.santa_claus.SetWorking(true)
	fmt.Println("[SANTA] Working!")

	t2s = MinHelpTime + rand.Float64() * (MaxHelpTime - MinHelpTime)

	time.Sleep(time.Duration(t2s) * time.Second)
	fmt.Println("[SANTA] Toy checked!")

	// Delete 'ElvesGroup' elves from problems:

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
	mutex.Unlock()
	tf.santa_claus.SetWorking(false)
}

func (tf * ToyFactoryType) runSantaBehavior(mutex * sync.RWMutex) {
	for {
		select{
		case <- tf.santa_claus.Reind_ch:
			// Santa finish his work and go to deal!

			tf.santa_claus.WaitForFinish()
			tf.Time_to_deal_ch <- struct{}{}

		case <- tf.santa_claus.Elv_ch:
			// Help to the elves

			go tf.helpToElves(mutex)
		}
	}
}

func (tf * ToyFactoryType) elfDoWork(elf * ElfType, mutex * sync.RWMutex) bool {
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
			tf.santa_claus.WakeUp()
		}

		return false
	}

	return true
}

func (tf * ToyFactoryType) runElfBehavior(elf * ElfType,
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
				working = tf.elfDoWork(elf, mutex)
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

func (tf * ToyFactoryType) initElves(mutex * sync.RWMutex) {
	for i := 0; i < NElves; i++ {
		// Init one elve:

		tf.elves[i].mutex = new(sync.RWMutex)
		tf.elves[i].problems = false
		tf.elves[i].is_working = false

		// Run its behavior:

		go tf.runElfBehavior(&tf.elves[i], mutex)
	}
}

func (tf * ToyFactoryType) reindeerBehavior(mutex * sync.RWMutex) {
	// When 'waiting_time' has finished, the reindeer arrives

	var waiting_time float64

	waiting_time = MinReindeerInterval +
		rand.Float64() * (MaxReindeerInterval - MinReindeerInterval)

	time.Sleep(time.Duration(waiting_time) * time.Second)

	fmt.Println("***Reindeer Arrived!***")
	tf.reindeer_available++
	if (tf.reindeer_available == NReindeer) {
		tf.santa_claus.Reind_ch <- struct{}{}
	}
}

func (tf * ToyFactoryType) runReindeerBehavior(mutex * sync.RWMutex) {
	for i := 0; i < NReindeer; i++ {
		tf.reindeerBehavior(mutex)
	}
}

func (tf * ToyFactoryType) Init(mutex * sync.RWMutex) {
	tf.reindeer_available = 0
	tf.Time_to_deal_ch = make(chan struct{}, 1)
	tf.elves_with_problems = 0
	tf.santa_claus.Init()
	tf.initElves(mutex)

	go tf.runSantaBehavior(mutex)
	go tf.runReindeerBehavior(mutex)
}