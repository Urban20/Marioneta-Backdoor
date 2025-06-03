
## Descripcion
Es una uilidad la cual se accede mediante parametros en consola (SE ENCUENTRAN EN LA SECCION USO) que nos permite acceder remotamente a una computadora a traves de la red privada (no esta pensado para el control remoto a traves de internet)

Consta de dos programas independientes: el cliente y el backdoor
* el backdor no tiene interfaz, intenta copiarse automaticamente al startup de windows si es la primera vez que se ejecuta, si falla y esta en modo sin consola (.pyw) se hara un error silencioso, sino se imprime en consola
* el cliente es el script utilizado para conectarse via TCP al backdoor y asi obtener acceso remoto a la maquina a traves de comandos de consola 

## uso:

1. ejecutar el backdoor (.pyw)
2. ejecutar el cliente de la siguiente manera (se debe tener python o se debe compilar el script a .exe):
- abrir un terminal en la ubicacion del repositorio
- ejecucion: python cliente.py --ip [ip de la maquina] -P [puerto] (utiliza el puerto 999, si deseas podes cambiarlo manualmente en el codigo)

## vista previa del cliente, script que se utiliza para acceder
<p align="center">
  <img src="https://i.postimg.cc/Z5rNhM02/demo.png" alt="Descripción" width="500"/>
</p>

la interfaz imprime algunos comandos basicos:
* enviar un mensaje: automatiza el comando msg * ...
(solo disponible en algunas ediciones de Windows)

* apagar: automatiza el comando shutdown

* borrar script: automatiza un cls para limpiar consola en el cliente

* ss: se le envia una señal al backdoor para que tome captura de pantalla, luego esa imagen se codifica a binario y se transmite por la red a traves de sockets, si todo sale bien se deberia ver una screenshot de la maquina victima en el directorio del cliente con el nombre de "screen.jpg"
* ademas se puede ejecutar casi cualquier comando, eso solo son algunos comandos automaticos pero podriamos ejecutar un whoami, un ls, dir, cd , etc
