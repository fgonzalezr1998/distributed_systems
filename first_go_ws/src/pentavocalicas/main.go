/*
 * Prints if the introduced words hav five vocals or not
 */

package main

import(
	"fmt"
	"os"
	"bufio"
)

func read_n_cases(n_words * int) bool {
	var elems int
	var err error

	elems, err = fmt.Scanln(n_words)
	
	if (elems != 1) {
		fmt.Print("[Error] ")
		fmt.Println(err)
	}

	return elems == 1
}

func is_petavocalic(str string) bool {
	return true
}

func print_if_pentavocalics(items int) {

	var str string

	reader := bufio.NewReader(os.Stdin)

	for i := 0; i < items ; i++ {
		str, _ = reader.ReadString('\n')

		if (is_petavocalic(str)) {
			fmt.Println("SI")
		} else {
			fmt.Println("NO")
		}
	}

}

func main() {
	var n_words int

	// Read number of trial cases:

	if (!read_n_cases(&n_words)) {
		os.Exit(1)
	}

	// Read words and print if they are pentavocalics:

	print_if_pentavocalics(n_words)

	os.Exit(0)
}