package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/tomnomnom/xtermcolor"
)

const (
	cols = 80
	rows = 40
)

const (
	escape = "\x1b"
	char   = "█"
)

const (
	PUSH = iota
	ADD
	SUBT
	SLEEP
	JMP
	JMPLT
	JMPGT
	JMPEQ
	CALL
	RET
	ARG
	PRINT
	SETPX
	RECT
	PAINT
	CLEAR
	LD
	ST
	HALT
)

type vm struct {
	code []int
	pc   int

	stack []int
	sp    int

	fp int

	ram  []int
	vram []uint8
}

func (v *vm) initVram() {
	v.vram = make([]uint8, cols*rows)

	for i := 0; i < cols*rows; i++ {
		v.vram[i] = 236
	}
}

func (v *vm) setPixel(x, y, c int) {
	offset := (y * cols) + x
	v.vram[offset] = xtermcolor.FromInt(uint32(c))
}

func (v *vm) paint() {
	buf := &bytes.Buffer{}
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			offset := (y * cols) + x
			c := v.vram[offset]
			buf.WriteString(fmt.Sprintf("%s[38;5;%dm%s", escape, c, char))
		}

		buf.WriteRune('\n')
	}

	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()

	os.Stdout.Write(buf.Bytes())
}

func (v *vm) run(code []int, pc int) {
	v.code = code
	v.pc = pc

	v.stack = make([]int, 128)
	v.sp = -1

	v.fp = -1

	v.ram = make([]int, 1024)
	v.initVram()

	for {
		switch v.nextOp() {

		case PUSH:
			v.push(v.nextOp())

		case ADD:
			v.push(v.pop() + v.pop())

		case SUBT:
			b := v.pop()
			a := v.pop()
			v.push(a - b)

		case SLEEP:
			time.Sleep(time.Duration(v.nextOp() * 1000000))

		case JMP:
			addr := v.nextOp()
			v.pc = addr

		case JMPLT:
			comp := v.nextOp()
			addr := v.nextOp()

			if v.peek() < comp {
				v.pc = addr
			}

		case JMPGT:
			comp := v.nextOp()
			addr := v.nextOp()

			if v.peek() > comp {
				v.pc = addr
			}

		case JMPEQ:
			comp := v.nextOp()
			addr := v.nextOp()

			if v.peek() == comp {
				v.pc = addr
			}

		case CALL:
			addr := v.nextOp()
			nargs := v.nextOp()

			v.push(nargs, v.fp, v.pc)

			v.fp = v.sp
			v.pc = addr

		case RET:
			val := v.pop()

			// reset the stack
			v.sp = v.fp

			v.pc = v.pop()
			v.fp = v.pop()
			nargs := v.pop()

			// ditch the args
			for i := 0; i < nargs; i++ {
				v.pop()
			}

			// put the return val at the top of the stack
			v.push(val)

		case ARG:
			offset := (v.fp - v.nextOp()) - 3
			v.push(v.stack[offset])

		case PRINT:
			fmt.Println(v.pop())

		case SETPX:
			c := v.pop()
			y := v.pop()
			x := v.pop()
			v.setPixel(x, y, c)

		case RECT:
			color := v.pop()
			h := v.pop()
			w := v.pop()
			y := v.pop()
			x := v.pop()

			for cY := y; cY < (y + h); cY++ {
				for cX := x; cX < (x + w); cX++ {
					v.setPixel(cX, cY, color)
				}
			}

		case PAINT:
			v.paint()

		case CLEAR:
			v.initVram()

		case LD:
			addr := v.nextOp()
			v.push(v.ram[addr])

		case ST:
			addr := v.nextOp()
			v.ram[addr] = v.pop()

		case HALT:
			return
		}
	}

}

func (v *vm) nextOp() int {
	op := v.code[v.pc]
	v.pc++
	return op
}

func (v *vm) push(nn ...int) {
	for _, n := range nn {
		v.sp++
		v.stack[v.sp] = n
	}
}

func (v *vm) pop() int {
	val := v.stack[v.sp]
	v.sp--
	return val
}

func (v *vm) peek() int {
	return v.stack[v.sp]
}

func main() {
	code := []int{

		// Initial state
		PUSH, 2, // x
		ST, 0,
		PUSH, 2, // y
		ST, 1,

		LD, 0, // x
		LD, 1, // y
		PUSH, 20, // w
		PUSH, 10, // h
		PUSH, 0xff000000, // color
		RECT,

		PAINT,
		CLEAR,

		// new coords
		LD, 0,
		PUSH, 3,
		ADD,
		ST, 0,

		LD, 1,
		PUSH, 1,
		ADD,
		ST, 1,

		SLEEP, 100,
		JMP, 8,
		HALT,
	}

	v := &vm{}
	v.run(code, 0)
}
