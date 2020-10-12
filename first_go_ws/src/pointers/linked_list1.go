/*
 * Implements a linked list as a FIFO list (next element at the end)
 */

package main

import "fmt"

type TypeCell struct {
	num int
	next *TypeCell
}

type TypeList struct {
	first *TypeCell
}

func insert_elements(list *TypeList, n_elems int) {
	for i := 0; i < n_elems; i++ {
		push(list, i)
	}
}

func push(list *TypeList, n int) {
	var aux *TypeCell
	var cell_ptr *TypeCell

	// Allocate memory for new cell and set it:

	cell_ptr = new(TypeCell)
	cell_ptr.num = n
	cell_ptr.next = nil

	// introduce the new cell in the list:

	if (list.first == nil) {
		list.first = cell_ptr
	} else {

		// Insert to the end of list:

		aux = list.first
		for aux.next != nil {
			aux = aux.next
		}
		aux.next = cell_ptr
	}
}

func print_list(list TypeList) {
	var aux *TypeCell = list.first

	for aux != nil {
		fmt.Println(aux.num)
		aux = aux.next
	}
}

func main() {
	var linked_list TypeList

	// Input elements in the linked list

	insert_elements(&linked_list, 10)

	// Print all elements from list

	print_list(linked_list)
}