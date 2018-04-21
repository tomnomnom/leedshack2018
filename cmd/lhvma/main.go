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

	sections := map[string][]string{}
	section := ""

	sc := bufio.NewScanner(r)
	for sc.Scan() {
		if strings.TrimSpace(sc.Text()) == "" {
			continue
		}

		// instruction
		if strings.HasPrefix(sc.Text(), "\t") {
			tokens := strings.Split(strings.TrimSpace(sc.Text()), " ")
			sections[section] = append(sections[section], tokens...)
			continue
		}

		// section def
		section = strings.Trim(sc.Text(), ": ")
		sections[section] = make([]string, 0)
	}

	// do the memory layout so we know the addresses
	layout := []string{}
	addrs := map[string]int{}
	pc := 0
	for s, toks := range sections {

		// get the addresses
		addrs[s] = pc

		for _, t := range toks {
			layout = append(layout, t)
			pc++
		}
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
