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

	"github.com/Snawoot/copeland"
)

const (
	ProgName = "copeland"
)

var (
	version = "undefined"

	showVersion   = flag.Bool("version", false, "show program version and exit")
	normalizeCase = flag.Bool("normalize-case", true, "normalize case")
	scoreWin      = flag.Float64("score-win", 1, "score for win against opponent")
	scoreTie      = flag.Float64("score-tie", .5, "score for tie against opponent")
	scoreLoss     = flag.Float64("score-loss", 0, "score for tie against opponent")
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

	cl, err := copeland.New(sortedNames)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Println()

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
			if err := cl.Update(names); err != nil {
				return fmt.Errorf("file %q: Copeland update failed: %w", filename, err)
			}
			return nil
		}(); err != nil {
			log.Fatalf("error: %v", err)
		}
	}

	cl.Dump()

	fmt.Println("Scores:")
	for _, entry := range cl.Score(&copeland.Scoring{
		Win:  *scoreWin,
		Tie:  *scoreTie,
		Loss: *scoreLoss,
	}) {
		fmt.Printf("\t%s\t%g\n", entry.Name, entry.Score)
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
