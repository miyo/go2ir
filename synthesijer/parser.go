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
			//fmt.Println("### function", td.Name)
			b := target.AddBoard(&Board{Name: td.Name.Name})
			b.AddSlot(&Slot{Id: b.NextSlotId}).AddItem(&SlotItem{Op: "METHOD_EXIT", StepIds: []int{1}})
			b.AddSlot(&Slot{Id: b.NextSlotId}).AddItem(&SlotItem{Op: "METHOD_ENTRY", StepIds: []int{2}})
			if td.Recv != nil {
				fmt.Println(td.Recv.List[0].Type)
			}
			if td.Type.Params != nil && td.Type.Params.NumFields() > 0 {
				for _, p := range td.Type.Params.List {
					for _, n := range p.Names {
						t := convTypeFromExpr(p.Type)
						b.AddVariable(&Variable{Name: n.Name, MethodParam: true, OriginalName: n.Name, MethodName: td.Name.Name, Type: t})
					}
				}
			}
			if td.Type.Results != nil && td.Type.Results.NumFields() > 0 {
				b.Type = convTypeFromFieldList(td.Type.Results)
			}else{
				b.Type = "VOID"
			}
			if td.Body != nil {
				ParseBlock(b, td.Body)
			}
			slot := b.AddSlot(&Slot{Id: b.NextSlotId})
			slot.Items = &SlotItem{Op: "JP", StepIds: []int{0}}
		default:
		}

	}
}

func convTypeFromExpr(e ast.Expr) string{
	switch t := e.(type){
	case *ast.Ident:
		return convType(fmt.Sprint(t))
	case *ast.ArrayType:
		elm := convTypeFromExpr(t.Elt)
		return "ArrayType::" + elm
	case *ast.ChanType:
		elm := convTypeFromExpr(t.Value)
		return "(CHANNEL " + elm + ")"
	default:
		fmt.Printf("convTypeFromExpr: %v(%T)\n", t, t)
		return "UNKNOWN"
	}
}

func convTypeFromFieldList(f *ast.FieldList) string{
	ret := ""
	if f.NumFields() == 1 {
		ret = convType(fmt.Sprintf("%v", f.List[0].Type))
	}else if f.NumFields() > 1{
		sep := "(MULTIPLE "
		for _, r := range f.List {
			ret += sep + convType(fmt.Sprintf("%v", r.Type))
			sep = " "
			//fmt.Println(r.Type, r.Names)
		}
		ret += ")"
	}else{
		ret = "VOID"
	}
	return ret
}

func convType(s string) string{
	switch s{
	case "int": return "INT"
	default: return "VOID"
	}
}

