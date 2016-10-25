package synthesijer

import (
	"fmt"
	"os"
)

func GenerateSlotItem(dest *os.File, item *SlotItem) string{
	str := fmt.Sprintf("(%v %v %v :next (", item.Op, item.Dest, item.Src);
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
	s := fmt.Sprintf("(VAR INT ")
	s += v.Name + " "
	s += fmt.Sprintf(":public %v ", v.PublicFlag)
	s += fmt.Sprintf(":global_constant %v ", v.GlobalConstant)
	s += fmt.Sprintf(":method_param %v ", v.MethodParam)
	s += ":original " + v.OriginalName + " "
	s += ":method " + v.MethodName + " "
	s += fmt.Sprintf(":private_method %v ", v.PrivateMethodFlag)
	s += fmt.Sprintf(":member %v ", v.MemberFlag)
	s += ")\n"
	dest.Write([]byte(s))
}

func GenerateBoard(dest *os.File, b *Board){
	dest.Write([]byte("  (BOARD INT " + b.Name + "\n"))
	dest.Write([]byte("    (VARIABLES \n"))
	fmt.Println(b)
	for v := b.Variables; v != nil; v = v.Next {
		GenerateVariable(dest, v)
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

func GenerateModule(m *Module, destfile string){

	dest, err := os.Create(destfile + ".ir")
	if err != nil {
		panic(err)
	}
	defer dest.Close()

	dest.Write([]byte("(MODULE " + m.Name + "\n"))
	dest.Write([]byte("  (VARIABLES \n"))
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