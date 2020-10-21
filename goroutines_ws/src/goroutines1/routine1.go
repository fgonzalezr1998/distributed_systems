package main

import(
	"fmt"
	"time"
)

func hi(s string) {
	for i := 0; i < 8; i++ {
		time.Sleep(100 * time.Millisecond)
		fmt.Println(s)
	}
}

func main() {

	/*
	 * ¡IMPORTANTE! Si el hilo principal acaba antes que la subrutina, ésta no
	 * se sigue ejecutando! ----> [Probar a comentar y descomentar la línea 23]
	 */

	go hi("Hola Mundo")

	for i := 0; i < 8; i++ {
		time.Sleep((100 * time.Millisecond) / 2)
		fmt.Println("Hii")
	}
}