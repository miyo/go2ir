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
				for _, sp := range td.Specs {
					s := sp.(*ast.ValueSpec)
					t := convTypeFromExpr(s.Type)
					for _, ss := range s.Names{
						target.AddVariable(&Variable{Name: ss.Name, OriginalName: ss.Name, Type: t, MethodName: "fefe"})
					}
				}
			case token.VAR:
				for _, sp := range td.Specs {
					s := sp.(*ast.ValueSpec)
					t := convTypeFromExpr(s.Type)
					for i, ss := range s.Names{
						//fmt.Printf("  var name=%v(%T) init=%v(%T) type=%v(%T)\n", s.Names[i], s.Names[i], s.Values[i], s.Values[i], s.Type, s.Type)
						//fmt.Println(convTypeFromExpr(s.Type))
						if convTypeFromExpr(s.Type) == "ArrayType::INT" {// TODO
							ar := target.AddArrayRef(&ArrayRef{Name: fmt.Sprintf("array_%04d", target.getUniqId()), Depth: 4, Words: 16})
							target.AddVariable(&Variable{Name: ss.Name, OriginalName: ss.Name, Type: t, MethodName: "fefe", PublicFlag: true, Init: fmt.Sprintf("(REF ARRAY %s)", ar.Name)})
						}else{
							ParseExpr(nil, nil, s.Values[i])
							target.AddVariable(&Variable{Name: ss.Name, OriginalName: ss.Name, Type: t, MethodName: "fefe", PublicFlag: true})
						}
					}
				}
			default:

			}
		case *ast.FuncDecl:
			//fmt.Println("### function", td.Name)
			b := target.AddBoard(&Board{Name: td.Name.Name, Module: target})
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
			fmt.Println("### otherwise")
			fmt.Printf("statement %v(%T)\n", decl, decl)
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
					init_val := board.AddConstant(&Variable{Name: fmt.Sprintf("init_val_%04d", board.Module.getUniqId()), Type: "INT", Init: "0"})
					v := board.AddVariable(&Variable{Name: fmt.Sprint(td.Lhs[i]), OriginalName: fmt.Sprint(td.Lhs[i]), Type: "INT", Init: fmt.Sprintf("(REF CONSTANT %v)", init_val.Name)})
					slot.AddItem(&SlotItem{Op: "SET", Dest: v, Src: AssignExpr{VarExpr{init_val}}, StepIds: []int{board.NextSlotId}})
				}else if td.Tok == token.ASSIGN {
					var rhs Expr
					var lhs Expr
					rhs, slot = ParseExpr(board, slot, td.Rhs[i])
					lhs, slot = ParseExpr(board, slot, td.Lhs[i])
					//slot.AddItem(&SlotItem{Op: "SET", Dest: lhs, Src: fmt.Sprintf("(ASSIGN %v)", rhs), StepIds: []int{board.NextSlotId}})
					// TODO 本当はlhsがIdentであることを確認して，定義済みのVariableインスタンスを探すべき
					switch e := rhs.(type){
					case VarExpr:
						slot.AddItem(&SlotItem{Op: "SET", Dest: &Variable{Name: lhs.ToSExp()}, Src: AssignExpr{e}, StepIds: []int{board.NextSlotId}})
					default:
						slot.AddItem(&SlotItem{Op: "SET", Dest: &Variable{Name: lhs.ToSExp()}, Src: rhs, StepIds: []int{board.NextSlotId}})
					}
					
				}else{
					var ret Expr
					var lhs Expr
					ret, slot = GenBinaryExpr(board, slot, td.Lhs[i], td.Rhs[i], td.Tok)
					lhs, slot = ParseExpr(board, slot, td.Lhs[i])
					// TODO 本当はlhsがIdentであることを確認して，定義済みのVariableインスタンスを探すべき
					switch e := ret.(type){
					case VarExpr:
						slot.AddItem(&SlotItem{Op: "SET", Dest: &Variable{Name: lhs.ToSExp()}, Src: AssignExpr{e}, StepIds: []int{board.NextSlotId}})
					default:
						slot.AddItem(&SlotItem{Op: "SET", Dest: &Variable{Name: lhs.ToSExp()}, Src: ret, StepIds: []int{board.NextSlotId}})
					}
				}
			}
		case *ast.SendStmt:
			//fmt.Println("### Send")
			//fmt.Printf("  chan:(%v:%T)\n", td.Chan, td.Chan)
			//fmt.Printf("  value:(%v:%T)\n", td.Value, td.Value)

			slot = board.AddSlot(&Slot{Id: board.NextSlotId})
			var rhs Expr
			rhs,slot = ParseExpr(board, slot, td.Value)
			// TODO 本当はtd.Chanに相当す定義済みのVariableインスタンスを探すべき
			switch e := rhs.(type){
			case VarExpr:
				slot.AddItem(&SlotItem{Op: "FIFO_WRITE", Dest: &Variable{Name: fmt.Sprint(td.Chan)}, Src: AssignExpr{e}, StepIds: []int{board.NextSlotId}})
			default:
				slot.AddItem(&SlotItem{Op: "FIFO_WRITE", Dest: &Variable{Name: fmt.Sprint(td.Chan)}, Src: rhs, StepIds: []int{board.NextSlotId}})
			}
			
		case *ast.GoStmt:
			//fmt.Println("### GoStmt")
			slot = board.AddSlot(&Slot{Id: board.NextSlotId})
			var ret CallExpr
			ret, slot = ParseCallExpr(board, slot, td.Call)
			
			vv := board.AddVariable(&Variable{Name: fmt.Sprintf("method_result_%04d", board.Module.getUniqId()), Type: "VOID"}) // TODO
			ret.NoWait = true
			slot.AddItem(&SlotItem{Op: "SET", Dest: vv, Src: ret, StepIds: []int{board.NextSlotId}})
			
		case *ast.ExprStmt:
			//fmt.Println("### ExprStmt")
			slot = board.AddSlot(&Slot{Id: board.NextSlotId})
			_, slot = ParseExpr(board, slot, td.X)

		case *ast.RangeStmt:
			//fmt.Println("### RangeStmt")
			//fmt.Printf("  key(%v:%T)\n", td.Key, td.Key)
			//fmt.Printf("  value(%v:%T)\n", td.Value, td.Value)
			//fmt.Printf("  x(%v:%T)\n", td.X, td.X) // value to range over
			//fmt.Printf("  body(%v:%T)\n", td.Body, td.Body)
			vv := board.AddVariable(&Variable{Name: fmt.Sprintf("%v_%04d", td.Value, board.Module.getUniqId()), OriginalName: fmt.Sprintf("%v", td.Value), Type: "INT"})
			v := board.AddVariable(&Variable{Name: fmt.Sprintf("range_index_%04d", board.Module.getUniqId()), Type: "INT"})
			v1 := board.AddVariable(&Variable{Name: fmt.Sprintf("range_compare_%04d", board.Module.getUniqId()), Type: "BOOLEAN"})
			v2 := board.AddVariable(&Variable{Name: fmt.Sprintf("field_length_%04d", board.Module.getUniqId()), Type: "INT"})
			c0 := board.AddConstant(&Variable{Name: fmt.Sprintf("const_zero_%04d", board.Module.getUniqId()), Type: "INT", Init: "0"})
			c1 := board.AddConstant(&Variable{Name: fmt.Sprintf("const_one_%04d", board.Module.getUniqId()), Type: "INT", Init: "1"})

			setup_slot := board.AddSlot(&Slot{Id: board.NextSlotId}); slot = setup_slot // range setup
			compare_slot := board.AddSlot(&Slot{Id: board.NextSlotId}); slot = compare_slot // range cond
			cond_slot := board.AddSlot(&Slot{Id: board.NextSlotId}); slot = cond_slot // range cond
			body_entry := board.NextSlotId
			body_slot := ParseBlock(board, td.Body); slot = body_slot // range body
			ret_slot := board.AddSlot(&Slot{Id: board.NextSlotId}); slot = ret_slot // range update and return
			
			// range setup
			setup_slot.AddItem(&SlotItem{Op: "SET", Dest: v, Src: BasicExpr{fmt.Sprintf("(ASSIGN %v)", c0.Name)}, StepIds: []int{compare_slot.Id}})
			setup_slot.AddItem(&SlotItem{Op: "SET", Dest: v2, Src: BasicExpr{fmt.Sprintf("(FIELD_ACCESS :obj %v :name %v)", td.X, "length")}, StepIds: []int{compare_slot.Id}})

			// range compare
			compare_slot.AddItem(&SlotItem{Op: "SET",
				                       Dest: v1,
				                       Src: BinaryExpr{"LT", IdentExpr{v.Name}, IdentExpr{v2.Name}},
				                       StepIds: []int{cond_slot.Id}})
			compare_slot.AddItem(&SlotItem{Op: "SET", Dest: vv, Src: BasicExpr{fmt.Sprintf("(ARRAY_ACCESS %v %v)", td.X, v.Name)}, StepIds: []int{compare_slot.Id}}) // TODO
			
			cond_slot.AddItem(&SlotItem{Op: "JT", Src: IdentExpr{v1.Name}, StepIds: []int{body_entry, board.NextSlotId}})

			ret_slot.AddItem(&SlotItem{Op: "SET", Dest: v, Src: BasicExpr{fmt.Sprintf("(ADD %v %v)", v.Name, c1.Name)}, StepIds: []int{compare_slot.Id}})
			ret_slot.AddItem(&SlotItem{Op: "SET", Dest: v2, Src: BasicExpr{fmt.Sprintf("(FIELD_ACCESS :obj %v :name %v)", td.X, "length")}, StepIds: []int{compare_slot.Id}})
			ret_slot.AddItem(&SlotItem{Op: "JP", StepIds: []int{compare_slot.Id}})
			
		case *ast.ReturnStmt:
			//fmt.Println("### ReturnStmt")
			slot = board.AddSlot(&Slot{Id: board.NextSlotId})
			var returns []Expr
			returns, slot = ParseExprList(board, slot, td.Results)
			if len(returns) == 1 {
				slot.AddItem(&SlotItem{Op: "RETURN", Src:returns[0], StepIds: []int{0} })
			}else if len(returns) > 1 {
				for i, ret := range returns{
					// TODO 本当はtd.Chanに相当す定義済みのVariableインスタンスを探すべき
					slot.AddItem(&SlotItem{Op: "MULTI_RETURN", Dest:&Variable{Name:fmt.Sprint(i)}, Src:ret, StepIds: []int{0} })
				}
			}
		default:
			fmt.Println("### otherwise")
			fmt.Printf("statement %v(%T)\n", s, s)
		}
	}
	return slot
}

