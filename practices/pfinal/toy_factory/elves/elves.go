package elves

import(
	"sync"
	"fmt"
	"math/rand"
	"time"
)

const NElvesBattalion = 16    // Number of elves per battalion
const NElvesBattalions = 2    // We have two elves battalions

// Main store configuration:

const CacheRows = 3;
const CacheCols = 30;

// Interval that leader get a present from the main store

const LeaderMinInterval = 5.5; // 1.0
const LeaderMaxInterval = 8.5; // 2.0


type CacheType struct {
	mutex * sync.RWMutex
	elems [CacheRows][CacheCols] bool
}

type ElfType struct {
	problems bool
	is_working bool
	battalion int32
	mutex * sync.RWMutex

	// Channel to wake up the elf after a failure:

	wake_up_ch chan struct{}
}

type ElvesBattalionType struct{
	// leader ElfLeaderType
	Elves [NElvesBattalion - 1]ElfType
	cache CacheType    // Cache
	presents_to_build int32    // Presents that have to build
	mutex * sync.RWMutex
}

type ElvesType struct {
	Battalions [NElvesBattalions] ElvesBattalionType
	Waiting_elves [] chan struct{}
	main_store CacheType	// Main memory
	main_store_empty_ch chan struct{}
	Start_working_ch chan struct{}
}

/*
 **********************
 * EXPORTED FUNCTIONS *
 **********************
 */

func (elves * ElvesType) Init(presents_finished_ch chan struct{}) {

	var main_store_readings * int32 = new(int32)
	var mutex * sync.RWMutex = new(sync.RWMutex)

	elves.main_store.setAll(true)
	elves.main_store.mutex = new(sync.RWMutex)
	elves.main_store_empty_ch = presents_finished_ch
	elves.Start_working_ch = make(chan struct{})
	for i := 0; i < NElvesBattalions ; i++ {
		elves.Battalions[i].initBattalion(int32(i))

		// Run leader behavior as Go routine:

		go elves.runLeaderBehavior(main_store_readings, mutex)	
	}
}

func (elf * ElfType) WaitForHelp(mutex * sync.RWMutex) {

	select {
	case <- elf.wake_up_ch:
		elf.mutex.RLock()
		elf.problems = false
		elf.is_working = true
		elf.mutex.RUnlock()
		fmt.Printf("[ELF] From Batallion %d, Back to work!\n",
			elf.battalion + 1)
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

func (elf * ElfType) GetBattalion() int32 {
	return elf.battalion
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

func rowUsed(row int32, used_rows [] int32) bool {
	for _, r := range used_rows {
		if (r == row) {
			return true
		}
	}
	return false
}

func emptyMem(row, col int32, main_store CacheType) bool {
	// Return if main_store[row][col] is empty

	defer main_store.mutex.Unlock()
	main_store.mutex.Lock()

	return !main_store.elems[row][col]
}

func randomRow(used_rows [] int32) int32 {
	var r int32

	r = int32(rand.Intn(CacheRows))
	for (rowUsed(r, used_rows)) {
		r = int32(rand.Intn(CacheRows))
	}
	return r
}

func randomCol(row int32, main_store CacheType) int32 {
	var i, col int32

	col = int32(rand.Intn(CacheCols))

	/*
	 * If the selected column is an empty position,
	 * I move until occupied position. If all collumns
	 * are empty, then -1 is returned and it is handled
	 * by the leader goroutine
	 */

	i = col
	for (i < CacheCols && emptyMem(row, i, main_store)) {
		i++
	}

	if (i == CacheCols) {
		return -1
	}
	return i
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

func (battalion * ElvesBattalionType) DeleteOneFromCache() {
	var i, j int
	var finish bool
	i = 0
	finish = false
	battalion.cache.mutex.RLock()
	for (i < CacheRows && !finish) {
		j = 0
		for (j < CacheCols && !finish) {
			finish = battalion.cache.elems[i][j]
			if (!finish) {
				j++
			}
		}

		if (!finish) {
			i++
		}
	}
	battalion.cache.mutex.RUnlock()
	if (finish) {
		battalion.cache.setAs(int32(i), int32(j), false)
		fmt.Printf("[WARN] (%d, %d) Deleted from Cache!\n", i, j)
	}
}

func (battalion * ElvesBattalionType) initBattalion(n_bat int32)  {
	initElves(battalion.Elves[:], n_bat)

	// Initialize the battalion cache as empty:

	battalion.cache.setAll(false)
	battalion.cache.mutex = new(sync.RWMutex)

	// Number of present to build start being 0:

	battalion.presents_to_build = 0
	battalion.mutex = new(sync.RWMutex)
}

func (cache * CacheType) setAll(occupied bool) {
	for i := 0; i < CacheRows; i++ {
		for j := 0; j < CacheCols; j++ {
			cache.elems[i][j] = occupied
		}
	}
}

func (cache * CacheType) setAs(row, col int32, occupied bool) {
	// Set data 'occupied' at (row,col) position
	cache.mutex.Lock()
	cache.elems[row][col] = occupied
	cache.mutex.Unlock()
}

func (elves * ElvesType) writeOnChaches(row, col int32) {
	fmt.Printf("[WARN] Writing over (%d, %d) on caches\n", row, col)
	for i := 0; i < NElvesBattalions; i++ {
		elves.Battalions[i].cache.setAs(row, col, true)
		elves.Start_working_ch <- struct{}{}
	}
}

func (elves * ElvesType) runLeaderBehavior(mem_readings * int32,
	mutex * sync.RWMutex) {
	var used_rows []int32
	var row, col, i int32
	var t2s float64

	i = 0
	for {
		mutex.RLock()
		fmt.Println(*mem_readings)
		if (*mem_readings == CacheRows * CacheCols) {
			elves.main_store_empty_ch <- struct{}{}
		}
		mutex.RUnlock()
		if (i == CacheRows - 1) {
			used_rows = used_rows[:0]	// Clear slice
			i = 0
		} else {
			i++
		}

		row = randomRow(used_rows)    // Get random Row:
		used_rows = append(used_rows, row)

		col = randomCol(row, elves.main_store)    // Get random col:
		if (col < 0) {
			continue
		}

		// Set position as empty because the data was readen

		elves.main_store.setAs(row, col, false)
		mutex.Lock()
		(*mem_readings)++
		mutex.Unlock()

		// Write data on caches

		elves.writeOnChaches(row, col)

		// Sleep the necessary time:

		t2s = LeaderMinInterval + rand.Float64() *
			(LeaderMaxInterval - LeaderMinInterval)
		time.Sleep(time.Duration(t2s) * time.Second)
	}
}