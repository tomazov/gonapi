package main

import (
	"fmt"
)

// This file provides a very small program that can be used to verify
// that the Go toolchain and project dependencies are configured
// correctly. Running `go run install.go` should print "installation OK".

func main() {
	fmt.Println("installation OK")
}
