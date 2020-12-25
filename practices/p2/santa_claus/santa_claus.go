package main

import (
	"fmt"
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

type SantaClausType struct {
	is_working bool;
	elf_ch chan struct{}	// Channel so that the elves can warn Santa
	reindeer_ch chan struct{}	// Channel so that the reindeer can warn Santa
}

type ToyFactory struct {
	santa_claus SantaClausType
	reindeer_available int32
	elves_with_problems int32
}

func main() {
	fmt.Println("Hello World")
}