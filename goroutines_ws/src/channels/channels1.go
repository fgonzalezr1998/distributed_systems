/*
 *                 --- ¡IMPORTANTE! ---
 * Los channels lidian con la concurrencia. Es decir, no necesito
 * usar locks
*/

package main

import "fmt"

func sum(s []int, c chan int) {
	sum := 0
	for _, v := range s {
		sum += v
	}

	// Introducimos el dato en el canal

	c <- sum	// send sum to c
}

func main() {

	// Se define un array de int:

	s := []int{7, 2, 8, -9, 4, 0}

	// Creamos el channel y llamamos a las goroutines:

	c := make(chan int)
	go sum(s[:len(s)/2], c)
	go sum(s[len(s)/2:], c)

	/*
	 *****************************************************************
	 ******************** ---- ¡PREGUNTA! ---- ***********************
	 * ¿Cómo puedo garantizar que el primer dato que recibo del canal*
	 * se corresponde con el de la primera rutina que se lanzó?*******
	 *****************************************************************
	*/

	// Obtenemos los valores devueltos por las subrutinas del canal:

	x := <-c
	y := <-c	// receive from c

	fmt.Println(x, y)
}
