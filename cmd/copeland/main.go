package main

import (
	"flag"
	"fmt"
	"os"
)

const (
	ProgName = "copeland"
)

var (
	version = "undefined"

	showVersion = flag.Bool("version", false, "show program version and exit")
)

func cmdVersion() int {
	fmt.Println(version)
	return 0
}

func usage() {
	out := flag.CommandLine.Output()
	fmt.Fprintln(out, "Usage:")
	fmt.Fprintln(out)
	fmt.Fprintf(out, "%s [OPTION]... [BALLOT FILE]...\n", ProgName)
	fmt.Fprintln(out)
	fmt.Fprintln(out, "  Process files with preference lists and output Copeland's ranking.")
	fmt.Fprintln(out)
	fmt.Fprintln(out, "Options:")
	flag.PrintDefaults()
}

func run() int {
	flag.CommandLine.Usage = usage
	flag.Parse()
		
	if *showVersion {
		fmt.Println(version)
		return 0
	}

	args := flag.Args()
	if len(args) == 0 {
		usage()
		return 2
	}

	return 0
}

func main() {
	os.Exit(run())
}
