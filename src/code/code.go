package code

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type Instructions []byte
type Opcode byte

// 操作码的定义
const (
	OpConstant Opcode = iota
	OpAdd
)

type Definition struct {
	Name          string // 操作码名称
	OperandWidths []int  // 包含每个操作数占用的字节数
}

// 所有操作码定义的字典
var definitions = map[Opcode]*Definition{
	OpConstant: {"OpConstant", []int{2}}, // OpConstant唯一的操作数有两字节宽
	OpAdd:      {"OpAdd", []int{}},       // add操作没有操作数
}

func Lookup(op Opcode) (*Definition, error) {
	// 从definitions字典中查找OpCode
	def, ok := definitions[op]
	if !ok {
		// 如果遇到识别不了的，返回空byte数组
		return nil, fmt.Errorf("invalid opcode: %d", op)
	}
	return def, nil
}

// Make 按操作码 操作数编码
func Make(op Opcode, operands ...int) []byte {
	// 从definitions字典中查找OpCode
	def, err := Lookup(op)
	if err != nil {
		// 如果遇到识别不了的，返回空byte数组
		return []byte{}
	}

	// 因为操作码占1字节，所以最开始的值为1
	instructionLen := 1
	for _, w := range def.OperandWidths {
		// 所有操作数的长度相加，得到整条指令的长度
		instructionLen += w
	}

	// 创建指令数组
	instruction := make([]byte, instructionLen)
	// 指令的第一个字节就是操作码
	instruction[0] = byte(op)

	// 当前位置
	offset := 1
	// 挨个处理每个操作数
	for i, o := range operands {
		// 获取第i个操作数的宽度
		width := def.OperandWidths[i]
		switch width {
		case 2:
			// 如果操作数的宽度是2，用大端序把操作数写入instruction
			binary.BigEndian.PutUint16(instruction[offset:], uint16(o))
		}
		// 移到下一个位置
		offset += width
	}
	// 返回指令数组
	return instruction
}

// ReadOperands 从指令中读取操作数
// 返回操作数数组和操作数总长度
func ReadOperands(def *Definition, ins Instructions) ([]int, int) {
	operands := make([]int, len(def.OperandWidths))
	offset := 0

	for i, width := range def.OperandWidths {
		switch width {
		case 2:
			operands[i] = int(ReadUint16(ins[offset:]))
		}
		offset += width
	}
	return operands, offset
}

func ReadUint16(ins Instructions) uint16 {
	return binary.BigEndian.Uint16(ins)
}

func (self Instructions) String() string {
	var out bytes.Buffer

	i := 0
	for i < len(self) {
		def, err := Lookup(Opcode(self[i]))
		if err != nil {
			fmt.Fprintf(&out, "ERROR: %s\n", err)
			continue
		}

		operands, read := ReadOperands(def, self[i+1:])

		fmt.Fprintf(&out, "%04d %s\n", i, self.fmtInstruction(def, operands))

		// i移动到下一个位置
		i += 1 + read
	}
	return out.String()
}

func (self Instructions) fmtInstruction(def *Definition, operands []int) string {
	// 有几个操作数
	operandCount := len(def.OperandWidths)

	if len(operands) != operandCount {
		return fmt.Sprintf("ERROR: operand len %d does not match defined %d\n",
			len(operands), operandCount)
	}

	switch operandCount {
	case 1:
		return fmt.Sprintf("%s %d", def.Name, operands[0])
	case 0:
		return def.Name // 如果没有操作数就直接返回指令名称
	}

	return fmt.Sprintf("ERROR: unhandled operandCount for %s\n", def.Name)
}
