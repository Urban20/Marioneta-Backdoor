// modulo que controla el input de la consola
package input

import (
	"bufio"
	"fmt"
	"os"
)

// funcionamiento similar a funcion input de python
func Input(msg string) string {
	fmt.Print(msg)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}
