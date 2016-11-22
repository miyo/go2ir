package synthesijer

import "fmt"

type Expr interface{
	ToSExp() string
}

type BasicExpr struct{
	str string
}

type IdentExpr struct{
	str string
}

type BinaryExpr struct{
	op string
	lhs Expr
	rhs Expr
}

type VarExpr struct{
	Var *Variable
}

type AssignExpr struct{
	Var VarExpr
}

func (e BasicExpr) ToSExp() string{
	return e.str
}

func (e IdentExpr) ToSExp() string{
	return e.str
}

func (e BinaryExpr) ToSExp() string{
	return fmt.Sprintf("(%v %v %v)", e.op, e.lhs.ToSExp(), e.rhs.ToSExp())
}

func (e VarExpr) ToSExp() string{
	//return fmt.Sprintf("(ASSIGN %v)", e.Var.Name)
	return e.Var.Name
}

func (e AssignExpr) ToSExp() string{
	return fmt.Sprintf("(ASSIGN %v)", e.Var.Var.Name)
}

type SlotItem struct{
	Next *SlotItem
	Op string
	Dest *Variable
	Src Expr
	StepIds []int
}

type Slot struct{
	Next *Slot
	Id int
	Items *SlotItem
}

type Variable struct{
	Next *Variable
	Name string
	Type string
	PublicFlag, GlobalConstant, MethodParam bool
	OriginalName string
	MethodName string
	PrivateMethodFlag, VolatileFlag, MemberFlag bool
	Constant bool
	Init string
}

type VariableRef struct{
	Next *VariableRef
	Name string
	Type string
	Ref string
	Ptr string
	MemberFlag bool
}

type Board struct{
	Next *Board
	Name string
	Type string
	Variables *Variable
	VariableRefs *VariableRef
	Slots *Slot
	NextSlotId int
}

type Module struct{
	Name string
	Variables *Variable
	Boards *Board
}


func (b *Board) AddSlot(slot *Slot) *Slot{
	b.Slots, slot.Next = slot, b.Slots
	defer func(){b.NextSlotId++}()
	return slot
}

func (s *Slot) AddItem(item *SlotItem) *SlotItem{
	s.Items, item.Next = item, s.Items
	return item
}

func (b *Board) AddVariable(v *Variable) *Variable{
	b.Variables, v.Next = v, b.Variables
	v.Constant = false
	return v
}

func (b *Board) AddVariableRef(v *VariableRef) *VariableRef{
	b.VariableRefs, v.Next = v, b.VariableRefs
	return v
}

func (b *Board) AddConstant(v *Variable) *Variable{
	b.Variables, v.Next = v, b.Variables
	v.Constant = true
	return v
}

func (m *Module) AddBoard(b *Board) *Board{
	m.Boards, b.Next = b, m.Boards
	return b
}

func (m *Module) AddVariable(v *Variable) *Variable{
	m.Variables, v.Next = v, m.Variables
	v.Constant = false
	return v
}

func (m *Module) AddConstant(v *Variable) *Variable{
	m.Variables, v.Next = v, m.Variables
	v.Constant = true
	return v
}
