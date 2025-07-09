package conexiones

import (
	"errors"
	"fmt"
	"net"
	"strings"
)

// funcion que intenta obtener la direccion ipv4 de la maquina
func Ipv4() (string, error) {
	con, dial_error := net.Dial("udp", "10.255.255.255:2")
	if dial_error != nil {
		return "", errors.New("error al obtener direccion ipv4 (Ipv4() fallo)")
	}
	con.Close()
	return strings.Split(con.LocalAddr().String(), ":")[0], nil

}

// se pone en escucha la maquina, se retorna la proxima conexion en caso de exito
func Server(ip string, puerto string) (net.Listener, error) {

	conexion, escucha_error := net.Listen("tcp", fmt.Sprintf("%s:%s", ip, puerto))
	if escucha_error != nil {
		return nil, errors.New("error al levantar el server (server () fallo)")
	}

	return conexion, nil

}
