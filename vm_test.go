package veem

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type testCase struct {
	name    string
	program []Inst
	out     int
	err     bool
}

func parseCode(code string) (Inst, error) {
	var name string
	var val int

	i, err := fmt.Sscanf(code, "%s %d", &name, &val)
	if err != nil && !errors.Is(err, io.EOF) {
		return nil, err
	}

	switch strings.ToUpper(name) {
	case "ADD":
		if i != 1 {
			return nil, fmt.Errorf("add does not take a value")
		}
		return Add, nil
	case "SUB":
		if i != 1 {
			return nil, fmt.Errorf("sub does not take a value")
		}
		return Sub, nil
	case "MUL":
		if i != 1 {
			return nil, fmt.Errorf("mul does not take a value")
		}
		return Mul, nil
	case "DIV":
		if i != 1 {
			return nil, fmt.Errorf("div does not take a value")
		}
		return Div, nil
	case "PUSH":
		if i != 2 {
			return nil, fmt.Errorf("push missing value")
		}

		return Push(val), nil
	case "MOD":
		if i != 1 {
			return nil, fmt.Errorf("add does not take a value")
		}

		return Mod, nil
	case "CMP":
		if i != 1 {
			return nil, fmt.Errorf("add does not take a value")
		}

		return Cmp, nil
	}

	return nil, fmt.Errorf("%q - unknown op", name)
}

func loadCase(t *testing.T, fileName string) testCase {
	file, err := os.Open(fileName)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	var tc struct {
		Program []string
		Out     int
		Error   bool
	}

	if err := json.NewDecoder(file).Decode(&tc); err != nil {
		t.Fatal(err)
	}

	prog := make([]Inst, 0, len(tc.Program))
	for _, v := range tc.Program {
		op, err := parseCode(v)
		if err != nil {
			t.Fatal(err)
		}
		prog = append(prog, op)
	}

	return testCase{
		name:    filepath.Base(fileName[:len(fileName)-5]),
		program: prog,
		out:     tc.Out,
		err:     tc.Error,
	}
}

func TestVM_Execute(t *testing.T) {
	files, err := filepath.Glob("testdata/*.json")
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

			if val != tc.out {
				t.Fatalf("expected %v, got %v", tc.out, val)
			}
		})
	}
}
