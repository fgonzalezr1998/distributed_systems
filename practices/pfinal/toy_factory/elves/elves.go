package elves

import(
	"sync"
)

const NElvesBattalion = 16    // Number of elves per battalion
const NElvesBattalions = 2    // We have two elves battalions

// Main store configuration:

const CacheTypeRows = 3;
const CacheTypeCols = 30;

type CacheTypeColsType struct {
	cols [CacheTypeCols] bool
}

type CacheType struct {
	rows [CacheTypeRows] CacheTypeColsType
}

type ElfType struct {
	problems bool
	is_working bool
	mutex * sync.RWMutex
}

type ElfLeaderType struct {
	main_store CacheType
}

type ElvesBattalionType struct{
	leader ElfLeaderType
	elves [NElvesBattalion - 1]ElfType
	cache CacheType
}

type ElvesType struct {
	battalions [NElvesBattalions] ElvesBattalionType
}

func (battalion * ElvesBattalionType) Init() {
	
}