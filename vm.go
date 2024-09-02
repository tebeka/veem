package main

import (
	"fmt"
)

type Op uint16 // low byte op, high byte value

const (
	InvalidOp Op = iota
	PushOp
	AddOp
	SubOp
)

type VM struct {
	stack []uint8
}

func (v *VM) Execute(program []Op) (uint8, error) {
	for i, op := range program {
		code := op & 0xFF
		switch code {
		case InvalidOp:
			return 0, fmt.Errorf("%d: invalid op", i)
		case PushOp:
			val := uint8(op >> 8)
			v.push(val)
		case AddOp:
			v1, v2 := v.pop(), v.pop()
			val := v1 + v2
			v.push(val)
		case SubOp:
			v1, v2 := v.pop(), v.pop()
			val := v1 - v2
			v.push(val)
		}
	}

	if len(v.stack) > 0 {
		return v.stack[len(v.stack)-1], nil
	}

	return 0, nil
}

func (v *VM) push(val uint8) {
	v.stack = append(v.stack, val)
}

func (v *VM) pop() uint8 {
	i := len(v.stack) - 1
	val := v.stack[i]
	v.stack = v.stack[:i]
	return val
}

func main() {
	program := []Op{
		PushOp | 3<<8,
		PushOp | 5<<8,
		PushOp | 8<<8,
		SubOp,
		AddOp,
	}
	var v VM
	val, err := v.Execute(program)
	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}
	fmt.Println(val)
}
