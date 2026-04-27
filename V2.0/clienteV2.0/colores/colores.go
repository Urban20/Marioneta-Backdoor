// colores y apartado visual del programa
package color

import (
	"fmt"
)

const LOGO = ` 
_  _ ____ ____ _ ____ _  _ ____ ___ ____ 
|\/| |__| |__/ | |  | |\ | |___  |  |__| 
|  | |  | |  \ | |__| | \| |___  |  |  | 
                                                                        
-----------------------------------------                                                                                                                  
 cliente V 2.0  - por Urb@n  
-----------------------------------------                           
`

func rgb(r, g, b int, fondo bool) string {

	var d int

	if fondo {

		d = 48

	} else {

		d = 38
	}

	return fmt.Sprintf("\033[%d;2;%d;%d;%dm", d, r, g, b)
}

// colores de terminal que se usan en los outputs de la consola
var Reset = "\033[0m"
var F_violeta = rgb(134, 72, 232, true)
var Amarillo = "\033[0;93m"
var Verde = "\033[0;32m"
var Rojo = "\033[0;31m"
var Violeta = rgb(134, 72, 232, false)
var Seleccion = "\033[30;47m"
