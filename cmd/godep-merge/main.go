package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/legionus/godep-utils/pkg/godeps"
)

const (
	PreHookEnv  = "GODEP_MERGE_PRE_HOOK"
	PostHookEnv = "GODEP_MERGE_POST_HOOK"
	DepsHookEnv = "GODEP_MERGE_DEPS_HOOK"
)

var (
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
	fmt.Fprintf(os.Stderr, "Use the following variables to change dependencies:\n")
	fmt.Fprintf(os.Stderr, "   %s=<path-to-script>\n", PreHookEnv)
	fmt.Fprintf(os.Stderr, "   %s=<path-to-script>\n", PostHookEnv)
	fmt.Fprintf(os.Stderr, "   %s=<path-to-script>\n", DepsHookEnv)
	fmt.Fprintf(os.Stderr, "\n")
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

	deps, err := godeps.Merge(diff,
		&godeps.MergeHooks{
			PreHook:  os.Getenv(PreHookEnv),
			DepHook:  os.Getenv(DepsHookEnv),
			PostHook: os.Getenv(PostHookEnv),
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	printJSON(OutputFile, deps)
}
