package main

import (
	"bufio"
	"fmt"
	"hash/fnv"
	"io"
	"os"

	"github.com/karrick/golf"
	"github.com/karrick/gorill"
)

func init() {
	// Rather than display the entire usage information for a parsing error,
	// merely allow golf library to display the error message, then print the
	// command the user may use to show command line usage information.
	golf.Usage = func() {
		fmt.Fprintf(os.Stderr, "Use `%s --help` for more information.\n", ProgramName)
	}
}

var (
	optHelp    = golf.BoolP('h', "help", false, "Print command line help and exit")
	optQuiet   = golf.BoolP('q', "quiet", false, "Do not print intermediate errors to stderr")
	optVerbose = golf.BoolP('v', "verbose", false, "Print verbose output to stderr")
)

func cmd() error {
	golf.Parse()

	if *optHelp {
		fmt.Printf("unique\n\n")
		fmt.Println(golf.Wrap("Like `uniq` but does not require input lines to be sorted."))
		fmt.Println(golf.Wrap("SUMMARY:  unique [options] [file1 [file2 ...]] [options]"))
		fmt.Println("EXAMPLES:")
		fmt.Println("\twho | unique")
		fmt.Println("\nCommand line options:")
		golf.PrintDefaults()
		return nil
	}

	h := fnv.New64a()
	e := make(map[uint64]struct{})

	var ior io.Reader
	if golf.NArg() == 0 {
		ior = os.Stdin
	} else {
		ior = &gorill.FilesReader{Pathnames: golf.Args()}
	}

	s := bufio.NewScanner(ior)

	for s.Scan() {
		b := s.Bytes()
		_, err := h.Write(b)
		if err != nil {
			warning("%s\n", err)
			continue
		}
		d := h.Sum64()
		h.Reset()
		if _, ok := e[d]; !ok {
			fmt.Println(string(b))
			e[d] = struct{}{}
		}
	}

	return s.Err()
}
