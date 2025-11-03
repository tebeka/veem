package veem

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/yaml/go-yaml"
)

func TestVM_PushPop(t *testing.T) {
	var vm VM
	vm.Push(1)
	vm.Push(2)
	vm.Push(3)

	if !slices.Equal([]Number{1, 2, 3}, vm.stack) {
		t.Fatal(vm.stack)
	}

	n := vm.Pop()
	if n != 3 || vm.err != nil {
		t.Fatal(n, vm.err)
	}

	vm.Push(4)

	expected := []Number{1, 2, 4}
	if !slices.Equal(expected, vm.stack) {
		t.Fatal(vm.stack)
	}

	slices.Reverse(expected)
	for _, v := range expected {
		n := vm.Pop()
		if n != v || vm.err != nil {
			t.Fatal(n, vm.err)
		}
	}

	_ = vm.Pop()
	if vm.err == nil {
		t.Fatal("no error")
	}
}

type testCase struct {
	Name  string
	Code  string
	Out   Number
	Error bool
}

func loadCases(t *testing.T, fileName string) []testCase {
	data, err := os.ReadFile(fileName)
	if err != nil {
		t.Fatal(err)
		return nil
	}

	var cases []testCase
	if err := yaml.Unmarshal(data, &cases); err != nil {
		t.Fatal(err)
		return nil
	}

	return cases
}

func TestVM_Execute(t *testing.T) {
	files, err := filepath.Glob("testdata/*.yml")
	if err != nil {
		t.Fatal(err)
	}

	for _, fileName := range files {
		suite := filepath.Base(fileName)
		suite = suite[:len(suite)-4]
		cases := loadCases(t, fileName)
		for _, tc := range cases {
			name := fmt.Sprintf("%s/%s", suite, tc.Name)
			t.Run(name, func(t *testing.T) {
				var vm VM
				code := strings.Split(tc.Code, "\n")
				vm.Execute(code)
				if tc.Error {
					if vm.err == nil {
						t.Fatal("expected error")
					}
					return
				}

				if vm.err != nil {
					t.Fatalf("unexpected error: %v", vm.err)
				}

				if len(vm.stack) == 0 {
					t.Fatalf("missing value")
				}

				n := vm.Pop()
				if fmt.Sprintf("%.5f", n) != fmt.Sprintf("%.5f", tc.Out) {
					t.Fatalf("expected %#v, got %#v", tc.Out, n)
				}
			})
		}
	}
}
