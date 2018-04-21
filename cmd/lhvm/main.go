package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	lhvm "github.com/tomnomnom/leedshack2018"
)

func main() {
	flag.Parse()

	fn := flag.Arg(0)
	if fn == "" {
		fmt.Println("usage: lhvm <executable>")
		return
	}

	f, err := os.Open(fn)
	if err != nil {
		log.Fatalf("error reading file: %s", err)
	}
	defer f.Close()

	code := []int{}
	sc := bufio.NewScanner(f)

	// entrypoint is the first line
	sc.Scan()
	entry, err := strconv.Atoi(sc.Text())
	if err != nil {
		log.Fatalf("entrypoint conv error: %s", err)
	}

	for sc.Scan() {
		op, err := strconv.Atoi(sc.Text())
		if err != nil {
			log.Fatalf("conv error: %s", err)
		}
		code = append(code, op)
	}

	v := &lhvm.VM{}

	v.Run(code, entry)
}
