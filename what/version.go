package what

import "fmt"

// Version string assigned by LDFLAGS on build time
var Version = "0.1.0"

// PrintVersion outputs tool version
func PrintVersion() {
	fmt.Printf("what-build version %s\n", Version)
}
