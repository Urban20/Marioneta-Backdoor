package consola

import (
	color "comando/colores"
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

func actualizar_seccion(n int) {

	for x := 0; x < n; x++ {

		fmt.Print("\033[F")

	}

}

func leer_tecla(i *int, tecla []byte) bool {
	os.Stdin.Read(tecla)
	flechas := tecla[2]

	if tecla[0] == 13 {

		return true

	} else if flechas == 65 {

		*i--
	} else if flechas == 66 {

		*i++
	}

	return false

}

// abstraigo la funcion para borrar consola
func Borrar_consola() {
	fmt.Println(strings.Repeat("\n", 50))
	fmt.Print("\033[H")
	Imprimir_logo()

}

func desplegar_opcion(opciones []string) string {

	var i int
	var op_largo = len(opciones)

	for {
		tecla := make([]byte, 3)

		for _, op := range opciones {

			if i > op_largo-1 {
				i = 0

			} else if i < 0 {
				i = op_largo - 1
			}

			if op == opciones[i] {
				fmt.Println(color.Seleccion + op + color.Reset)
			} else {
				fmt.Println(op)
			}
		}

		if leer_tecla(&i, tecla) {

			return opciones[i]
		}

		actualizar_seccion(op_largo)

	}
}
func Imprimir_logo() {
	fmt.Println(color.Violeta + color.LOGO + color.Reset)
}

func Menu(opciones []string) string {

	fd := int(os.Stdin.Fd())

	st, _ := term.MakeRaw(fd)

	defer term.Restore(fd, st)

	return desplegar_opcion(opciones)

}
