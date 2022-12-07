package vm

import (
	"MyCompiler/src/code"
	"MyCompiler/src/compiler"
	"MyCompiler/src/object"
	"fmt"
)

// 栈大小
const StackSize = 2048

// VM 虚拟机结构体
type VM struct {
	constants    []object.Object   // 常量池
	instructions code.Instructions // 指令集
	stack        []object.Object   // 虚拟机栈
	sp           int               // 栈指针
}

func New(bytecode *compiler.ByteCode) *VM {
	return &VM{
		instructions: bytecode.Instructions,
		constants:    bytecode.Constants,
		stack:        make([]object.Object, StackSize),
		sp:           0,
	}
}

// 返回栈顶对象
func (vm *VM) StackTop() object.Object {
	if vm.sp == 0 {
		return nil
	}
	return vm.stack[vm.sp-1]
}

func (vm *VM) Run() error {
	for ip := 0; ip < len(vm.instructions); ip++ {
		op := code.Opcode(vm.instructions[ip])
		// 分别处理每种操作码
		switch op {
		case code.OpConstant:
			// 获取常量索引
			constIndex := code.ReadUint16(vm.instructions[ip+1:])
			ip += 2
			// 找到常量，并压入栈中
			err := vm.push(vm.constants[constIndex])
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (vm *VM) push(obj object.Object) error {
	if vm.sp >= StackSize {
		// 栈溢出
		return fmt.Errorf("stack overflow")
	}

	// 压栈
	vm.stack[vm.sp] = obj
	vm.sp++
	return nil
}
