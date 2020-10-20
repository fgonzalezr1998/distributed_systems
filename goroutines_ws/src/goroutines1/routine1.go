package main

import(
	"fmt"
)

func hi(s string) {
	/*
	for i := 0; i < 8; i++ {
		time.Sleep(100 * time.Millisecond)
		fmt.Println(s)
	}
	*/
	fmt.Println(s)
}

func main() {

	/*
	 * ¡IMPORTANTE! Si el hilo principal acaba antes que la subrutina, ésta no
	 * se sigue ejecutando! ----> [Probar a comentar y descomentar la línea 24]
	 */

	go hi("Hello World")
	go hi("Hello World")
	go hi("Hello World")
	for {
		
	}
	/*
	for i := 0; i < 8; i++ {
		// time.Sleep(100 * time.Millisecond)
		fmt.Println("Hii")
	}
	*/
}