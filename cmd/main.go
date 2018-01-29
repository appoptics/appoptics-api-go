//
// This package exists only to print the version number
//
package main

import (
	"fmt"

	"github.com/appoptics/appoptics-api-go"
)

func main() {
	fmt.Printf(printVersion())
}

func printVersion() string {
	return fmt.Sprintf(appoptics.Version())
}
