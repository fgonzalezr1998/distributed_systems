package main

import(
	"fmt"
)

func change_num(p * int) {
	/*
	 * I can across 'p' because 'p' was initialized before
	 */

	*p = 42
}

func main() {
	// Declare new integer

	var p int = 21

	// Print integer value

	fmt.Println(p)

	// Call to the function with the pointer to 'p' (its memodry address)

	change_num(&p)

	// Print integer value

	fmt.Println(p)
}