package elves

import(
	"sync"
	"fmt"
)

const NElvesBattalion = 16    // Number of elves per battalion
const NElvesBattalions = 2    // We have two elves battalions

// Main store configuration:

const CacheRows = 3;
const CacheCols = 30;

type CacheType struct {
	elems [CacheRows][CacheCols] bool
}

type ElfType struct {
	problems bool
	is_working bool
	battalion int32
	mutex * sync.RWMutex
	wake_up_ch chan struct{}
}

type ElvesBattalionType struct{
	// leader ElfLeaderType
	Elves [NElvesBattalion - 1]ElfType
	cache CacheType
	presents_to_build int32    // Presents that have to build
	mutex * sync.RWMutex
}

type ElvesType struct {
	Battalions [NElvesBattalions] ElvesBattalionType
	Waiting_elves [] chan struct{}
	main_store CacheType
}

/*
 **********************
 * EXPORTED FUNCTIONS *
 **********************
 */

func (elves * ElvesType) Init() {
	elves.main_store.setAll(true)
	for i := 0; i < NElvesBattalions ; i++ {
		elves.Battalions[i].initBattalion(int32(i))
	}
}

func (elf * ElfType) WaitForHelp(mutex * sync.RWMutex) {

	select {
	case <- elf.wake_up_ch:
		elf.mutex.RLock()
		elf.problems = false
		elf.is_working = true
		elf.mutex.RUnlock()
		fmt.Println("[ELF] Back to work!")
	}
}

func (elf * ElfType) SetWorking(working bool) {
	elf.mutex.RLock()
	elf.is_working = working;
	elf.mutex.RUnlock()
}

func (elf * ElfType) SetProblems(problems bool) {
	elf.mutex.Lock()
	elf.problems = problems
	elf.mutex.Unlock()
}

func (elf * ElfType) IsWorking() bool {
	return elf.is_working
}

func (elves * ElvesType) AddWaitingElf(elf ElfType, mutex * sync.RWMutex) {
	mutex.Lock()
	elves.Waiting_elves = append(elves.Waiting_elves, elf.wake_up_ch)
	mutex.Unlock()
}

/*
 ************************
 * UNEXPORTED FUNCTIONS *
 ************************
 */

func (battalion * ElvesBattalionType) initBattalion(n_bat int32)  {
	initElves(battalion.Elves[:], n_bat)

	// Initialize the battalion cache as empty:

	battalion.cache.setAll(false)

	// Number of present to build start being 0:
	battalion.presents_to_build = 0
	battalion.mutex = new(sync.RWMutex)
}

func initElves(elves [] ElfType, n_bat int32) {
	for i := 0; i < NElvesBattalion - 1; i++ {
		// Init one elve:

		elves[i].mutex = new(sync.RWMutex)
		elves[i].problems = false
		elves[i].is_working = false
		elves[i].battalion = n_bat
		elves[i].wake_up_ch = make(chan struct{}, 1)
	}
}

func (cache * CacheType) setAll(occupied bool) {
	for i := 0; i < CacheRows; i++ {
		for j := 0; j < CacheCols; j++ {
			cache.elems[i][j] = occupied
		}
	}
}