func ParseExprList(board *Board, slot *Slot, exprs []ast.Expr) ([]Expr, *Slot){
	var results []Expr
	new_slot := slot
	var tmp Expr
	for _,expr := range exprs {
		tmp, new_slot = ParseExpr(board, slot, expr)
		results = append(results, tmp)
	}
	return results, new_slot
}

func ParseExpr(board *Board, slot *Slot, expr ast.Expr) (Expr, *Slot){
	var ret Expr
	ret_slot := slot
	switch td := expr.(type) {
	case *ast.BinaryExpr:
		ret, ret_slot = ParseBinaryExpr(board, slot, td)
	case *ast.Ident:
		ret, ret_slot = ParseIdent(td, slot)
	case *ast.BasicLit:
		ret, ret_slot = BasicExpr{td.Value}, slot
	case *ast.CallExpr:
		ret, ret_slot = ParseCallExpr(board, slot, td)
	case *ast.ChanType:
		fmt.Printf(" ** expr ChanType %v(%T)\n", expr, expr)
		ret, ret_slot = BasicExpr{""}, slot
	case *ast.SliceExpr:
		fmt.Printf(" ** expr Slice %v(%T)\n", expr, expr)
		ret, ret_slot = BasicExpr{""}, slot
	default:
		fmt.Printf(" ** expr otherwise expr %v(%T)\n", expr, expr)
	}
	return ret, ret_slot
}

