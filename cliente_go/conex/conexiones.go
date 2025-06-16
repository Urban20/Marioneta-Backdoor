// modulo que maneja la conexion tcp del comando
package conexiones

import (
	color "comando/colores"
	remoto "comando/conex/comandos"
	"errors"
	"fmt"
	"net"
	"time"
)

const TIMEOUT = 3

var timeout_err net.Error

// funcion que se encarga de establecer conexion TCP con el host
func Conexion(ip string, tiempo time.Duration) {
	conec, error := net.DialTimeout("tcp", ip, time.Second*tiempo)

	if errors.As(error, &timeout_err) && timeout_err.Timeout() {
		fmt.Println(color.Rojo + "\n[!]tiempo agotado\n" + color.Reset)

	} else if error != nil { // si hay algun error
		fmt.Printf(color.Rojo+"\n[!] error:\n%s \n", error.Error()+color.Reset)

	} else {
		remoto.Comando(conec)
	}

}
