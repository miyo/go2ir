package synthesijer

import (
	"fmt"
	"os"
)

func GenerateSlotItem(dest *os.File, item *SlotItem) string{
	str := "("
	str += fmt.Sprintf("%v ", item.Op)
	if item.Dest != nil {
		str += fmt.Sprintf("%v ", item.Dest.Name)
	}
	if item.Src != nil {
		str += fmt.Sprintf("%v ", item.Src.ToSExp())
	}
	str += ":next ("
	for _,id := range item.StepIds{
		str += fmt.Sprintf(" %v", id)
	}
	str += " ))\n"
	return str
}

func GenerateSlot(dest *os.File, slot *Slot){
	str := fmt.Sprintf("(SLOT %v\n", slot.Id)
	for item := slot.Items; item != nil; item = item.Next{
		str += GenerateSlotItem(dest, item)
	}
	str += ")\n"
	dest.Write([]byte(str))
}

func GenerateVariable(dest *os.File, v *Variable){
	s := ""
	if v.Constant == true {
		s += fmt.Sprintf("(CONSTANT ")
		s += v.Type + " "
		s += v.Name + " "
		s += v.Init + " "
		s += ")\n"
	}else{
		s += fmt.Sprintf("(VAR ")
		s += v.Type + " "
		s += v.Name + " "
		s += fmt.Sprintf(":public %v ", v.PublicFlag)
		s += fmt.Sprintf(":global_constant %v ", v.GlobalConstant)
		s += fmt.Sprintf(":method_param %v ", v.MethodParam)
		if v.OriginalName != "" {
			s += ":original " + v.OriginalName + " "
		}
		s += ":method " + v.MethodName + " "
		s += fmt.Sprintf(":private_method %v ", v.PrivateMethodFlag)
		s += fmt.Sprintf(":member %v ", v.MemberFlag)
		if v.Init != "" {
			s += fmt.Sprintf(":init %v ", v.Init)
		}
		s += ")\n"
	}
	dest.Write([]byte(s))
}

func GenerateVariableRef(dest *os.File, v *VariableRef){
	s := ""
	s += fmt.Sprintf("(VAR-REF ")
	s += v.Type + " "
	s += v.Name + " "
	s += fmt.Sprintf(":ref %v ", v.Ref)
	s += fmt.Sprintf(":ptr %v ", v.Ptr)
	s += fmt.Sprintf(":member %v", v.MemberFlag)
	s += ")\n"
	dest.Write([]byte(s))
}

func GenerateBoard(dest *os.File, b *Board){
	dest.Write([]byte("  (BOARD " + b.Type + " " + b.Name + "\n"))
	dest.Write([]byte("    (VARIABLES \n"))
	for v := b.Variables; v != nil; v = v.Next {
		GenerateVariable(dest, v)
	}
	for v := b.VariableRefs; v != nil; v = v.Next {
		GenerateVariableRef(dest, v)
	}
	dest.Write([]byte("    )\n"))
	dest.Write([]byte("    (SEQUENCER " + b.Name + "\n"))
	
	slots,_ := reverseSlots(b.Slots)
	for v := slots; v != nil; v = v.Next{
		GenerateSlot(dest, v)
	}
	
	dest.Write([]byte("    )\n"))
	dest.Write([]byte("  )\n"))
}

func GenerateArrayRef(dest *os.File, v *ArrayRef){
	s := ""
	s += fmt.Sprintf("(ARRAY-REF ")
	s += "ARRAY "
	s += v.Name + " "
	s += fmt.Sprintf(":depth %v ", v.Depth)
	s += fmt.Sprintf(":words %v ", v.Words)
	s += ")\n"
	dest.Write([]byte(s))
}

func GenerateModule(m *Module, destfile string){

	dest, err := os.Create(destfile + ".ir")
	if err != nil {
		panic(err)
	}
	defer dest.Close()

	dest.Write([]byte("(MODULE " + m.Name + "\n"))
	dest.Write([]byte("  (VARIABLES \n"))
	for v := m.ArrayRefs; v != nil; v = v.Next {
		GenerateArrayRef(dest, v)
	}
	for v := m.Variables; v != nil; v = v.Next {
		GenerateVariable(dest, v)
	}
	dest.Write([]byte("  )\n"))
	for b := m.Boards; b != nil; b = b.Next {
		GenerateBoard(dest, b)
	}
	dest.Write([]byte(")\n"))

}

func reverseSlots(s *Slot) (*Slot, *Slot){
	if s == nil {
		return nil, nil
	}else{
		s0, s1 := reverseSlots(s.Next)
		if s0 == nil {
			return s, s
		}else{
			s1.Next, s.Next = s, nil
			return s0, s
		}
	}
}
