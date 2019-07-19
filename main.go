package main

import (
	"bufio"
	"fmt"
	"hash/fnv"
	"os"
	"path/filepath"
)

var progname = filepath.Base(os.Args[0])

func main() {
	h := fnv.New64a()
	e := make(map[uint64]struct{})
	s := bufio.NewScanner(os.Stdin)

	for s.Scan() {
		b := s.Bytes()
		_, err := h.Write(b)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", progname, err)
			continue
		}
		d := h.Sum64()
		h.Reset()
		if _, ok := e[d]; !ok {
			fmt.Println(string(b))
			e[d] = struct{}{}
		}
	}

	if err := s.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", progname, err)
		os.Exit(1)
	}
}