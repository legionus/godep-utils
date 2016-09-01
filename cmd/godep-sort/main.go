package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/legionus/godep-utils/pkg/godeps"
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: godep-sort Godeps.json\n")
	fmt.Fprintf(os.Stderr, "\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func printJSON(fd *os.File, v interface{}) {
	out, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(fd, "%s\n", string(out))
}

func main() {
	var err error

	flag.Usage = usage
	flag.Parse()

	args := flag.Args()

	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "ERROR: More arguments required\n")
		usage()
	}

	deps, err := godeps.Parse(args[0])
	if err != nil {
		log.Fatal(err)
	}

	sort.Sort(godeps.SortByImportPath(deps.Deps))

	printJSON(os.Stdout, deps)
}
