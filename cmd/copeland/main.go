package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"slices"
	"strings"
)

const (
	ProgName = "copeland"
)

var (
	version = "undefined"

	showVersion   = flag.Bool("version", false, "show program version and exit")
	normalizeCase = flag.Bool("normalize-case", true, "normalize case")
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

	uniqNames := make(map[string]struct{})
	for _, filename := range args {
		if err := func() error {
			f, err := os.Open(filename)
			if err != nil {
				return fmt.Errorf("unable to open file %q: %w", filename, err)
			}
			defer f.Close()
			names, err := readBallot(f)
			if err != nil {
				return fmt.Errorf("ballot %q reading failed: %w", filename, err)
			}
			for _, name := range names {
				uniqNames[name] = struct{}{}
			}
			return nil
		}(); err != nil {
			log.Fatalf("error: %v", err)
		}
	}

	sortedNames := make([]string, 0, len(uniqNames))
	for name := range uniqNames {
		sortedNames = append(sortedNames, name)
	}
	slices.Sort(sortedNames)
	fmt.Println("Registered names:")
	for _, name := range sortedNames {
		fmt.Printf("\t%s\n", name)
	}

	return 0
}

func readBallot(input io.Reader) ([]string, error) {
	scanner := bufio.NewScanner(input)
	var res []string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		if normalizeCase != nil && *normalizeCase {
			line = strings.ToUpper(line)
		}
		res = append(res, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("line scanning failed: %w", err)
	}
	return res, nil
}

func main() {
	os.Exit(run())
}
