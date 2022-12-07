package code

import (
	"MyCompiler/src/code"
	"testing"
)

func TestMake(t *testing.T) {
	tests := []struct {
		op       code.Opcode
		operands []int
		expected []byte
	}{
		{code.OpConstant, []int{65534}, []byte{byte(code.OpConstant), 255, 254}},
	}

	for _, tt := range tests {
		instruction := code.Make(tt.op, tt.operands...)

		// 长度应该对得上
		// 对于OpConstant，有三个字节，指令本身一个字节，操作数两个字节
		if len(instruction) != len(tt.expected) {
			t.Errorf("instruction has wrong length. expect=%d, got=%d", len(tt.expected), len(instruction))
		}

		// 逐个检查每个字节
		for i, b := range tt.expected {
			if instruction[i] != b {
				t.Errorf("wrong byte at pos %d, expect=%d, got=%d", i, b, instruction[i])
			}
		}
	}
}