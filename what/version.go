package what

import "fmt"

// Version string assigned by LDFLAGS on build time
var Version string

// PrintVersion outputs tool version
func PrintVersion() {
	fmt.Printf("what-build version %s\n", Version)
}
