package pentavocalicas

import "fmt"

type CellType struct {
	str string
	next * CellType
}

type WordsListType struct {
	First * CellType
	len int
}

// Functions that begins by mayus are exported functions (can be called by other Go file)

func (l * WordsListType) Push(s string) {
	var cell * CellType = new(CellType)

	cell.str = s
	cell.next = nil

	if (l.First == nil) {
		l.First = cell
		l.len = 1
	} else {
		aux := l.First
		for (aux.next != nil) {
			aux = aux.next
		}
		aux.next = cell
		l.len++
	}
}

func (l * WordsListType) Pop(is_ok * bool) (s string) {
	cell := l.First
	if (cell == nil) {
		*is_ok = false
		s = ""
	} else {
		*is_ok = true

		l.First = l.First.next
		l.len--

		s = cell.str
	}

	return s
}

func Who() {
	fmt.Println("I am thw Words List package!")
}