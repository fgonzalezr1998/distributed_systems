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
}
/*
type ElfLeaderType struct {
	// 'true' for occupied and 'false' for empty

	main_store CacheType
}
*/
type ElvesBattalionType struct{
	// leader ElfLeaderType
	Elves [NElvesBattalion - 1]ElfType
	cache CacheType
}

type ElvesType struct {
	Battalions [NElvesBattalions] ElvesBattalionType
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
	var problems bool

	elf.mutex.RLock()
	problems = elf.problems
	elf.mutex.RUnlock()

	if (!problems) {
		elf.is_working = true
		fmt.Println("[ELF] Back to work!")
	}
}

func (elf * ElfType) SetWorking(working bool) {
	elf.is_working = working;
}

func (elf * ElfType) SetProblems(problems bool) {
	elf.mutex.Lock()
	elf.problems = problems
	elf.mutex.Unlock()
}

func (elf * ElfType) IsWorking() bool {
	return elf.is_working
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
}

func initElves(elves [] ElfType, n_bat int32) {
	for i := 0; i < NElvesBattalion; i++ {
		// Init one elve:

		elves[i].mutex = new(sync.RWMutex)
		elves[i].problems = false
		elves[i].is_working = false
		elves[i].battalion = n_bat
	}
}

func (cache * CacheType) setAll(occupied bool) {
	for i := 0; i < CacheRows; i++ {
		for j := 0; j < CacheCols; j++ {
			cache.elems[i][j] = occupied
		}
	}
}