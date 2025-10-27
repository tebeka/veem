package veem

import (
	"cmp"
	"errors"
)

type VM struct {
	stack []int
	sp    int
	err   error
}

func (vm *VM) Execute(code []Inst) (int, error) {
	ip := 0

	for ip < len(code) {
		ip += code[ip](vm)
		if vm.err != nil {
			return 0, vm.err
		}
	}

	if vm.sp == 0 {
		return 0, nil
	}

	return vm.stack[vm.sp-1], nil
}

type Inst func(*VM) int

func Push(n int) Inst {
	return func(vm *VM) int {
		vm.stack = append(vm.stack, n)
		vm.sp++
		return 1
	}
}

var ErrStackUnderflow = errors.New("stack underflow")
var ErrDivisionByZero = errors.New("division by zero")

func binOp(vm *VM, fn func(int, int) int) int {
	if vm.sp < 2 {
		vm.err = ErrStackUnderflow
		return 0
	}

	vm.stack[vm.sp-2] = fn(vm.stack[vm.sp-2], vm.stack[vm.sp-1])
	vm.sp--
	vm.stack = vm.stack[:vm.sp]
	return 1
}

func Add(vm *VM) int {
	return binOp(vm, func(a, b int) int { return a + b })
}

func Sub(vm *VM) int {
	return binOp(vm, func(a, b int) int { return a - b })
}

func Mul(vm *VM) int {
	return binOp(vm, func(a, b int) int { return a * b })
}

func Div(vm *VM) int {
	return binOp(
		vm,
		func(a, b int) int {
			if b == 0 {
				vm.err = ErrDivisionByZero
				return 0
			}
			return a / b
		},
	)
}

func Mod(vm *VM) int {
	return binOp(
		vm,
		func(a, b int) int {
			if b == 0 {
				vm.err = ErrDivisionByZero
				return 0
			}
			return a % b
		},
	)
}

func Cmp(vm *VM) int {
	return binOp(vm, cmp.Compare)
}
