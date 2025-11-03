package veem

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

type Number float64

type VM struct {
	stack []Number
	err   error
}

func (vm *VM) Pop() Number {
	n := len(vm.stack)
	if n == 0 {
		vm.err = errors.New("stack underflow")
		return 0
	}

	v := vm.stack[n-1]
	vm.stack = vm.stack[:n-1]
	return v
}

func (vm *VM) Push(v Number) {
	vm.stack = append(vm.stack, v)
}

func cleanCode(code string) string {
	i := strings.Index(code, ";")
	if i != -1 {
		code = code[:i]
	}

	return strings.ToUpper(strings.TrimSpace(code))
}

const (
	OpAdd = "ADD"
	OpSub = "SUB"
	OpMul = "MUL"
	OpDiv = "DIV"
	OpMod = "MOD"

	OpPush = "PUSH"
)

func (vm *VM) Execute(code []string) {
	ip := 0

	for ip < len(code) && vm.err == nil {
		line := cleanCode(code[ip])

		if line == "" {
			ip++
			continue
		}

		fields := strings.Fields(line)
		switch op := fields[0]; op {
		case OpAdd, OpSub, OpMul, OpDiv, OpMod:
			if len(fields) != 1 {
				vm.err = fmt.Errorf("bad inst: %q", code)
				return
			}

			binOp(vm, op)
			ip++
		case OpPush:
			if len(fields) != 2 {
				vm.err = fmt.Errorf("bad inst: %q", code)
				return
			}

			n, err := strconv.ParseFloat(fields[1], 64)
			if err != nil {
				vm.err = fmt.Errorf("%s: %w", code, err)
				return
			}
			vm.Push(Number(n))
			ip++
		}

	}
}

func binOp(vm *VM, op string) {
	b, a := vm.Pop(), vm.Pop()
	if vm.err != nil {
		return
	}

	var n Number
	switch op {
	case OpAdd:
		n = a + b
	case OpSub:
		n = a - b
	case OpMul:
		n = a * b
	case OpDiv:
		n = a / b
	case OpMod:
		if math.Round(float64(a)) != float64(a) || math.Round(float64(b)) != float64(b) {
			vm.err = fmt.Errorf("MOD not on round numbers")
			return
		}

		if b == 0 {
			vm.err = errors.New("mod by 0")
		} else {
			n = Number(int(a) % int(b))
		}
	}

	vm.Push(n)
}
