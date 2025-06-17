// modulo que maneja la conexion tcp de la herramienta
package conexiones

import (
	color "comando/colores"
	remoto "comando/conex/comandos"
	"errors"
	"fmt"
	"net"
	"os"
	"os/exec"
	"time"
)

const TIMEOUT = 3

var timeout_err net.Error

/*
	en caso de algun error se llama a esta funcion para reiniciar la conexion y no arrastrar errores

es una solucion que encontre para no saturar el programa y que se sigan generando errores
*/
func Reconexion(net net.Conn, ip string, tiempo time.Duration) {
	exec.Command("powershell", "-command", "clear").Run()
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
func Conexion(ip string, tiempo time.Duration) {
	conec, error := net.DialTimeout("tcp", ip, time.Second*tiempo)

	fmt.Printf(color.F_violeta+"[#] conexion establecida %s --> %s\n"+color.Reset, conec.LocalAddr(), conec.RemoteAddr())

	if errors.As(error, &timeout_err) && timeout_err.Timeout() {
		fmt.Println(color.Rojo + "\n[!]tiempo agotado\n" + color.Reset)

	} else if error != nil { // si hay algun error
		fmt.Printf(color.Rojo+"\n[!] error:\n%s \n", error.Error()+color.Reset)

	} else {
		err := remoto.Comando(conec)
		if err != nil {
			fmt.Println(err)
			Reconexion(conec, ip, tiempo)
		}
	}

}
