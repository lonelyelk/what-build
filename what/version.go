package what

import "fmt"

// Version string assigned by LDFLAGS on build time
var Version = "0.5.1"

// PrintVersion outputs tool version
func PrintVersion() {
	fmt.Printf("what-build version %s\n", Version)
}
