package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/legionus/godep-utils/pkg/godeps"
)

var (
	hooksDir  = flag.String("hooks-dir", "", "Specifies directory with hooks")
	output    = flag.String("output", "-", "Write new Godeps to file")
	showMerge = flag.Bool("show-merge", false, "Show only merge")
)

func printJSON(fd *os.File, v interface{}) {
	out, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(fd, "%s\n", string(out))
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: godep-merge [options] Godeps-old.json Godeps-new.json\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "Options:\n")
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "Use the following hooks to change dependencies:\n")
	fmt.Fprintf(os.Stderr, "   <hooks-dir>/pre.sh  -- runs before merge\n")
	fmt.Fprintf(os.Stderr, "   <hooks-dir>/post.sh -- runs after deps merge\n")
	fmt.Fprintf(os.Stderr, "   <hooks-dir>/dep.sh  -- runs for every dep\n")
	fmt.Fprintf(os.Stderr, "\n")
	os.Exit(2)
}

func isAccessable(file string) bool {
	fd, err := os.Open(file)
	if err != nil {
		return false
	}
	fd.Close()
	return true
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

	var OutputFile *os.File

	if *output != "-" {
		OutputFile, err = os.OpenFile(*output, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		OutputFile = os.Stdout
	}

	oldGodeps, err := godeps.Parse(args[0])
	if err != nil {
		log.Fatal(err)
	}

	newGodeps, err := godeps.Parse(args[1])
	if err != nil {
		log.Fatal(err)
	}

	diff := godeps.MakeDiff(oldGodeps, newGodeps)

	if *showMerge {
		printJSON(OutputFile, diff)
		os.Exit(0)
	}

	hooks := &godeps.MergeHooks{}

	if *hooksDir != "" {
		if isAccessable(*hooksDir + "/godep-merge-pre-hook") {
			hooks.PreHook = *hooksDir + "/godep-merge-pre-hook"
		}
		if isAccessable(*hooksDir + "/godep-merge-post-hook") {
			hooks.PostHook = *hooksDir + "/godep-merge-post-hook"
		}
		if isAccessable(*hooksDir + "/godep-merge-dep-hook") {
			hooks.DepHook = *hooksDir + "/godep-merge-dep-hook"
		}
	}

	deps, err := godeps.Merge(diff, hooks)
	if err != nil {
		log.Fatal(err)
	}

	printJSON(OutputFile, deps)
}
