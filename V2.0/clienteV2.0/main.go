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
	color "comando/colores"
	conexiones "comando/conex"
	"flag"
	"fmt"
	"os"
)

var arg = flag.String("IP", "", "[ip]:[puerto del host]")

const timeout = 10

func main() {

	flag.Parse()

	ip := *arg

	if ip != "" {
		fmt.Println(color.Violeta + color.LOGO + color.Reset)
		error := conexiones.Conexion(ip, timeout)
		if error != nil {

			fmt.Println("error: ", error)
			os.Exit(1)
		}

	} else {
		fmt.Println(color.Rojo + "[!] ingresar un valor de ip (ip:puerto)" + color.Reset)
		os.Exit(0)

	}

}
