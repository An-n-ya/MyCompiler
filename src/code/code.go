package code

import "encoding/binary"

type Instructions []byte
type Opcode byte

// 操作码的定义
const (
	OpConstant Opcode = iota
)

type Definition struct {
	Name          string // 操作码名称
	OperandWidths []int  // 包含每个操作数占用的字节数
}

// 所有操作码定义的字典
var definitions = map[Opcode]*Definition{
	OpConstant: {"OpConstant", []int{2}}, // OpConstant唯一的操作数有两字节宽
}

func Make(op Opcode, operands ...int) []byte {
	// 从definitions字典中查找OpCode
	def, ok := definitions[op]
	if !ok {
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