package synthesijer

import(
	"fmt"
	"go/ast"
	"go/token"
)

func Parse(file *ast.File, target *Module) {

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
			b := Board{Name: td.Name.Name}
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
						v := Variable{Name: n.Name, MethodParam: true, OriginalName: n.Name, MethodName: td.Name.Name}
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
				ParseBlock(td.Body)
			}
		default:
		}

		fmt.Println()
	}

}

func ParseBlock(block *ast.BlockStmt){
	
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
