/*
	este modulo alberga las todas las llamadas al sistema que son necesarias

* funciones de bajo nivel (win api)
*/
package sistema

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

const (
	LOGO_WARNING = 0x00000030
	NOTIFICACION = 0x00200000
	DLL          = "User32.dll"
	PROC         = "MessageBoxW"
)

// crea un popup con un mensaje en windows
func MsgCartel(titulo string, msg string) error {

	sys, loaderr := windows.LoadDLL(DLL)

	if loaderr != nil {
		return loaderr
	}

	defer func() error { // liberar la dll
		relerr := sys.Release()
		if relerr != nil {
			return relerr
		}
		return nil
	}()

	sysfind, procerr := sys.FindProc(PROC)
	if procerr != nil {
		return procerr
	}

	sysfind.Call(
		0,
		uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(msg))),
		uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(titulo))), // transformo string a utf16 y devuelvo el puntero, lo trato con usafe y lo transformo en uintprt (valor solicitado)
		LOGO_WARNING|NOTIFICACION,
	)

	return nil

}
