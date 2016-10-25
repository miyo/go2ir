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
			b := target.AddBoard(&Board{Name: td.Name.Name})
			b.AddSlot(&Slot{Id: b.NextSlotId}).AddItem(&SlotItem{Op: "METHOD_EXIT", StepIds: []int{1}})
			b.AddSlot(&Slot{Id: b.NextSlotId}).AddItem(&SlotItem{Op: "METHOD_ENTRY", StepIds: []int{2}})
			fmt.Println(td.Name)
			if td.Recv != nil {
				fmt.Println(td.Recv.List[0].Type)
			}
			if td.Type.Params != nil && td.Type.Params.NumFields() > 0 {
				fmt.Println("##### args")
				for _, p := range td.Type.Params.List {
					fmt.Println(p.Type, p.Names)
					for _, n := range p.Names {
						b.AddVariable(&Variable{Name: n.Name, MethodParam: true, OriginalName: n.Name, MethodName: td.Name.Name})
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
				ParseBlock(b, td.Body)
			}
			slot := b.AddSlot(&Slot{Id: b.NextSlotId})
			slot.Items = &SlotItem{Op: "JP", StepIds: []int{0}}
		default:
		}

		fmt.Println()
	}

}

func ParseBlock(board *Board, block *ast.BlockStmt){
	
	for _, s := range block.List{
		switch td := s.(type) {
		case *ast.AssignStmt:
			fmt.Println("### Assign")
		case *ast.ReturnStmt:
			fmt.Println("### ReturnStmt")
			slot := board.AddSlot(&Slot{Id: board.NextSlotId})
			returns := ParseExprList(board, slot, td.Results)
			retSlot := board.AddSlot(&Slot{Id: board.NextSlotId})
			for _, ret := range returns{
				retSlot.AddItem(&SlotItem{Op: "RETURN", Src:ret, StepIds: []int{0} })
			}
		default:
			fmt.Println("### otherwise")
			fmt.Printf("statement %v(%T)\n", s, s)
		}
	}

}

func ParseExprList(board *Board, slot *Slot, exprs []ast.Expr) []string{
	var results []string
	for _,expr := range exprs {
		results = append(results, ParseExpr(board, slot, expr))
	}
	return results
}

func ParseExpr(board *Board, slot *Slot, expr ast.Expr) string{
	var ret string
	switch td := expr.(type) {
	case *ast.BinaryExpr:
		ret = ParseBinaryExpr(board, slot, td)
	case *ast.Ident:
		ret = ParseIdent(td)
	default:
		fmt.Println("### otherwise")
		fmt.Printf("expr %v(%T)\n", expr, expr)
	}
	return ret
}

func ParseBinaryExpr(board *Board, slot *Slot, expr *ast.BinaryExpr) string{
	fmt.Printf("BinaryExpr %v(%T)\n", expr, expr)
	rhs := ParseExpr(board, slot, expr.Y) // rhs
	lhs := ParseExpr(board, slot, expr.X) // lhs
	op := ParseOp(expr.Op)
	v := board.AddVariable(&Variable{Name: "binary_expr"})

	fmt.Printf("(SET %v (%v %v %v))\n", v.Name, op, lhs, rhs)
	slot.AddItem(&SlotItem{Op: "SET", Dest: v.Name, Src: fmt.Sprintf("(%v %v %v)", op, lhs, rhs), StepIds: []int{board.NextSlotId}})
	return v.Name
}

func ParseIdent(expr *ast.Ident) string{
	fmt.Printf("Ident %v(%T)\n", expr, expr)
	return expr.Name
}

func ParseOp(op token.Token) string{
	switch op {
	case token.ADD:
		return "ADD"
	default:
		fmt.Printf("(Op %v[%T])\n", op, op)
	}
	return fmt.Sprintf("(Op %v[%T])", op, op)
}
