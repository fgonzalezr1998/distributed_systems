/*
 * Prints if the introduced words hav five vocals or not
 */

package main

import(
	"fmt"
	"os"
	"bufio"
	"strings"
	pvlcas "my_pkgs/pentavocalicas"
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
	var a, e, i, o, u bool
	
	s := strings.ToUpper(str)

	for _, c := range s {

		if (!a) {
			a = c == 'A'
		}
		if (!e) {
			e = c == 'E'
		}
		if (!i) {
			i = c == 'I'
		}
		if (!o) {
			o = c == 'O'
		}
		if (!u) {
			u = c == 'U'
		}
	}
	
	return a && e && i && o && u
}

func print_if_pentavocalics(l * pvlcas.WordsListType) {

	var success bool

	str := l.Pop(&success)

	for success{

		if (is_petavocalic(str)) {
			fmt.Println("SI")
		} else {
			fmt.Println("NO")
		}

		str = l.Pop(&success)
	}

}

func read_words(l * pvlcas.WordsListType, items int) {
	var str string

	reader := bufio.NewReader(os.Stdin)

	for i := 0; i < items; i++ {
		str, _ = reader.ReadString('\n')
		str = str[:len(str) - 1]
		l.Push(str)
	}
}

func main() {
	var n_words int
	var list pvlcas.WordsListType

	list.First = nil

	pvlcas.Who()

	// Read number of trial cases:

	if (!read_n_cases(&n_words)) {
		os.Exit(1)
	}

	// Read words and print if they are pentavocalics:

	read_words(&list, n_words)
	print_if_pentavocalics(&list)

	os.Exit(0)
}