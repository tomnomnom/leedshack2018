package main

import "fmt"

const (
	PUSH = iota
	ADD
	SUB
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
		case SUB:
			b := v.pop()
			a := v.pop()
			v.push(a - b)
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

func main() {
	code := []int{
		PUSH, 2,
		PUSH, 2,
		ADD,
		PUSH, 1,
		SUB,
		PRINT,
		HALT,
	}

	v := &vm{}
	v.run(code)
}
