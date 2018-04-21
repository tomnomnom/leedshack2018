package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	lhvm "github.com/tomnomnom/leedshack2018"
)

func main() {
	flag.Parse()

	inf := flag.Arg(0)
	outf := flag.Arg(1)

	if inf == "" || outf == "" {
		fmt.Println("usage: lhvma <infile> <outfile>")
		return
	}

	r, err := os.Open(inf)
	if err != nil {
		fmt.Printf("failed opening input file: %s", err)
		return
	}
	defer r.Close()

	w, err := os.Create(outf)
	if err != nil {
		fmt.Printf("failed creating output file: %s", err)
		return
	}

	layout := []string{}
	addrs := map[string]int{}
	pc := 0
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		trimmed := strings.TrimSpace(sc.Text())
		if trimmed == "" || strings.HasPrefix(trimmed, "//") {
			continue
		}

		if strings.HasPrefix(sc.Text(), "\t") || strings.HasPrefix(sc.Text(), " ") {
			// instruction
			tokens := strings.Split(trimmed, " ")
			layout = append(layout, tokens...)
			pc += len(tokens)
			continue
		}

		// label
		label := strings.Trim(sc.Text(), ": ")
		addrs[label] = pc
	}

	// decode
	code := []int{}

	for _, t := range layout {

		// Check for OP
		if op, ok := lhvm.Ops[t]; ok {
			code = append(code, op)
			continue
		}

		// Check for label
		if addr, ok := addrs[t]; ok {
			code = append(code, addr)
			continue
		}

		// assume constant
		n, err := strconv.ParseInt(t, 0, 0)
		if err != nil {
			log.Fatalf("invalid constant: %s", err)
		}
		code = append(code, int(n))

	}

	// write to file
	fmt.Fprintf(w, "%d\n", addrs["main"])
	for _, c := range code {
		fmt.Fprintf(w, "%d\n", c)
	}

	fmt.Printf("%#v\n", code)

}
