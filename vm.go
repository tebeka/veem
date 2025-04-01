package veem

import "errors"

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

func Add(vm *VM) int {
	if vm.sp < 2 {
		vm.err = ErrStackUnderflow
		return 0
	}

	vm.stack[vm.sp-2] += vm.stack[vm.sp-1]
	vm.sp--
	return 1
}

func Sub(vm *VM) int {
	if vm.sp < 2 {
		vm.err = ErrStackUnderflow
		return 0
	}

	vm.stack[vm.sp-2] -= vm.stack[vm.sp-1]
	vm.sp--
	return 1
}
