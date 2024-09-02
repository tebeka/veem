package veem

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

var executeCases = []struct {
	name    string
	program []Op
	val     int8
	err     bool
}{
	{
		"add",
		[]Op{
			PushOp | 3<<8,
			PushOp | 5<<8,
			AddOp,
		},
		8, false,
	},
	{
		"sub",
		[]Op{
			PushOp | 3<<8,
			PushOp | 5<<8,
			SubOp,
		},
		int8(-2), false,
	},
	{
		"add sub",
		[]Op{
			PushOp | 3<<8,
			PushOp | 5<<8,
			PushOp | 8<<8,
			AddOp,
			SubOp,
		},
		6, false,
	},
}

type testCase struct {
	name    string
	program []Op
	val     int8
	err     bool
}

func parseOp(op string) (Op, error) {
	var name string
	var val int8
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

		return Op(val | int8(PushOp)), nil
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
		err     bool `yaml:"error"`
	}
	if err := yaml.NewDecoder(file).Decode(&tc); err != nil {
		t.Fatal(err)
	}
}

func TestVM_Execute(t *testing.T) {
	for _, tc := range executeCases {
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
