/*
ss -> screenshot

en este modulo se maneja:

* la logica para capturar el screenshot proveniente del host

* la escritura de la imagen
*/
package ss

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"os"
)

const (
	BUFFER_IMG    = 1_000_000
	BUFFER_TAMAÑO = 8
)

func Escribir_img(img_byte []byte, nombre_arch string) {
	arch, error := os.Create(fmt.Sprintf("%s.png", nombre_arch))
	if error != nil {
		fmt.Println("error al escribir imagen: ", error)

	} else {
		arch.Write(img_byte)
		arch.Close()
	}
}

/*
funcion encargada de la obtencion de los bytes de la imagen proveniente del host

obtener_img : net.conn -> byte

1. obtener el tamaño de la imagen

2. obtener la imagen por partes hasta completarla
*/
func Obtener_img(conn net.Conn) ([]byte, error) {
	// retornar bytes
	tamaño_img := make([]byte, BUFFER_TAMAÑO)
	img := make([]byte, BUFFER_IMG)
	contador := 0

	_, error := conn.Write([]byte("ss"))
	if error != nil {
		return nil, errors.New("[!] error al obtener la imagen")

	}
	n, error := conn.Read(tamaño_img) // lectura del tamaño de la imagen
	if error != nil {
		return nil, errors.New("error en la lectura del tamaño")

	}
	total := binary.BigEndian.Uint64(tamaño_img[:n]) //tamaño total de la imagen

	fmt.Println("cantidad total de la imagen: ", total)
	for {
		n_imagen, error := conn.Read(img[contador:])

		if error != nil {
			return nil, errors.New("[!] error al obtener la imagen")

		}
		contador += n_imagen

		fmt.Println("cantidad obtenida: ", contador)

		if contador >= int(total) {
			fmt.Println("[*] se obtuvo la imagen")
			return img, nil
		}

	}

}
