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
				if td.Type.Results.NumFields() == 1 {
					b.Type = convType(fmt.Sprintf("%v", td.Type.Results.List[0].Type))
				}else if td.Type.Results.NumFields() > 1{
					sep := "(MULTI "
					for _, r := range td.Type.Results.List {
						b.Type += sep + convType(fmt.Sprintf("%v", r.Type))
						sep = " "
						fmt.Println(r.Type, r.Names)
					}
					b.Type += ")"
				}else{
					b.Type = "VOID"
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

func convType(s string) string{
	switch s{
	case "int": return "INT"
	default: return "VOID"
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
			returns, slot := ParseExprList(board, slot, td.Results)
			if len(returns) == 1 {
				slot.AddItem(&SlotItem{Op: "RETURN", Src:returns[0], StepIds: []int{0} })
			}else if len(returns) > 1 {
				for _, ret := range returns{
					slot.AddItem(&SlotItem{Op: "MULTI_RETURN", Src:ret, StepIds: []int{0} })
				}
			}
		default:
			fmt.Println("### otherwise")
			fmt.Printf("statement %v(%T)\n", s, s)
		}
	}

}

func ParseExprList(board *Board, slot *Slot, exprs []ast.Expr) ([]string, *Slot){
	var results []string
	new_slot := slot
	var tmp string
	for _,expr := range exprs {
		tmp, new_slot = ParseExpr(board, slot, expr)
		results = append(results, tmp)
	}
	return results, new_slot
}

func ParseExpr(board *Board, slot *Slot, expr ast.Expr) (string, *Slot){
	var ret string
	ret_slot := slot
	switch td := expr.(type) {
	case *ast.BinaryExpr:
		ret, ret_slot = ParseBinaryExpr(board, slot, td)
	case *ast.Ident:
		ret, ret_slot = ParseIdent(td, slot)
	default:
		fmt.Println("### otherwise")
		fmt.Printf("expr %v(%T)\n", expr, expr)
	}
	return ret, ret_slot
}

func ParseBinaryExpr(board *Board, slot *Slot, expr *ast.BinaryExpr) (string, *Slot){
	fmt.Printf("BinaryExpr %v(%T)\n", expr, expr)
	new_slot := slot
	var rhs, lhs string
	rhs, new_slot = ParseExpr(board, slot, expr.Y) // rhs
	lhs, new_slot = ParseExpr(board, slot, expr.X) // lhs
	op := ParseOp(expr.Op)
	v := board.AddVariable(&Variable{Name: "binary_expr"})

	fmt.Printf("(SET %v (%v %v %v))\n", v.Name, op, lhs, rhs)
	new_slot.AddItem(&SlotItem{Op: "SET", Dest: v.Name, Src: fmt.Sprintf("(%v %v %v)", op, lhs, rhs), StepIds: []int{board.NextSlotId}})
	new_slot = board.AddSlot(&Slot{Id: board.NextSlotId})
	return v.Name, new_slot
}

func ParseIdent(expr *ast.Ident, slot *Slot) (string, *Slot){
	fmt.Printf("Ident %v(%T)\n", expr, expr)
	return expr.Name, slot
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
