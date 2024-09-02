package veem

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

type testCase struct {
	name    string
	program []Op
	val     int8
	err     bool
}

func parseOp(op string) (Op, error) {
	var name string
	var val Op
	i, _ := fmt.Sscanf(op, "%s %d", &name, &val)
	if i == 0 {
		return 0, fmt.Errorf("parse %q", op)
	}

	switch strings.ToUpper(name) {
	case "ADD":
		return AddOp, nil
	case "SUB":
		return SubOp, nil
	case "PUSH":
		if i != 2 {
			return 0, fmt.Errorf("push missing value")
		}
		val <<= 8

		return Op(val | PushOp), nil
	}

	return 0, fmt.Errorf("%q - unknown op", name)
}

func loadCase(t *testing.T, fileName string) testCase {
	file, err := os.Open(fileName)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	var tc struct {
		program []string
		val     int8
		error   bool
	}
	if err := yaml.NewDecoder(file).Decode(&tc); err != nil {
		t.Fatal(err)
	}

	prog := make([]Op, 0, len(tc.program))
	for _, v := range tc.program {
		op, err := parseOp(v)
		if err != nil {
			t.Fatal(err)
		}
		prog = append(prog, op)
	}

	return testCase{
		name:    filepath.Base(fileName[:len(fileName)-4]),
		program: prog,
		val:     tc.val,
		err:     tc.error,
	}
}

func TestVM_Execute(t *testing.T) {
	files, err := filepath.Glob("testdata/*.yml")
	if err != nil {
		t.Fatal(err)
	}

	for _, fileName := range files {
		tc := loadCase(t, fileName)
		t.Run(tc.name, func(t *testing.T) {
			var v VM
			val, err := v.Execute(tc.program)
			if tc.err {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if val != tc.val {
				t.Fatalf("expected %v, got %v", tc.val, val)
			}
		})
	}
}