func GenBinaryExpr(board *Board, slot *Slot, lhs ast.Expr, rhs ast.Expr, tok token.Token) (Expr, *Slot){
	var new_rhs, new_lhs Expr
	new_slot := slot
	new_rhs, new_slot = ParseExpr(board, new_slot, rhs) // rhs
	new_lhs, new_slot = ParseExpr(board, new_slot, lhs) // lhs
	op := ParseOp(tok)
	v := board.AddVariable(&Variable{Name: fmt.Sprintf("binary_expr_%04d", board.Module.getUniqId()), Type: "INT"})
	//fmt.Printf("(SET %v (%v %v %v))\n", v.Name, op, lhs_str, rhs_str)
	new_slot.AddItem(&SlotItem{Op: "SET", Dest: v, Src: BinaryExpr{op, new_lhs, new_rhs}, StepIds: []int{board.NextSlotId}})
	new_slot = board.AddSlot(&Slot{Id: board.NextSlotId})
	return VarExpr{v}, new_slot
}

func ParseBinaryExpr(board *Board, slot *Slot, expr *ast.BinaryExpr) (Expr, *Slot){
	//fmt.Printf("BinaryExpr %v(%T)\n", expr, expr)
	return GenBinaryExpr(board, slot, expr.X, expr.Y, expr.Op)
}

func ParseIdent(expr *ast.Ident, slot *Slot) (Expr, *Slot){
	v := slot.Board.searchVariable(expr.Name)
	if v != nil {
		return VarExpr{v}, slot
	}else {
		// TODO 本当はtd.Chanに相当する定義済みのVariableインスタンスを探すべき
		return VarExpr{&Variable{Name: expr.Name}}, slot
	}
}

func ParseCallExpr(board *Board, slot *Slot, expr *ast.CallExpr) (CallExpr, *Slot){

	fmt.Println("ParseCallExpr")
	fmt.Printf(" ** call fun=%v(%T), args=%v(%T), ellipsis=%v(%T)\n",
		expr.Fun, expr.Fun, expr.Args, expr.Args, expr.Ellipsis, expr.Ellipsis)

	for _, e := range expr.Args {
		ParseExpr(board, slot, e)
	}

	fmt.Printf("%v\n", expr.Fun)
	return CallExpr{Name: fmt.Sprintf("%v", expr.Fun)}, slot
}

func ParseOp(op token.Token) string{
	switch op {
	case token.ADD:
		return "ADD"
	case token.ADD_ASSIGN:
		return "ADD"
	case token.MUL_ASSIGN:
		return "MUL32"
	default:
		fmt.Printf("(Op %v[%T])\n", op, op)
	}
	return fmt.Sprintf("(Op %v[%T])", op, op)
}
