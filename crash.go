package core
import (
	"fmt"
	"os"
	"strings"
)
const (
	// error codes
	ExitWithError = 1
	ExitSuccess   = 0
)


//Crash prints errors and exit with code 2. First line is printed in bold red
func Crash(err error) {
	lines := strings.Split(err.Error(), "\n")
	if len(lines) > 0 {
		//"\033[31;1;4m turn color red and bold. \033[0m reset colors"
		fmt.Fprintf(os.Stderr, "\033[31;1m%s\n\033[0m", lines[0])
		for i := 1; i < len(lines); i++ {
			fmt.Fprintf(os.Stderr, "%s\n", lines[i])
		}
	}
	os.Exit(ExitWithError)
}

//Crashf prints errors and exit with code 2. First line is printed in bold red
func Crashf(format string, a ...interface{}) {
	err := fmt.Sprintf(format, a...)
	fmt.Fprintf(os.Stderr, "\033[31;1m%s\n\033[0m", err)
	os.Exit(ExitWithError)
}