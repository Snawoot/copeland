package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
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
	names         = flag.String("names", "", "filename of list of names in the voting. If not specified names inferred from first ballot")
	skipErrors    = flag.Bool("skip-errors", false, "skip ballot errors, but still report them")
)

func cmdVersion() int {
	fmt.Println(version)
	return 0
}

func usage() {
	out := flag.CommandLine.Output()
	fmt.Fprintln(out, "Usage:")
	fmt.Fprintln(out)
	fmt.Fprintf(out, "%s [OPTION]... [BALLOT FILE OR DIRECTORY]...\n", ProgName)
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
	var nameListFilename string
	var err error
	if *names != "" {
		nameListFilename = *names
	} else {
		nameListFilename, err = getFirstFilename(args)
		if err != nil {
			log.Fatalf("unable to get first filename: %v", err)
		}
	}
	if err := func() error {
		f, err := os.Open(nameListFilename)
		if err != nil {
			return fmt.Errorf("unable to open file %q: %w", nameListFilename, err)
		}
		defer f.Close()
		names, err := readBallot(f)
		if err != nil {
			return fmt.Errorf("ballot %q reading failed: %w", nameListFilename, err)
		}
		for _, name := range names {
			uniqNames[name] = struct{}{}
		}
		return nil
	}(); err != nil {
		log.Fatalf("error: %v", err)
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

	if err := walkFiles(args, func(filename string) error {
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
			if *skipErrors {
				log.Printf("file %q: Copeland update failed: %v", filename, err)
				return nil
			}
			return fmt.Errorf("file %q: Copeland update failed: %w", filename, err)
		}
		return nil
	}); err != nil {
		log.Fatalf("file walk failed: %v", err)
	}

	fmt.Println("Scores:")
	for i, rank := range copeland.RankScore(cl.Score(&copeland.Scoring{
		Win:  *scoreWin,
		Tie:  *scoreTie,
		Loss: *scoreLoss,
	})) {
		fmt.Printf("\tRank %d:\n", i+1)
		for _, entry := range rank {
			fmt.Printf("\t\t%g\t%s\n", entry.Score, entry.Name)
		}
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

func getFirstFilename(args []string) (string, error) {
	var res string
	for _, root := range args {
		err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.Type().IsRegular() {
				res = path
				return fs.SkipAll
			}
			return nil
		})
		if err != nil {
			return "", err
		}
		if res != "" {
			return res, nil
		}
	}
	return "", errors.New("no files were found!")
}

func walkFiles(args []string, f func(path string) error) error {
	for _, root := range args {
		err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
			if d.Type().IsRegular() {
				return f(path)
			}
			return nil
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	log.Default().SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)
	log.Default().SetPrefix(strings.ToUpper(ProgName) + ": ")
	os.Exit(run())
}
