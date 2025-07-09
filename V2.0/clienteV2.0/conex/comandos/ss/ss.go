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
	"io"
	"net"
	"os"
)

const (
	BUFFER_TAMAÑO = 8
)

func Escribir_img(img_byte []byte, nombre_arch string) {
	arch, arch_error := os.Create(fmt.Sprintf("%s.png", nombre_arch))
	if arch_error != nil {
		fmt.Println("error al escribir imagen: ", arch_error)

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

	_, escritura_error := conn.Write([]byte("ss"))
	if escritura_error != nil {
		return nil, errors.New("[!] error al obtener la imagen")

	}
	_, err := io.ReadFull(conn, tamaño_img) // lectura del tamaño de la imagen
	if err != nil {
		return nil, errors.New("error en la lectura del tamaño")

	}
	total := binary.BigEndian.Uint64(tamaño_img) //tamaño total de la imagen

	fmt.Printf("cantidad total de la imagen: %f MB aprox\n", float32(total)/1_000_000)

	img := make([]byte, total) // crear buffer de la imagen

	_, erro := io.ReadFull(conn, img) // leer toda la imagen

	if erro != nil {
		return nil, errors.New("[!] error al obtener la imagen")

	}
	fmt.Println("[*] se obtuvo la imagen")
	return img, nil

}
