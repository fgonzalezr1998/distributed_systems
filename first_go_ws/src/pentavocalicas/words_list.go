
package words_list

type CellType struct {
	str string
	next * CellType
}

type WordsListType struct {
	first * CellType
}