package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	lhvm "github.com/tomnomnom/leedshack2018"
)

func main() {
	var debug bool
	flag.BoolVar(&debug, "d", false, "debug mode")

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

	v.Log = ioutil.Discard
	if debug {
		lf, err := os.Create("debug.log")
		if err == nil {
			v.Log = lf
		}
	}

	v.Run(code, entry)
}
