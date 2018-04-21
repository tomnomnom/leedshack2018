package main

import "fmt"

const (
	PUSH = iota
	ADD
	SUBT
	JMPLT
	JMPGT
	JMPEQ
	PRINT
	HALT
)

type vm struct {
	code []int
	pc   int

	stack []int
	sp    int
}

func (v *vm) run(code []int) {
	v.code = code
	v.pc = 0

	v.stack = make([]int, 128)
	v.sp = -1

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

func (v *vm) push(n int) {
	v.sp++
	v.stack[v.sp] = n
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
		PUSH, 2,
		PUSH, 2,
		ADD,
		PUSH, 1,
		SUBT,
		JMPEQ, 3, 2,
		PRINT,
		HALT,
	}

	v := &vm{}
	v.run(code)
}
