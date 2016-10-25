package synthesijer

import (
	"fmt"
	"os"
)

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
}

type Module struct{
	Name string
	Variables *Variable
	Boards *Board
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
	
	dest.Write([]byte("      (SLOT 0 \n"))
	dest.Write([]byte("        (METHOD_EXIT :next (1))\n"))
	dest.Write([]byte("      )"))
	
	dest.Write([]byte("      (SLOT 1 \n"))
	dest.Write([]byte("        (METHOD_ENTRY :next (2))\n"))
	dest.Write([]byte("      )"))
	
	dest.Write([]byte("      (SLOT 2 \n"))
	dest.Write([]byte("        (JP :next (0))\n"))
	dest.Write([]byte("      )"))
	
	dest.Write([]byte("    )\n"))
	dest.Write([]byte("  )\n"))
}

func GenerateModule(m *Module){

	dest, err := os.Create(m.Name + ".ir")
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

