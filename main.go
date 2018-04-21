package main

import "fmt"

const (
	PUSH = iota
	ADD
	SUBT
	JMPLT
	JMPGT
	JMPEQ
	CALL
	RET
	ARG
	PRINT
	HALT
)

type vm struct {
	code []int
	pc   int

	stack []int
	sp    int

	fp int
}

func (v *vm) run(code []int, pc int) {
	v.code = code
	v.pc = pc

	v.stack = make([]int, 128)
	v.sp = -1

	v.fp = -1

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
		// sub @ 0
		ARG, 0,
		ARG, 1,
		ADD,
		RET,

		// main
		PUSH, 2,
		PUSH, 2,
		CALL, 0, 2,
		PUSH, 1,
		SUBT,
		PRINT,
		HALT,
	}

	v := &vm{}
	v.run(code, 6)
}
