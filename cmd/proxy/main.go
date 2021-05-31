package main

import (
	"fmt"
	"os"

	"github.com/crochee/proxy-go/cmd"
)

func main() {
	if err := cmd.Server(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	os.Exit(0)
}
