// modulo pricnipal
package main

/*
cliente hecho en golang:
- se intenta hacer mas robusto el codigo
- mas rapidez
- pensado para mayor compatibilidad en computadoras

-- Autor : Urb@n "estamos hack" -- https://www.github.com/Urban20
*/

import (
	conexiones "comando/conex"
	"comando/consola"
	"fmt"
	"os"
)

var ansierr = consola.Iniciar_ANSI()

const TIMEOUT = 5

func main() {

	args := os.Args

	if len(args) != 2 {
		fmt.Print("\n")
		fmt.Println("USO: [programa] [ip:puerto]")
		fmt.Println("Ejemplo: ./[Programa] 192.168.0.50:9999\nNota: ya no es necesaria la flag --IP")
		fmt.Print("\n")
		os.Exit(1)
	}

	if ansierr != nil {
		panic("esta terminal es incompatible con el programa")
	}

	ip := args[1]

	if conex_error := conexiones.Conexion(ip, TIMEOUT); conex_error != nil {

		fmt.Println(conex_error)

		os.Exit(1)
	}

}
