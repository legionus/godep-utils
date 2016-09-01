package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/legionus/godep-utils/pkg/godeps"
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: godep-inspect Godeps.json script\n")
	fmt.Fprintf(os.Stderr, "\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	var err error

	flag.Usage = usage
	flag.Parse()

	args := flag.Args()

	if len(args) < 2 {
		fmt.Fprintf(os.Stderr, "ERROR: More arguments required\n")
		usage()
	}

	deps, err := godeps.Parse(args[0])
	if err != nil {
		log.Fatal(err)
	}

	if err := godeps.Inspect(deps, args[1]); err != nil {
		log.Fatal(err)
	}
}