func ParseBlock(board *Board, block *ast.BlockStmt) *Slot{
	
	var slot  *Slot
	for _, s := range block.List{
		switch td := s.(type) {
		case *ast.AssignStmt:
			//fmt.Println("### Assign")
			//fmt.Printf("  lhs:(%v:%T)\n", td.Lhs, td.Lhs)
			//fmt.Printf("  rhs:(%v:%T)\n", td.Rhs, td.Rhs)
			slot = board.AddSlot(&Slot{Id: board.NextSlotId})
			for i,_ := range td.Lhs{
				//fmt.Printf("   lhs[%v]:(%v:%T)\n", i, td.Lhs[i], td.Lhs[i])
				//fmt.Printf("   rhs[%v]:(%v:%T)\n", i, td.Rhs[i], td.Rhs[i])
				//fmt.Printf("   tok:(%v:%T)\n", td.Tok, td.Tok)

				if td.Tok == token.DEFINE {
					init_val := board.AddConstant(&Variable{Name: "init_val", Type: "INT", Init: "0"})
					board.AddVariable(&Variable{Name: fmt.Sprint(td.Lhs[i]), Type: "INT", Init: fmt.Sprintf("(REF CONSTANT %v)", init_val.Name)})
				}else if td.Tok == token.ASSIGN {
					rhs, lhs := "", ""
					rhs, slot = ParseExpr(board, slot, td.Rhs[i])
					lhs, slot = ParseExpr(board, slot, td.Lhs[i])
					slot.AddItem(&SlotItem{Op: "SET", Dest: lhs, Src: fmt.Sprintf("(ASSIGN %v)", rhs), StepIds: []int{board.NextSlotId}})
				}else{
					ret := ""
					lhs := ""
					ret, slot = GenBinaryExpr(board, slot, td.Lhs[i], td.Rhs[i], td.Tok)
					lhs, slot = ParseExpr(board, slot, td.Lhs[i])
					slot.AddItem(&SlotItem{Op: "SET", Dest: lhs, Src: fmt.Sprintf("(ASSIGN %v)", ret), StepIds: []int{board.NextSlotId}})
				}
			}
		case *ast.SendStmt:
			fmt.Println("### Send")
			fmt.Printf("  chan:(%v:%T)\n", td.Chan, td.Chan)
			fmt.Printf("  value:(%v:%T)\n", td.Value, td.Value)

			slot = board.AddSlot(&Slot{Id: board.NextSlotId})
			slot.AddItem(&SlotItem{Op: "FIFO_WRITE", Dest: fmt.Sprint(td.Chan), Src: fmt.Sprint("(ASSIGN ", td.Value, " )"), StepIds: []int{board.NextSlotId}})
			
		case *ast.RangeStmt:
			//fmt.Println("### RangeStmt")
			//fmt.Printf("  key(%v:%T)\n", td.Key, td.Key)
			//fmt.Printf("  value(%v:%T)\n", td.Value, td.Value)
			//fmt.Printf("  x(%v:%T)\n", td.X, td.X) // value to range over
			//fmt.Printf("  body(%v:%T)\n", td.Body, td.Body)
			board.AddVariable(&Variable{Name: fmt.Sprint(td.Value), Type: "INT"})
			v := board.AddVariable(&Variable{Name: "range_index", Type: "INT"})
			v1 := board.AddVariable(&Variable{Name: "range_compare", Type: "INT"})
			v2 := board.AddVariable(&Variable{Name: "field_length", Type: "INT"})
			c0 := board.AddConstant(&Variable{Name: "const_zero", Type: "INT", Init: "0"})
			c1 := board.AddConstant(&Variable{Name: "const_one", Type: "INT", Init: "1"})

			setup_slot := board.AddSlot(&Slot{Id: board.NextSlotId}); slot = setup_slot // range setup
			compare_slot := board.AddSlot(&Slot{Id: board.NextSlotId}); slot = compare_slot // range cond
			cond_slot := board.AddSlot(&Slot{Id: board.NextSlotId}); slot = cond_slot // range cond
			body_entry := board.NextSlotId
			body_slot := ParseBlock(board, td.Body); slot = body_slot // range body
			ret_slot := board.AddSlot(&Slot{Id: board.NextSlotId}); slot = ret_slot // range update and return
			
			// range setup
			setup_slot.AddItem(&SlotItem{Op: "SET", Dest: v.Name, Src: fmt.Sprintf("(ASSIGN %v)", c0.Name), StepIds: []int{compare_slot.Id}})
			setup_slot.AddItem(&SlotItem{Op: "SET", Dest: v2.Name, Src: fmt.Sprintf("(FIELD_ACCESS :obj %v :name %v)", td.X, "length"), StepIds: []int{compare_slot.Id}})

			// range compare
			compare_slot.AddItem(&SlotItem{Op: "SET",
				                       Dest: v1.Name,
				                       Src: fmt.Sprintf("(LT %v %v)", v.Name, v2.Name),
				                       StepIds: []int{cond_slot.Id}})
			
			cond_slot.AddItem(&SlotItem{Op: "JT", Src: v1.Name, StepIds: []int{board.NextSlotId, body_entry}})

			ret_slot.AddItem(&SlotItem{Op: "SET", Dest: v.Name, Src: fmt.Sprintf("(ADD %v %v)", v.Name, c1.Name), StepIds: []int{compare_slot.Id}})
			ret_slot.AddItem(&SlotItem{Op: "SET", Dest: v2.Name, Src: fmt.Sprintf("(FIELD_ACCESS :obj %v :name %v)", td.X, "length"), StepIds: []int{compare_slot.Id}})
			ret_slot.AddItem(&SlotItem{Op: "JP", StepIds: []int{compare_slot.Id}})
			
		case *ast.ReturnStmt:
			//fmt.Println("### ReturnStmt")
			slot = board.AddSlot(&Slot{Id: board.NextSlotId})
			var returns []string
			returns, slot = ParseExprList(board, slot, td.Results)
			if len(returns) == 1 {
				slot.AddItem(&SlotItem{Op: "RETURN", Src:returns[0], StepIds: []int{0} })
			}else if len(returns) > 1 {
				for i, ret := range returns{
					slot.AddItem(&SlotItem{Op: "MULTI_RETURN", Dest:fmt.Sprint(i), Src:ret, StepIds: []int{0} })
				}
			}
		default:
			fmt.Println("### otherwise")
			fmt.Printf("statement %v(%T)\n", s, s)
		}
	}
	return slot
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
	case *ast.BasicLit:
		ret, ret_slot = td.Value, slot
	default:
		fmt.Println("### expr otherwise")
		fmt.Printf("expr %v(%T)\n", expr, expr)
	}
	return ret, ret_slot
}

func GenBinaryExpr(board *Board, slot *Slot, lhs ast.Expr, rhs ast.Expr, tok token.Token) (string, *Slot){
	rhs_str, lhs_str := "", ""
	new_slot := slot
	rhs_str, new_slot = ParseExpr(board, new_slot, rhs) // rhs
	lhs_str, new_slot = ParseExpr(board, new_slot, lhs) // lhs
	op := ParseOp(tok)
	v := board.AddVariable(&Variable{Name: "binary_expr", Type: "INT"})
	//fmt.Printf("(SET %v (%v %v %v))\n", v.Name, op, lhs_str, rhs_str)
	new_slot.AddItem(&SlotItem{Op: "SET", Dest: v.Name, Src: fmt.Sprintf("(%v %v %v)", op, lhs_str, rhs_str), StepIds: []int{board.NextSlotId}})
	new_slot = board.AddSlot(&Slot{Id: board.NextSlotId})
	return v.Name, new_slot
}

func ParseBinaryExpr(board *Board, slot *Slot, expr *ast.BinaryExpr) (string, *Slot){
	//fmt.Printf("BinaryExpr %v(%T)\n", expr, expr)
	return GenBinaryExpr(board, slot, expr.X, expr.Y, expr.Op)
}

func ParseIdent(expr *ast.Ident, slot *Slot) (string, *Slot){
	//fmt.Printf("Ident %v(%T)\n", expr, expr)
	return expr.Name, slot
}

func ParseOp(op token.Token) string{
	switch op {
	case token.ADD:
		return "ADD"
	case token.ADD_ASSIGN:
		return "ADD"
	default:
		fmt.Printf("(Op %v[%T])\n", op, op)
	}
	return fmt.Sprintf("(Op %v[%T])", op, op)
}
