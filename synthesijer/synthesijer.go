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

type CallExpr struct{
	Name string
	Args *Expr
	NoWait bool
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

func (e CallExpr) ToSExp() string{
	return fmt.Sprintf("(CALL :no_wait %v :name %v :args ())", e.NoWait, e.Name) // TODO
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
	Board *Board
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

type ArrayRef struct{
	Next *ArrayRef
	Name string
	Depth int
	Words int
}

type Board struct{
	Next *Board
	Name string
	Type string
	Variables *Variable
	VariableRefs *VariableRef
	Slots *Slot
	NextSlotId int
	Module *Module
}

type Module struct{
	Name string
	Variables *Variable
	ArrayRefs *ArrayRef
	Boards *Board
	UniqId int
}

func (m *Module) getUniqId() int{
	defer func(){m.UniqId++}()
	return m.UniqId
}

func (b *Board) AddSlot(slot *Slot) *Slot{
	b.Slots, slot.Next = slot, b.Slots
	defer func(){b.NextSlotId++}()
	slot.Board = b
	return slot
}

func (s *Slot) AddItem(item *SlotItem) *SlotItem{
	s.Items, item.Next = item, s.Items
	return item
}

func (b *Board) AddVariable(v *Variable) *Variable{
	b.Variables, v.Next = v, b.Variables
	v.Constant = false
	v.MethodName = b.Name
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

func (m *Module) AddArrayRef(v *ArrayRef) *ArrayRef{
	m.ArrayRefs, v.Next = v, m.ArrayRefs
	return v
}

func (m *Module) AddConstant(v *Variable) *Variable{
	m.Variables, v.Next = v, m.Variables
	v.Constant = true
	return v
}

func (b *Board) searchVariable(n string) *Variable{
	for v := b.Variables; v != nil; v = v.Next {
		if(v.OriginalName == n){
			return v
		}
	}
	return nil
}
