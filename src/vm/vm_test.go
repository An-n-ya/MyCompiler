package vm

import (
	"MyCompiler/src/ast"
	"MyCompiler/src/compiler"
	"MyCompiler/src/lexer"
	"MyCompiler/src/object"
	"MyCompiler/src/parser"
	"fmt"
	"testing"
)

type vmTestCase struct {
	input    string
	expected interface{}
}

// region 测试函数

func TestIntegerArithmetic(t *testing.T) {
	tests := []vmTestCase{
		{"100", 100},
		{"65535", 65535},
		{"1 + 2", 2},
	}
	runVmTests(t, tests)
}

// region

// region 帮助函数
func parse(input string) *ast.Program {
	l := lexer.New(input)
	p := parser.New(l)
	return p.ParseProgram()
}

func runVmTests(t *testing.T, tests []vmTestCase) {
	t.Helper()

	for _, tt := range tests {
		program := parse(tt.input)

		comp := compiler.New()
		err := comp.Compile(program)
		if err != nil {
			t.Fatalf("compiler error: %s", err)
		}

		// 用字节码建立虚拟机线程
		vm := New(comp.Bytecode())
		err = vm.Run()
		if err != nil {
			t.Fatalf("vm error: %s", err)
		}

		stackElem := vm.StackTop()
		testExpectedObject(t, tt.expected, stackElem)
	}
}

// 按测试对象类型分派测试方法
func testExpectedObject(t *testing.T, expected interface{}, actual object.Object) {
	t.Helper()

	switch expected := expected.(type) {
	case int:
		err := testIntegerObject(int64(expected), actual)
		if err != nil {
			t.Errorf("testIntegerObject failed: %s", err)
		}
	}
}

// 测试整数对象
func testIntegerObject(expected int64, actual object.Object) error {
	result, ok := actual.(*object.Integer)
	if !ok {
		return fmt.Errorf("object is not Integer. got=%T (%+v) ", actual, actual)
	}

	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got=%d, expect=%d",
			result.Value, expected)
	}

	return nil
}

// endregion
