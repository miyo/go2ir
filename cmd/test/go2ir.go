package main

import (
	"fmt"
	"strings"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"github.com/codegangsta/cli"
	"github.com/miyo/go2ir/synthesijer"
)

func generate_var(dest *os.File, v *synthesijer.Variable){
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

func generate_board(dest *os.File, b *synthesijer.Board){
	dest.Write([]byte("  (BOARD INT " + b.Name + "\n"))
	dest.Write([]byte("    (VARIABLES \n"))
	fmt.Println(b)
	for v := b.Variables; v != nil; v = v.Next {
		generate_var(dest, v)
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

func generate(m synthesijer.Module){

	dest, err := os.Create(m.Name + ".ir")
	if err != nil {
		panic(err)
	}
	defer dest.Close()

	dest.Write([]byte("(MODULE " + m.Name + "\n"))
	dest.Write([]byte("  (VARIABLES \n"))
	for v := m.Variables; v != nil; v = v.Next {
		generate_var(dest, v)
	}
	dest.Write([]byte("  )\n"))
	for b := m.Boards; b != nil; b = b.Next {
		generate_board(dest, b)
	}
	dest.Write([]byte(")\n"))

}

func parse_block(block *ast.BlockStmt){
	
	for _, s := range block.List{
		fmt.Printf("statement %v(%T)\n", s, s)
		switch s.(type) {
		case *ast.AssignStmt:
			fmt.Println("### Assign")
		default:
			fmt.Println("### otherwise")
		}
	}

}

func parse(src string) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, src, nil, 0)
	if err != nil {
		panic(err)
	}

	target_name := src[:strings.LastIndex(src, ".")]
	target := synthesijer.Module{Name: target_name}

	for _, decl := range file.Decls {
		switch td := decl.(type) {
		case *ast.GenDecl:
			switch td.Tok {
			case token.IMPORT:
				fmt.Println("### import")
				for _, sp := range td.Specs {
					s := sp.(*ast.ImportSpec)
					fmt.Println(s.Path.Value)
				}
			case token.TYPE:
				fmt.Println("### type")
				for _, sp := range td.Specs {
					s := sp.(*ast.TypeSpec)
					fmt.Println(s.Name)

					switch t := s.Type.(type) {
					case *ast.InterfaceType:
						for _, m := range t.Methods.List {
							fmt.Println(m)
						}
					case *ast.StructType:
						for _, f := range t.Fields.List {
							fmt.Println(f)
						}
					default:
						fmt.Println(3, t)
					}
				}
			case token.CONST:
			case token.VAR:
				fmt.Println("### var")
				for _, sp := range td.Specs {
					s := sp.(*ast.ValueSpec)
					fmt.Println(s.Names)
					fmt.Println(s.Type)
					fmt.Printf("type= %T(%v)\n", s.Type, s.Type)
				}

			default:

			}
		case *ast.FuncDecl:
			fmt.Println("### function")
			b := synthesijer.Board{Name: td.Name.Name}
			if target.Boards == nil{
				target.Boards = &b
			}else{
				b.Next = target.Boards
				target.Boards = &b
			}
			fmt.Println(td.Name)
			if td.Recv != nil {
				fmt.Println(td.Recv.List[0].Type)
			}
			if td.Type.Params != nil && td.Type.Params.NumFields() > 0 {
				fmt.Println("##### args")
				for _, p := range td.Type.Params.List {
					fmt.Println(p.Type, p.Names)
					for _, n := range p.Names {
						v := synthesijer.Variable{Name: n.Name, MethodParam: true, OriginalName: n.Name, MethodName: td.Name.Name}
						if b.Variables == nil {
							b.Variables = &v
						}else{
							v.Next = b.Variables
							b.Variables = &v
						}
					}
				}
			}
			if td.Type.Results != nil && td.Type.Results.NumFields() > 0 {
				fmt.Println("##### returns")
				for _, r := range td.Type.Results.List {
					fmt.Println(r.Type, r.Names)
				}
			}
			if td.Body != nil {
				parse_block(td.Body)
			}
		default:
		}

		fmt.Println()
	}
	generate(target)

}

func main() {
	app := cli.NewApp()
	app.Name = "go2ir"
	app.Usage = "Generating Synthesijer-IR from Go programming"
	app.Version = "0.1.1"
	app.Action = func(c *cli.Context) {
		if len(c.Args()) == 0 {
			fmt.Println("usage: go2ir sources")
			return
		}
		for _, src := range c.Args(){
			parse(src)
		}
	}
	app.Run(os.Args)
}

