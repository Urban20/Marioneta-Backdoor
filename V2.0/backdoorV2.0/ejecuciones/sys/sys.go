/*
	este modulo alberga las todas las llamadas al sistema que son necesarias

* funciones de bajo nivel (win api)
*/
package sistema

import (
	"syscall"
	"unsafe"
)

// crea un popup con un mensaje en windows
func MsgCartel(titulo string, msg string) error {
	sys, loaderr := syscall.LoadDLL("User32.dll")
	if loaderr != nil {
		return loaderr
	}

	sysfind, finderr := sys.FindProc("MessageBoxW")

	if finderr != nil {
		return finderr
	}

	const LOGO_WARNING = 0x00000030

	sysfind.Call(
		0,
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(msg))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(titulo))), // transformo string a utf16 y devuelvo el puntero, lo trato con usafe y lo transformo en uintprt (valor solicitado)
		LOGO_WARNING,
	)

	return nil

}
