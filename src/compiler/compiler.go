package compiler

import (
	"MyCompiler/src/ast"
	"MyCompiler/src/code"
	"MyCompiler/src/object"
	"fmt"
)

type Compiler struct {
	instructions code.Instructions // 编译器编译后的指令存放在这里
	constants    []object.Object   // 编译器计算后的常量放在这里
}

type ByteCode struct {
	Instructions code.Instructions
	Constants    []object.Object
}

func New() *Compiler {
	return &Compiler{
		instructions: code.Instructions{},
		constants:    []object.Object{},
	}
}

// 进行编译
func (self *Compiler) Compile(node ast.Node) error {
	// 根据node的类别编译
	switch node := node.(type) {
	case *ast.Program:
		for _, s := range node.Statement {
			// 递归对每条语句编译
			err := self.Compile(s)
			if err != nil {
				return err
			}
		}
	case *ast.ExpressionStatement:
		err := self.Compile(node.Expression)
		if err != nil {
			return err
		}
	case *ast.InfixExpression:
		err := self.Compile(node.Left)
		if err != nil {
			return err
		}
		err = self.Compile(node.Right)
		if err != nil {
			return err
		}

		switch node.Operator {
		case "+":
			self.emit(code.OpAdd)
		default:
			return fmt.Errorf("unknown operator %s", node.Operator)
		}
	case *ast.IntegerLiteral:
		integer := &object.Integer{Value: node.Value}
		// 将integer加入常量池，并得到它的位置
		pos := self.addConstant(integer)
		// 将指令写入指令集, 操作数就是integer在常量池的索引
		self.emit(code.OpConstant, pos)
	}
	return nil
}

// 将编译结果转化成字节码结构输出
func (self *Compiler) Bytecode() *ByteCode {
	return &ByteCode{
		Instructions: self.instructions, // 将编译器生成的指令给到字节码结构
		Constants:    self.constants,    // 将编译器计算的常量给字节码结构
	}
}

// 将obj加入到常量池
// 返回ob在常量池的位置
func (self *Compiler) addConstant(obj object.Object) int {
	self.constants = append(self.constants, obj)
	return len(self.constants) - 1
}

// 将指令加入到Instructions中
// 返回instruction在指令集合中的位置
func (self *Compiler) addInstruction(ins []byte) int {
	// 记住原来的位置
	ret := len(self.instructions)
	self.instructions = append(self.instructions, ins...)
	return ret
}

// 将操作码和操作数转换成指令加入Instruction中
// 返回指令在指令集合中的位置
func (self *Compiler) emit(op code.Opcode, operands ...int) int {
	ins := code.Make(op, operands...)
	pos := self.addInstruction(ins)
	return pos
}
