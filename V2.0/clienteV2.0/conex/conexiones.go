// modulo que maneja la conexion tcp de la herramienta
package conexiones

import (
	color "comando/colores"
	remoto "comando/conex/comandos"
	"errors"
	"fmt"
	"net"
	"os"
	"time"
)

/*
en caso de algun error se llama a esta funcion para reiniciar la conexion y no arrastrar errores

es una solucion que encontre para no saturar el programa y que se sigan generando errores
*/
func Reconexion(net net.Conn, ip string, tiempo time.Duration) {
	remoto.Borrar_consola()
	fmt.Println("reconectando...")
	error := net.Close()
	if error != nil {
		fmt.Println("error fatal: ", error)
		os.Exit(1)
	} else {
		Conexion(ip, tiempo)
	}

}

// funcion que se encarga de establecer conexion TCP con el host
func Conexion(ip string, tiempo time.Duration) error {
	conec, error := net.DialTimeout("tcp", ip, time.Second*tiempo)

	if error != nil { // si hay algun error

		return errors.New("\n[!]hubo un error al establecer conexion")

	} else {
		fmt.Printf(color.F_violeta+"[#] conexion establecida %s --> %s\n\r\r"+color.Reset, conec.LocalAddr(), conec.RemoteAddr())

		err := remoto.Comando(conec)
		if err != nil {
			fmt.Println(err)
			Reconexion(conec, ip, tiempo)
		}
	}
	return nil

}
