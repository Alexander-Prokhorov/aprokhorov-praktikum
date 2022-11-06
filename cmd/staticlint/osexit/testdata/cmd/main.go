package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("Hello World")
	os.Exit(0) // want "os.Exit used in func main of package main"
}
