package synthesijer

type SlotItem struct{
	Next *SlotItem
	Op string
	Dest string
	Src string
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
}

type Board struct{
	Next *Board
	Name string
	Variables *Variable
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
	b.NextSlotId++
	return slot
}

func (s *Slot) AddItem(item *SlotItem) *SlotItem{
	s.Items, item.Next = item, s.Items
	return item
}

func (b *Board) AddVariable(v *Variable) *Variable{
	b.Variables, v.Next = v, b.Variables
	return v
}

func (m *Module) AddBoard(b *Board) *Board{
	m.Boards, b.Next = b, m.Boards
	return b
}

